package share

import (
	"dbs"
	"fmt"
	"hash/crc32"
	"libs/hook"
	"libs/passport"
	"logs"
	"reflect"
	"strconv"
	"sync"
	"time"
	"utils"
	"utils/ssdb"
)

type ShareCommentEvent interface {
	SC_CommentedEvent(sc *ShareComment) //提交后
	SC_RevokedEvent(sc *ShareComment)
}

var commentEvents = make(map[string]ShareCommentEvent)

func RegisterShareCommentEvent(eventName string, eventHandler ShareCommentEvent) {
	if _, ok := commentEvents[eventName]; ok {
		return
	}
	commentEvents[eventName] = eventHandler
}

var cmtTbls map[string]bool = make(map[string]bool)
var c_locker *sync.Mutex = new(sync.Mutex)
var previewCmtCounts int = 5

const (
	SHARE_COMMENT_LIST_KEY = "share_cmt_ids:%d"
	SHARE_COMMENT_KEY      = "share_cmt:%d"
	cmt_db_mark            = "10" //库默认编号
	SHARE_COMMENTS_TOP     = "share_cmt_tops:%d"
)

type ShareComments struct{}

func NewShareComments() *ShareComments {
	return &ShareComments{}
}

//库分表
func (s *ShareComments) hash_cmt_tbl(sid int64) string {
	str := strconv.FormatInt(sid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	pfx := strconv.FormatUint(uint64(hs), 10)
	tbl := fmt.Sprintf("%s_%s", share_cmt_tbl, pfx[:2])
	if _, ok := cmtTbls[tbl]; ok {
		return tbl
	}
	c_locker.Lock()
	defer c_locker.Unlock()
	if _, ok := cmtTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(share_cmt_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(id bigint(18) NOT NULL,
	  sid bigint(20) NOT NULL,
	  suid int(11) NOT NULL,
	  su_nickname char(25) NOT NULL,
	  uid int(11) NOT NULL,
	  u_nickname char(25) NOT NULL,
	  ruid int(11) NOT NULL,
	  ru_nickname char(25) NOT NULL,
	  t smallint(6) NOT NULL,
	  content varchar(200) NOT NULL,
	  ts bigint(15) NOT NULL,
	  ex1 varchar(50) NOT NULL DEFAULT '0',
	  ex2 varchar(50) NOT NULL DEFAULT '0',
	  ex3 varchar(50) NOT NULL DEFAULT '0',
	  PRIMARY KEY (id),
	  KEY idx_sid_ts (sid,ts) USING BTREE) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		cmtTbls[tbl] = true
	} else {
		logs.Errorf("share's comment create table %s fail:%s", tbl, err.Error())
	}
	return tbl
}

func (s *ShareComments) selDbTbl(id int64) (db string, tbl string) {
	str := strconv.FormatInt(id, 10)
	return str[:2], fmt.Sprintf("%s_%s", share_cmt_tbl, str[2:4])
}

//id生成模式:db(2位)+tbl(2位)+ts(毫秒13位)=17位
func (s *ShareComments) getId(sid int64, t time.Time) int64 {
	str := strconv.FormatInt(sid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	pfx := strconv.FormatUint(uint64(hs), 10)
	ts := utils.TimeMillisecond(t)
	id, _ := strconv.ParseInt(fmt.Sprintf("%s%s%d", cmt_db_mark, pfx[:2], ts), 10, 64)
	return id
}

func (s *ShareComments) Create(sc *ShareComment) (int64, error) {
	if len(sc.Content) == 0 {
		return 0, fmt.Errorf("评论内容不能为空")
	}
	if sc.Sid <= 0 {
		return 0, fmt.Errorf("评论的分享不存在或发布者不存在")
	}
	if sc.Uid <= 0 {
		return 0, fmt.Errorf("回复人未指定")
	}

	ss := &Shares{}
	sobj := ss.Get(sc.Sid)
	if sobj == nil {
		return 0, fmt.Errorf("评论的分享对象不存在")
	}

	now := time.Now()
	id := s.getId(sc.Sid, now)
	sc.Id = id
	sc.SUid = sobj.Uid
	sc.Ts = utils.TimeMillisecond(now)
	if sc.RUid > 0 {
		sc.T = SHARE_COMMENT_TYPE_REPLY
	} else {
		sc.T = SHARE_COMMENT_TYPE_SINGLE
	}
	pst := passport.NewMemberProvider()
	sc.SUNickname = pst.GetNicknameByUid(sc.SUid)
	sc.UNickname = pst.GetNicknameByUid(sc.Uid)
	sc.RUNickname = pst.GetNicknameByUid(sc.RUid)
	sc.Content = utils.ReplaceRepeatString(sc.Content, `\n`, 1, " ")

	//加入nosql库
	_, err := ssdb.New(use_ssdb_cmt_db).Zadd(fmt.Sprintf(SHARE_COMMENT_LIST_KEY, sc.Sid), sc.Id, sc.Ts)
	if err != nil {
		fmt.Println("publish share commet fail:", err)
		logs.Errorf("publish share commet fail:%s", err.Error())
	}
	ssdb.New(use_ssdb_cmt_db).Set(fmt.Sprintf(SHARE_COMMENT_KEY, sc.Id), *sc)

	go func() {
		o := dbs.NewOrm(share_cmt_db)
		tbl := s.hash_cmt_tbl(sc.Sid)
		_, err := o.Raw(fmt.Sprintf("insert into %s(id,sid,suid,su_nickname,uid,u_nickname,ruid,ru_nickname,t,content,ts,ex1,ex2,ex3) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)", tbl),
			sc.Id,
			sc.Sid,
			sc.SUid,
			sc.SUNickname,
			sc.Uid,
			sc.UNickname,
			sc.RUid,
			sc.RUNickname,
			sc.T,
			sc.Content,
			sc.Ts,
			sc.Ex1,
			sc.Ex2,
			sc.Ex3).Exec()
		if err != nil {
			logs.Errorf("share's comment add to mysql fail:%s", err.Error())
		}
		//更新评论数
		ss.UpdateCommentCount(sc.Sid, 1)

		for _, handler := range commentEvents {
			handler.SC_CommentedEvent(sc)
		}
		//添加参与讨论者
		go ss.JoinDiscuss(sc.Sid, []int64{sc.Uid}, sc.Ts)

	}()
	//钩子事件
	go func() {
		if sc.SUid != sc.Uid {
			hook.Do("share_qt_comment", sc.Uid, 1)
		}
	}()
	return sc.Id, nil
}

func (s *ShareComments) TopComments(sid int64, top int) []*ShareComment {
	scs := []*ShareComment{}
	ids, err := ssdb.New(use_ssdb_cmt_db).Zrevrange(fmt.Sprintf(SHARE_COMMENT_LIST_KEY, sid), 0, top-1, reflect.TypeOf(int64(0)))
	if err != nil {
		logs.Errorf("share's comments get tops fail:%v", err)
		return scs
	}
	if len(ids) == 0 {
		return scs
	}
	keys := []string{}
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf(SHARE_COMMENT_KEY, *(id.(*int64))))
	}
	_scs := ssdb.New(use_ssdb_cmt_db).MultiGet(keys, reflect.TypeOf(ShareComment{}))
	for _, sc := range _scs {
		_sc, ok := sc.(*ShareComment)
		if !ok {
			continue
		}
		scs = append(scs, _sc)
	}
	return scs
}

func (s *ShareComments) GetsByTsDesc(sid int64, page int, size int, ts time.Time) []*ShareComment {
	ids, err := ssdb.New(use_ssdb_cmt_db).Zrscan(fmt.Sprintf(SHARE_COMMENT_LIST_KEY, sid), utils.TimeMillisecond(ts)-1, 0, size, reflect.TypeOf(int64(0)))
	scs := []*ShareComment{}
	if err != nil {
		logs.Errorf("share's comments get by ts_desc tops fail:%v", err)
		return scs
	}
	keys := []string{}
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf(SHARE_COMMENT_KEY, *(id.(*int64))))
	}
	if len(keys) == 0 {
		return scs
	}
	_scs := ssdb.New(use_ssdb_cmt_db).MultiGet(keys, reflect.TypeOf(ShareComment{}))
	for _, sc := range _scs {
		_sc, ok := sc.(*ShareComment)
		if !ok {
			continue
		}
		scs = append(scs, _sc)
	}
	return scs
}

func (s *ShareComments) Delete(id int64, currentUid int64) (sid int64, err error) {
	_, tbl := s.selDbTbl(id)
	o := dbs.NewOrm(share_cmt_db)
	result := s.Get(id)
	if result == nil {
		result := ShareComment{}
		o.Raw("select id,sid,suid,su_nickname,uid,u_nickname,ruid,ru_nickname,t,content,ts,ex1,ex2,ex3 from "+tbl+" where id=?", id).QueryRow(&result)
		if result.Uid <= 0 {
			return 0, fmt.Errorf("评论不存在")
		}
	}
	if result.Uid > 0 && result.Uid != currentUid {
		return 0, fmt.Errorf("只能删除自己的评论")
	}
	//删除nosql库中引用
	ssdb.New(use_ssdb_cmt_db).Zrem(fmt.Sprintf(SHARE_COMMENT_LIST_KEY, result.Sid), result.Id)

	_, err = o.Raw("delete from "+tbl+" where id=?", id).Exec()
	if err != nil {
		logs.Errorf("share's comment delete on mysql fail:%s", err.Error())
		return result.Sid, err
	}
	return result.Sid, nil
}

func (s *ShareComments) Get(id int64) *ShareComment {
	sc := &ShareComment{}
	err := ssdb.New(use_ssdb_cmt_db).Get(fmt.Sprintf(SHARE_COMMENT_KEY, id), sc)
	if err == nil {
		return sc
	}
	return nil
}

func (s *ShareComments) GetsAll(sid int64) []*ShareComment {
	col := fmt.Sprintf(SHARE_COMMENT_LIST_KEY, sid)
	ids, err := ssdb.New(use_ssdb_cmt_db).Zrange(col, 0, -1, reflect.TypeOf(int64(0)))
	scs := []*ShareComment{}
	if err != nil {
		logs.Errorf("share's comments gets fail:%v", err)
		return scs
	}
	if len(ids) == 0 {
		return scs
	}
	keys := []string{}
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf(SHARE_COMMENT_KEY, *(id.(*int64))))
	}
	_scs := ssdb.New(use_ssdb_cmt_db).MultiGet(keys, reflect.TypeOf(ShareComment{}))
	for _, sc := range _scs {
		_sc, ok := sc.(*ShareComment)
		if !ok {
			continue
		}
		scs = append(scs, _sc)
	}
	return scs
}

func (s *ShareComments) Gets(currentUid int64, sid int64) []*ShareComment {
	scs := []*ShareComment{}
	all := s.GetsAll(sid)
	//sortutil.AscByField(scs, "Ts")
	for _, sc := range all {
		if sc.T == SHARE_COMMENT_TYPE_SINGLE {
			if sc.Uid == currentUid {
				scs = append(scs, sc)
				continue
			}
			if passport.IsBothFriend(currentUid, sc.Uid) {
				scs = append(scs, sc)
				continue
			}
		} else if sc.T == SHARE_COMMENT_TYPE_REPLY {
			if sc.Uid != currentUid && sc.RUid != currentUid {
				if passport.IsBothFriend(currentUid, sc.Uid) && passport.IsBothFriend(currentUid, sc.RUid) {
					scs = append(scs, sc)
					continue
				}
			}
			if sc.Uid == currentUid && passport.IsBothFriend(currentUid, sc.RUid) {
				scs = append(scs, sc)
				continue
			}
			if sc.RUid == currentUid && passport.IsBothFriend(currentUid, sc.Uid) {
				scs = append(scs, sc)
				continue
			}
		}
	}
	return scs
}
