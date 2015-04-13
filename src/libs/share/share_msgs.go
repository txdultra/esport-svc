package share

import (
	"dbs"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"libs"
	"libs/passport"
	"reflect"
	"strconv"
	"sync"
	"time"
	"utils"
	"utils/redis"
	"utils/ssdb"
)

const (
	SHARE_MEMBER_MSGS_FMT = "share.member.msgs:%d"       //朋友圈内容
	SHARE_MEMBER_MSGS_C   = "share.member.msgs.count:%d" //朋友圈列表数
	SHARE_LASTNEW_MSG_FMT = "share.member.newmsg:%d"
)

var msgTbls map[string]bool = make(map[string]bool)
var m_locker *sync.Mutex = new(sync.Mutex)
var sobj *Shares = new(Shares)

type ShareMsgs struct{}

func NewShareMsgService() *ShareMsgs {
	return &ShareMsgs{}
}

//获取某人的消息
func (n *ShareMsgs) Gets(uid int64, page int, size int, ts time.Time) (int, []*Share) {
	ss := []*Share{}
	if uid <= 0 {
		return 0, ss
	}
	c := 0
	err := ssdb.New(use_ssdb_share_db).Get(fmt.Sprintf(SHARE_MEMBER_MSGS_C, uid), &c)
	if err != nil {
		type Result struct {
			C int64
		}
		o := dbs.NewOrm(msg_db)
		tbl := n.hash_share_msg_tbl(uid, msg_tbl)
		sql := fmt.Sprintf("select count(*) as c from %s where uid=?", tbl)
		var res Result
		o.Raw(sql, uid).QueryRow(&res)
		ssdb.New(use_ssdb_share_db).Set(fmt.Sprintf(SHARE_MEMBER_MSGS_C, uid), res.C)
		c = int(res.C)
	}

	offset := page * size
	if offset <= mbox_share_length {
		col_name := fmt.Sprintf(SHARE_MEMBER_MSGS_FMT, uid)
		kvs, _ := redis.ZRevRangeByScore(nil, col_name, reflect.TypeOf(int64(0)), utils.TimeMillisecond(ts)-1, 0, "LIMIT", 0, size)

		//redis数据可能被清空,重新从db中获取数据后导入
		if c > 0 && len(kvs) == 0 {
			tbl := n.hash_share_msg_tbl(uid, msg_tbl)
			type Result struct {
				Sid int64
				Ts  int64
			}
			sql := fmt.Sprintf("select sid,ts from %s where uid=? order by ts desc limit %d", tbl, mbox_share_length)
			var res []Result
			o := dbs.NewOrm(msg_db)
			o.Raw(sql, uid).QueryRows(&res)
			if len(res) > 0 {
				ids := []interface{}{}
				for _, r := range res {
					ids = append(ids, r.Ts)
					ids = append(ids, r.Sid)
				}
				redis.ZMultiAdd(nil, col_name, ids...)
				kvs, _ = redis.ZRevRangeByScore(nil, col_name, reflect.TypeOf(int64(0)), utils.TimeMillisecond(ts)-1, 0, "LIMIT", 0, size)
			}
		}

		for _, key := range kvs {
			if id, ok := key.(*int64); ok {
				_s := sobj.Get(*id)
				if _s != nil {
					ss = append(ss, _s)
				}
			}
		}
	} else {
		type Result struct {
			Sid int64
		}
		o := dbs.NewOrm(msg_db)
		tbl := n.hash_share_msg_tbl(uid, msg_tbl)
		sql := fmt.Sprintf("select sid from %s where uid=? and ts <= ? order by ts desc limit %d", tbl, size)
		var res []Result
		o.Raw(sql, uid, utils.TimeMillisecond(ts)-1).QueryRows(&res)
		for _, r := range res {
			_s := sobj.Get(r.Sid)
			if _s != nil {
				ss = append(ss, _s)
			}
		}
	}
	return c, ss
}

//重置新消息列表 IEventCounter interface
func (n ShareMsgs) ResetEventCount(uid int64) bool {
	c, _ := redis.Del(nil, fmt.Sprintf(SHARE_LASTNEW_MSG_FMT, uid))
	return c > 0
}

//IEventCounter interface
func (n ShareMsgs) NewEventCount(uid int64) int {
	exist, _ := redis.Exists(nil, fmt.Sprintf(SHARE_LASTNEW_MSG_FMT, uid))
	if exist {
		return 1
	}
	return 0
}

//收到最新的一条新消
func (n *ShareMsgs) LastNewMsg(uid int64) *Share {
	key := fmt.Sprintf(SHARE_LASTNEW_MSG_FMT, uid)
	idstr, err := redis.Get(nil, key)
	if err != nil {
		return nil
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return nil
	}
	s := &Shares{}
	_s := s.Get(id)
	if _s != nil {
		return _s
	}
	return nil
}

//用户消息库分表
func (n *ShareMsgs) hash_share_msg_tbl(uid int64, pfx_table string) string {
	str := strconv.FormatInt(uid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	tbl := fmt.Sprintf("%s_%d", pfx_table, hs%500)
	if _, ok := msgTbls[tbl]; ok {
		return tbl
	}
	m_locker.Lock()
	defer m_locker.Unlock()
	if _, ok := msgTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(msg_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(uid int(11) NOT NULL,
		  sid bigint(20) NOT NULL,
		  suid int(11) NOT NULL,
		  ts bigint(15) NOT NULL,
		  st smallint(6) NOT NULL,
		  PRIMARY KEY (uid,sid,suid),
  		  KEY idx_sid (sid),
		  KEY idx_uid_st_ts (uid,st,ts)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		msgTbls[tbl] = true
	}
	return tbl
}

func (n ShareMsgs) SSE_Save(s *Share) {
	//朋友圈中
	groupcol := fmt.Sprintf(SHARE_MEMBER_MSGS_FMT, s.Uid)
	redis.ZAdd(nil, groupcol, utils.TimeMillisecond(s.CreateTime), s.Id)
	redis.ZRemRangeByRank(nil, groupcol, 0, -mbox_share_length) //限制box大小

	o := dbs.NewOrm(msg_db)
	tbl := n.hash_share_msg_tbl(s.Uid, msg_tbl)
	sql := fmt.Sprintf("insert into %s(uid,sid,suid,ts,st) values(?,?,?,?,?)", tbl)
	o.Raw(sql, s.Uid, s.Id, s.Uid, utils.TimeMillisecond(s.CreateTime), s.ShareType).Exec()

	ssdb.New(use_ssdb_share_db).Incr(fmt.Sprintf(SHARE_MEMBER_MSGS_C, s.Uid))
}

func (n ShareMsgs) SDE_Revoke(s *Share) {
	//o := dbs.NewOrm(msg_db)
	//tbl := n.hash_share_msg_tbl(s.Uid, msg_tbl)
	//sql := fmt.Sprintf("delete from %s where sid=? and uid=?", tbl)
	//o.Raw(sql, s.Id, s.Uid).Exec()

	//朋友圈中
	//groupcol := fmt.Sprintf(SHARE_MEMBER_MSGS_FMT, s.Uid)
	//redis.ZRem(nil, groupcol, s.Id)
}

func (n ShareMsgs) NotifyFans(s *Share) {
	pst := &passport.FriendShips{}

	//发送给双向关系的人
	//notify_uids := pst.BothFriendUids(s.Uid)
	ts := time.Now()
	page := 1
	size := 500
	o := dbs.NewOrm(msg_db)
	for {
		//发送给粉丝
		_, notify_uids, last_t := pst.FollowerIdsP(s.Uid, size, page, ts)
		if len(notify_uids) == 0 {
			return
		}
		for _, uid := range notify_uids {
			msg_key := fmt.Sprintf(SHARE_MEMBER_MSGS_FMT, uid)
			//加入朋友圈消息列表
			redis.ZAdd(nil, msg_key, utils.TimeMillisecond(s.CreateTime), s.Id)
			redis.ZRemRangeByRank(nil, msg_key, 0, -mbox_share_length) //限制box大小

			if sns_msg_db_insert_queue_open {
				go n.SendMsgToQueue(s, uid)
			} else {
				tbl := n.hash_share_msg_tbl(uid, msg_tbl)
				sql := fmt.Sprintf("insert into %s(uid,sid,suid,ts,st) values(?,?,?,?,?)", tbl)
				o.Raw(sql, uid, s.Id, s.Uid, utils.TimeMillisecond(s.CreateTime), s.ShareType).Exec()
			}

			ssdb.New(use_ssdb_share_db).Incr(fmt.Sprintf(SHARE_MEMBER_MSGS_C, uid))
			redis.Set(nil, fmt.Sprintf(SHARE_LASTNEW_MSG_FMT, uid), s.Id)
		}
		page++
		ts = utils.MsToTime(last_t)
	}
}

func (n *ShareMsgs) SendMsgToQueue(s *Share, toUid int64) {
	tbl := n.hash_share_msg_tbl(toUid, msg_tbl)
	sql := fmt.Sprintf("insert into %s(uid,sid,suid,ts,st) values(%d,%d,%d,%d,%d)", tbl, toUid, s.Id, s.Uid, utils.TimeMillisecond(s.CreateTime), s.ShareType)
	dbmsg := &dbMsg{
		DbName: "share_msgs",
		Sql:    sql,
	}
	body, _ := json.Marshal(dbmsg)
	msqCfg := &libs.MsqQueueConfig{
		MsqId:      "share_msg_db_batch_insert",
		MsqConnUrl: libs.MsqConnectionUrl(),
		QueueName:  msq_db_batch_queue_name,
		Durable:    true,
		AutoAck:    false,
		Exchange:   "",
		QueueMode:  libs.MsqWorkQueueMode,
	}
	msq := libs.NewMsq()
	msq_service, _ := msq.CreateMsqService(msqCfg)

	msg := &libs.MsqMessage{
		DataType: "json",
		Ts:       utils.TimeMillisecond(time.Now()),
		MsgType:  libs.MSG_TYPE_DBBATCH,
		Data:     body,
	}
	msq_service.Send(msg)
}

func (n ShareMsgs) RecoverFans(s *Share) {
	//不做操作

	//pst := &passport.FriendShips{}
	//both_uids := pst.BothFriendIds(s.Uid)
	//if len(both_uids) == 0 {
	//	return
	//}
	//for _, uid := range both_uids {
	//	msg_key := fmt.Sprintf(SHARE_MEMBER_MSGS_FMT, uid)
	//	//加入朋友圈消息列表
	//	redis.ZRem(nil, msg_key, s.Id)
	//}
}

//friendship interfaces
func (n ShareMsgs) FSFriendDo(suid int64, tuid int64) {
}

func (n ShareMsgs) FSUnFriendDo(suid int64, tuid int64) {
	//老数据依旧保留

	//不保留代码
	//o := dbs.NewOrm(msg_db)
	//tbl := n.hash_share_msg_tbl(tuid, msg_tbl)
	//sql := fmt.Sprintf("delete from %s where suid=%d and uid=%d", tbl_a, tuid, suid)

	//type Result struct {
	//	Sid int64
	//	Ts  int64
	//}

	//sql = fmt.Sprintf("select sid,ts from %s where uid=%d order by ts desc limit %d", tbl_a, suid, mbox_share_length)
	//var res []Result
	//o.Raw(sql).QueryRows(&res)

	//my_colname := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS, suid)
	//adds := []interface{}{}
	//for _, r := range res {
	//	adds = append(adds, r.Ts)
	//	adds = append(adds, r.Sid)
	//}
	//redis.ZMultiAdd(nil, my_colname, adds...)
}
