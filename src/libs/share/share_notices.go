package share

import (
	"crypto/md5"
	"crypto/rand"
	"dbs"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"hash/crc32"
	"io"
	"libs"
	"libs/passport"
	"libs/vod"
	"logs"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"utils"
	"utils/redis"
	"utils/ssdb"
)

const (
	SHARE_NEW_NOTICES = "share.new.notices:%s"
	SHARE_NOTICE      = "share.notice:%s"
	notice_db_mark    = "10"
)

var noticeTbls = make(map[string]bool)
var n_locker *sync.Mutex = new(sync.Mutex)

type ShareNotices struct{}

func NewShareNoticeService() *ShareNotices {
	return &ShareNotices{}
}

func (n *ShareNotices) tbl_tag(uid int64) string {
	str := strconv.FormatInt(uid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	mod := hs % 500
	mtag := ""
	if mod > 99 {
		mtag = fmt.Sprintf("%d", mod)
	} else if mod < 10 {
		mtag = fmt.Sprintf("00%d", mod)
	} else {
		mtag = fmt.Sprintf("0%d", mod)
	}
	return mtag
}

//用户通知库分表
func (n *ShareNotices) hash_share_notice_tbl(uid int64) string {
	mtag := n.tbl_tag(uid)
	tbl := fmt.Sprintf("%s_%s", share_notice_tbl, mtag)
	if _, ok := noticeTbls[tbl]; ok {
		return tbl
	}
	n_locker.Lock()
	defer n_locker.Unlock()
	if _, ok := noticeTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(share_notice_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(id char(30) NOT NULL,
		  sid bigint(20) NOT NULL,
		  uid int(11) NOT NULL,
		  luid int(11) NOT NULL,
		  lu_nickname char(25) NOT NULL,
		  ruid int(11) NOT NULL,
		  ru_nickname char(25) NOT NULL,
		  content varchar(200) NOT NULL,
		  pic int(11) NOT NULL,
		  ts bigint(15) NOT NULL,
		  t smallint(6) NOT NULL,
		  st smallint(6) NOT NULL,
		  PRIMARY KEY (id),
		  KEY idx_uid_ts (uid,ts)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		noticeTbls[tbl] = true
	}
	return tbl
}

func (n *ShareNotices) selDbTbl(id string) (db string, tbl string) {
	return id[:2], fmt.Sprintf("%s_%s", share_notice_tbl, id[2:5])
}

//id生成模式:db(2位)+tbl(2位)+ts(毫秒13位)=17位
func (n *ShareNotices) newId(uid int64) string {
	mtag := n.tbl_tag(uid)
	newid := n.hashId()
	id := fmt.Sprintf("%s%s%s", notice_db_mark, mtag, newid)
	return id
}

var machineId = readMachineId()
var objectIdCounter uint32 = 0

func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

func (n *ShareNotices) hashId() string {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return hex.EncodeToString(b[:])
}

func (n ShareNotices) getSharePic(s *Share) int64 {
	nts := &Shares{}
	ress := nts.GetResources(s.Resources)
	for _, res := range ress {
		if res.Kind == SHARE_KIND_VOD {
			_id, _ := strconv.ParseInt(res.Id, 10, 64)
			vods := &vod.Vods{}
			v := vods.Get(_id, false)
			if v == nil {
				return 0
			}
			return v.Img
		}
		if res.Kind == SHARE_KIND_PIC {
			_id, _ := strconv.ParseInt(res.Id, 10, 64)
			return _id
		}
	}
	return 0
}

//IEventCounter interface
func (n *ShareNotices) ResetEventCount(uid int64) bool {
	redis.Del(nil, fmt.Sprintf(SHARE_NEW_NOTICES, uid))
	return true
}

//IEventCounter interface
func (n *ShareNotices) NewEventCount(uid int64) int {
	c, _ := redis.ZCard(nil, fmt.Sprintf(SHARE_NEW_NOTICES, uid))
	return int(c)
}

//发布新分享时被提及到的对象得到通知
func (n ShareNotices) SSE_Save(s *Share) {
	refuids := s.GetRefUids()
	pst := passport.NewMemberProvider()
	o := dbs.NewOrm(share_notice_db)

	for _, uid := range refuids {
		sn := &ShareNotice{
			Id:         n.newId(uid),
			Uid:        uid,
			Sid:        s.Id,
			LUid:       s.Uid,
			LUNickname: pst.GetNicknameByUid(s.Uid),
			Content:    "提到了你",
			Pic:        n.getSharePic(s),
			Ts:         utils.TimeMillisecond(s.CreateTime),
			T:          SHARE_NOTICE_TYPE_TD,
			ST:         s.ShareType,
		}
		go n.addNotice(o, sn)
	}
}

func (n *ShareNotices) addNotice(o orm.Ormer, sn *ShareNotice) {
	if sns_notice_db_insert_queue_open {
		go n.sendNoticeToQueue(sn)
	} else {
		tbl := n.hash_share_notice_tbl(sn.Uid)
		_, err := o.Raw(fmt.Sprintf("insert into %s(id,sid,uid,luid,lu_nickname,ruid,ru_nickname,content,pic,ts,t,st) values(?,?,?,?,?,?,?,?,?,?,?,?)", tbl),
			sn.Id,
			sn.Sid,
			sn.Uid,
			sn.LUid,
			sn.LUNickname,
			sn.RUid,
			sn.RUNickname,
			sn.Content,
			sn.Pic,
			sn.Ts,
			sn.T,
			sn.ST).Exec()
		if err != nil {
			logs.Errorf("share's notice add to mysql fail:%s", err.Error())
		}
	}

	col_name := fmt.Sprintf(SHARE_NEW_NOTICES, sn.Uid)
	redis.ZAdd(nil, col_name, sn.Ts, sn.Id)
	redis.ZRemRangeByRank(nil, col_name, 0, -mbox_share_length) //限制box大小

	//存入ssdb
	ssdb.New(use_ssdb_notice_db).Set(fmt.Sprintf(SHARE_NOTICE, sn.Id), *sn)
}

func (n *ShareNotices) sendNoticeToQueue(sn *ShareNotice) {
	tbl := n.hash_share_notice_tbl(sn.Uid)
	sql := fmt.Sprintf("insert into %s(id,sid,uid,luid,lu_nickname,ruid,ru_nickname,content,pic,ts,t,st) values('%s',%d,%d,%d,'%s',%d,'%s','%s',%d,%d,%d,%d)",
		tbl,
		sn.Id,
		sn.Sid,
		sn.Uid,
		sn.LUid,
		sn.LUNickname,
		sn.RUid,
		sn.RUNickname,
		sn.Content,
		sn.Pic,
		sn.Ts,
		sn.T,
		sn.ST)
	dbmsg := &dbMsg{
		DbName: "share_notices",
		Sql:    sql,
	}
	body, _ := json.Marshal(dbmsg)
	msqCfg := &libs.MsqQueueConfig{
		MsqId:      "notice_db_batch_insert",
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

func (n *ShareNotices) getNoticeType(t SHARE_COMMENT_TYPE) SHARE_NOTICE_TYPE {
	switch t {
	case SHARE_COMMENT_TYPE_REPLY:
		return SHARE_NOTICE_TYPE_HF
	case SHARE_COMMENT_TYPE_SINGLE:
		return SHARE_NOTICE_TYPE_PL
	default:
		return SHARE_NOTICE_TYPE_TD
	}
}

func (n ShareNotices) SC_CommentedEvent(sc *ShareComment) {
	uidmaps := make(map[int64]bool)
	s := &Shares{}
	share := s.Get(sc.Sid)
	if share == nil {
		return
	}
	discussers := s.GetDiscussers(sc.Sid)

	for _, uid := range discussers {
		uidmaps[uid] = true
	}
	delete(uidmaps, sc.Uid) //自己不会收到新消息

	if len(uidmaps) == 0 {
		return
	}

	o := dbs.NewOrm(share_notice_db)

	for uid, _ := range uidmaps {
		is_send := true

		//相互关注才能收到
		//if sc.T == SHARE_COMMENT_TYPE_SINGLE {
		//	if passport.IsBothFriend(sc.Uid, uid) {
		//		is_send = true
		//	}
		//} else if sc.T == SHARE_COMMENT_TYPE_REPLY {
		//	if passport.IsBothFriend(sc.Uid, uid) {
		//		if sc.RUid == uid {
		//			is_send = true
		//		} else if passport.IsBothFriend(sc.RUid, uid) {
		//			is_send = true
		//		}
		//	}
		//}
		if is_send {
			sn := &ShareNotice{
				Id:         n.newId(uid),
				Uid:        uid,
				Sid:        sc.Sid,
				LUid:       sc.Uid,
				LUNickname: sc.UNickname,
				RUid:       sc.RUid,
				RUNickname: sc.RUNickname,
				Content:    sc.Content,
				Pic:        n.getSharePic(share),
				Ts:         sc.Ts,
				T:          n.getNoticeType(sc.T),
				ST:         share.ShareType,
			}
			go n.addNotice(o, sn)
		}
	}
}

func (n ShareNotices) SC_RevokedEvent(sc *ShareComment) {

}

func (n *ShareNotices) GetNewNotices(uid int64) []*ShareNotice {
	col_name := fmt.Sprintf(SHARE_NEW_NOTICES, uid)
	lst, _ := redis.ZRevRangeByScore(nil, col_name, reflect.TypeOf(""), utils.TimeMillisecond(time.Now())-1, 0)
	keys := []string{}
	for _, id := range lst {
		_id := *(id.(*string))
		key := fmt.Sprintf(SHARE_NOTICE, _id)
		keys = append(keys, key)
	}
	arrs := ssdb.New(use_ssdb_notice_db).MultiGet(keys, reflect.TypeOf(ShareNotice{}))
	sns := []*ShareNotice{}
	for _, ar := range arrs {
		sn, ok := ar.(*ShareNotice)
		if !ok {
			continue
		}
		sns = append(sns, sn)
	}
	return sns
}

func (n *ShareNotices) GetNotices(uid int64, page int, size int, ts time.Time) (int, []*ShareNotice) {
	o := dbs.NewOrm(share_notice_db)
	tbl := n.hash_share_notice_tbl(uid)

	type Result struct {
		C int64
	}
	sql := fmt.Sprintf("select count(*) as c from %s where uid=? and ts <=?", tbl)
	var res_a Result
	o.Raw(sql, uid, utils.TimeMillisecond(ts)-1).QueryRow(&res_a)
	c := int(res_a.C)

	offset := (page - 1) * size

	sql = fmt.Sprintf("select id,sid,uid,luid,lu_nickname,ruid,ru_nickname,content,pic,ts,t,st from %s where uid=? and ts <= ? order by ts desc limit %d,%d", tbl, offset, size)
	var res []*ShareNotice
	o.Raw(sql, uid, utils.TimeMillisecond(ts)-1).QueryRows(&res)
	return c, res
}

func (n *ShareNotices) DelNotices(uid int64, ids []string) {
	o := dbs.NewOrm(share_notice_db)
	tbl := n.hash_share_notice_tbl(uid)
	if len(ids) == 0 {
		return
	}
	sql_ids := []string{}
	keys := []string{}
	for _, id := range ids {
		if len(id) == 0 {
			continue
		}
		sql_ids = append(sql_ids, fmt.Sprintf("'%s'", id))
		keys = append(keys, fmt.Sprintf(SHARE_NOTICE, id))
	}
	sql := fmt.Sprintf("delete from %s where id in (%s)", tbl, strings.Join(sql_ids, ","))
	o.Raw(sql).Exec()
	ssdb.New(use_ssdb_notice_db).MultiDel(keys)
}

func (n *ShareNotices) EmptyNotices(uid int64) {
	o := dbs.NewOrm(share_notice_db)
	tbl := n.hash_share_notice_tbl(uid)
	sql := fmt.Sprintf("delete from %s where uid=?", tbl)
	o.Raw(sql, uid).Exec()
	n.ResetEventCount(uid)
}
