package share

import (
	"dbs"
	"fmt"
	"hash/crc32"
	"libs/hook"
	"libs/message"
	"libs/passport"
	"libs/stat"
	"libs/vars"
	"logs"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"
	"utils/redis"
	"utils/ssdb"
)

var shareTbls map[string]bool = make(map[string]bool)
var s_locker *sync.Mutex = new(sync.Mutex)

const (
	SHARE_CACHE_KEY_FMT         = "share.obj:%d"
	SHARE_CACHE_DELTAG          = "share.deleted:%d"
	SHARE_MEMBER_LIST_FMT       = "share.member.list:%d" //自己发送的列表
	SHARE_DISCUSSERS            = "share_discussers:%d"  //参与者s
	share_table_max_rows  int64 = 100000000
)

type ShareSubscriptionEvent interface {
	NotifyFans(s *Share)
	RecoverFans(s *Share)
}

type ShareSaveEvent interface {
	SSE_Save(s *Share)
}

type ShareRevokeEvent interface {
	SDE_Revoke(s *Share)
}

var shareSubscriptionEvents map[string]ShareSubscriptionEvent = make(map[string]ShareSubscriptionEvent)
var shareSaveEvents map[string]ShareSaveEvent = make(map[string]ShareSaveEvent)
var shareRevokeEvents map[string]ShareRevokeEvent = make(map[string]ShareRevokeEvent)

//通知事件
func RegisterShareNotifyEvent(eventName string, eventHandler ShareSubscriptionEvent) {
	if _, ok := shareSubscriptionEvents[eventName]; ok {
		return
	}
	shareSubscriptionEvents[eventName] = eventHandler
}

//保存事件
func RegisterShareSaveEvent(eventName string, eventHandler ShareSaveEvent) {
	if _, ok := shareSaveEvents[eventName]; ok {
		return
	}
	shareSaveEvents[eventName] = eventHandler
}

//删除事件
func RegisterShareRevokeEvent(eventName string, eventHandler ShareRevokeEvent) {
	if _, ok := shareRevokeEvents[eventName]; ok {
		return
	}
	shareRevokeEvents[eventName] = eventHandler
}

type Shares struct{}

//重置新消息列表 IEventCounter interface
func (n Shares) ResetEventCount(uid int64) bool {
	return message.ResetEventCount(uid, MSG_TYPE_TEXT)
}

//IEventCounter interface
func (n Shares) NewEventCount(uid int64) int {
	return message.NewEventCount(uid, MSG_TYPE_TEXT)
}

//share库分表
func (n *Shares) Sel_tbl(id int64, pfx_table string) string {
	idstr := strconv.FormatInt(id, 10)
	if len(idstr) < 2 {
		return ""
	}
	return fmt.Sprintf("%s_%s", pfx_table, idstr[:2])
}

func (n *Shares) Hash_tbl(uid int64, pfx_table string) string {
	str := strconv.FormatInt(uid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	pfx := strconv.FormatUint(uint64(hs), 10)
	tbl := fmt.Sprintf("%s_%s", pfx_table, pfx[:2])
	if _, ok := shareTbls[tbl]; ok {
		return tbl
	}
	s_locker.Lock()
	defer s_locker.Unlock()
	if _, ok := shareTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(share_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	  uid int(11) NOT NULL,
	  st smallint(6) NOT NULL,
	  t tinyint(2) NOT NULL,
	  source varchar(32) NOT NULL,
	  geo varchar(64) NOT NULL,
	  txt varchar(300) NOT NULL,
	  create_time datetime NOT NULL,
	  ts bigint(20) NOT NULL,
	  res varchar(300) NOT NULL,
	  cmted_count int(11) NOT NULL DEFAULT '0',
	  cmt_count int(11) NOT NULL DEFAULT '0',
	  transferred_count int(11) NOT NULL DEFAULT '0',
	  transfer_count int(11) NOT NULL DEFAULT '0',
	  attituded_count int(11) NOT NULL DEFAULT '0',
	  ref_uids varchar(200) NOT NULL,
	  PRIMARY KEY (id),
	  KEY idx_uid_st_ts (uid,st,ts) USING BTREE) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		shareTbls[tbl] = true
	}
	return tbl
}

func (n *Shares) newShareId(uid int64) int64 {
	str := strconv.FormatInt(uid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	pfx := strconv.FormatUint(uint64(hs), 10)
	idstr := fmt.Sprintf("%s%d", pfx[:2], utils.TimeMillisecond(time.Now()))
	id, _ := strconv.ParseInt(idstr, 10, 64)
	return id
}

func (n *Shares) Create(s *Share, msg_notice bool) (error, int64) {
	if s.Uid <= 0 {
		return fmt.Errorf("Uid不能为0"), 0
	}

	_uidsps := strings.Split(s.RefUids, ",")
	ats := make(map[int64]int64)
	for _, _uid := range _uidsps {
		uid, err := strconv.ParseInt(_uid, 10, 64)
		if err != nil {
			continue
		}
		if uid > 0 {
			if _, ok := ats[uid]; ok {
				continue
			}
			ats[uid] = uid
		}
	}
	//@计数器
	go func() {
		for _, atuid := range ats {
			passport.IncrAtCount(s.Uid, atuid)
		}
	}()

	s.CreateTime = time.Now()
	s.Ts = utils.TimeMillisecond(s.CreateTime)

	//处理图片数据
	if s.ShareType&int(SHARE_KIND_PIC) == int(SHARE_KIND_PIC) {
		_tags := strings.Split(s.Resources, ",")
		var wait sync.WaitGroup
		for _, _tag := range _tags {
			res_tag := _tag
			wait.Add(1)
			go func() {
				defer wait.Done()
				kind, _sid := n.ResourceTranform(res_tag)
				if kind != SHARE_KIND_PIC {
					return
				}
				_picid, err := strconv.ParseInt(_sid, 10, 64)
				if err != nil {
					return
				}
				file := file_storage.GetFile(_picid)
				if file == nil {
					return
				}
				//				if file.ExtName == "gif" { //gif不支持
				//					return
				//				}
				sspic := NewShareViewPics()
				//创建预览图
				sspic.Create(file, NeedSharePicSizes)
			}()
		}
		wait.Wait()
	}

	////过滤敏感字
	s.Text = utils.CensorWords(utils.StripSQLInjection(s.Text))
	s.Text = utils.ReplaceRepeatString(s.Text, `\n`, 3, "\n\n")

	oa := dbs.NewOrm(share_db)
	tbl_a := n.Hash_tbl(s.Uid, share_tbl)
	newSid := n.newShareId(s.Uid)
	_, err := oa.Raw("insert into "+tbl_a+"(id,uid,st,t,source,geo,txt,create_time,ts,res,ref_uids) values(?,?,?,?,?,?,?,?,?,?,?)",
		newSid,
		s.Uid,
		s.ShareType,
		s.Type,
		s.Source,
		s.Geo,
		s.Text,
		s.CreateTime,
		s.Ts,
		s.Resources,
		s.RefUids).Exec()
	if err != nil {
		return fmt.Errorf("插入错误:%s", err), 0
	}
	s.Id = newSid // res.LastInsertId()

	//添加参与者
	ref_uids := s.GetRefUids()
	ref_uids = append(ref_uids, s.Uid)
	n.JoinDiscuss(s.Id, ref_uids, utils.TimeMillisecond(s.CreateTime))

	//缓存和事件
	n.saveToCacheAndMsgEvent(s, ats, msg_notice)

	go stat.GetCounter(passport.MOD_NAME).DoC(s.Uid, 1, "notes")
	if s.ShareType&int(SHARE_KIND_VOD) == int(SHARE_KIND_VOD) {
		go stat.GetCounter(passport.MOD_NAME).DoC(s.Uid, 1, "vods")
	}

	//事件钩子
	go func() {
		hook.Do("create_share", s.Uid, 1)
	}()

	return nil, s.Id
}

func (n *Shares) JoinDiscuss(sid int64, discussers []int64, ts int64) {
	if len(discussers) > 0 {
		kvs := make(map[string]int64) //uid:ts
		for _, uid := range discussers {
			uidstr := strconv.FormatInt(uid, 10)
			kvs[uidstr] = ts
		}
		ssdb.New(use_ssdb_share_db).Hmset(fmt.Sprintf(SHARE_DISCUSSERS, sid), kvs)
	}
}

func (n *Shares) GetDiscussers(sid int64) []int64 {
	arr, err := ssdb.New(use_ssdb_share_db).Hkeys(fmt.Sprintf(SHARE_DISCUSSERS, sid))
	if err != nil {
		logs.Errorf("get discussers fail:%s", err.Error())
	}
	uids := []int64{}
	for _, v := range arr {
		_uid, _ := strconv.ParseInt(v, 10, 64)
		if _uid > 0 {
			uids = append(uids, _uid)
		}
	}
	return uids
}

func (n *Shares) GetResources(res string) []*ShareResource {
	ress := strings.Split(res, ",")
	resources := []*ShareResource{}
	for _, r := range ress {
		sk, id := n.ResourceTranform(r)
		if sk != SHARE_KIND_EMPTY {
			resources = append(resources, &ShareResource{Id: id, Kind: sk})
		}
	}
	return resources
}

func (n *Shares) CombinResources(res []string) string {
	return strings.Join(res, ",")
}

func (n *Shares) UpdateCommentCount(id int64, c int) {
	s := n.Get(id)
	if s == nil {
		return
	}
	s.CommentedCount += c
	oa := dbs.NewOrm(share_db)
	tbl_a := n.Hash_tbl(s.Uid, share_tbl)
	oa.Raw(fmt.Sprintf("update "+tbl_a+" set cmted_count=cmted_count+%d where id=?", c), s.Id).Exec()
	cache := utils.GetCache()
	cache.Replace(fmt.Sprintf(SHARE_CACHE_KEY_FMT, s.Id), *s, 96*time.Hour)
}

func (n *Shares) ResourceTranform(resource string) (SHARE_KIND, string) {
	funcMaps := ShareResourceTransformFuncs()
	for _, _f := range funcMaps {
		kind, id, err := _f(resource)
		if err == nil {
			return kind, id
		}
	}

	//	if ok, _ := regexp.MatchString(share_tag_vod_regex, resource); ok {
	//		rep := regexp.MustCompile(share_tag_vod_regex)
	//		arr := rep.FindStringSubmatch(resource)
	//		return SHARE_KIND_VOD, arr[1]
	//	}
	//	if ok, _ := regexp.MatchString(share_tag_pic_regex, resource); ok {
	//		rep := regexp.MustCompile(share_tag_pic_regex)
	//		arr := rep.FindStringSubmatch(resource)
	//		return SHARE_KIND_PIC, arr[1]
	//	}
	return SHARE_KIND_EMPTY, ""
}

func (n *Shares) TranformResource(shareKind SHARE_KIND, id string, args ...string) string {
	_f := ShareTransformResourceFunc(shareKind)
	if _f == nil {
		return ""
	}
	return _f(id)
	//	switch shareKind {
	//	case SHARE_KIND_VOD:
	//		return fmt.Sprintf(share_tag_vod_fmt, id)
	//	case SHARE_KIND_PIC:
	//		return fmt.Sprintf(share_tag_pic_fmt, id)
	//	default:
	//		return ""
	//	}
}

func (n *Shares) saveToCacheAndMsgEvent(s *Share, atuids map[int64]int64, msgNotice bool) {
	cache := utils.GetCache()
	cache.Set(fmt.Sprintf(SHARE_CACHE_KEY_FMT, s.Id), *s, 96*time.Hour)

	mycol := fmt.Sprintf(SHARE_MEMBER_LIST_FMT, s.Uid)
	redis.ZAdd(nil, mycol, utils.TimeMillisecond(s.CreateTime), s.Id)
	redis.ZRemRangeByRank(nil, mycol, 0, -mbox_share_length) //限制box大小

	for _, event := range shareSaveEvents {
		go event.SSE_Save(s)
	}

	for _, event := range shareSubscriptionEvents {
		go event.NotifyFans(s)
	}

	//触发事件
	if msgNotice {
		msg_type := n.getMsgType(s)
		for _, to_uid := range atuids {
			go message.SendMsgV2(s.Uid, to_uid, msg_type, "提到了你", strconv.FormatInt(s.Id, 10), nil)
		}
	}
}

func (n *Shares) getMsgType(s *Share) vars.MSG_TYPE {
	if s.ShareType&int(SHARE_KIND_VOD) == int(SHARE_KIND_VOD) {
		return MSG_TYPE_VOD
	} else if s.ShareType&int(SHARE_KIND_PIC) == int(SHARE_KIND_PIC) {
		return MSG_TYPE_PICS
	} else {
		return MSG_TYPE_TEXT
	}
}

func (n *Shares) Delete(id int64) error {
	if id <= 0 || id < 100000000 {
		return fmt.Errorf("参数id格式错误")
	}
	s := n.Get(id)
	if s == nil {
		return fmt.Errorf("删除的分享记录不存在")
	}
	oa := dbs.NewOrm(share_db)
	tbl := n.Sel_tbl(id, share_tbl)
	res, err := oa.Raw("delete from "+tbl+" where id=?", s.Id).Exec()
	if err != nil {
		return fmt.Errorf("删除失败:%s", err)
	}
	_c, _ := res.RowsAffected()
	if _c <= 0 {
		return fmt.Errorf("删除失败:%s", err)
	}

	for _, event := range shareRevokeEvents {
		go event.SDE_Revoke(s)
	}
	for _, event := range shareSubscriptionEvents {
		go event.RecoverFans(s)
	}

	cache := utils.GetCache()
	cache.Delete(fmt.Sprintf(SHARE_CACHE_KEY_FMT, id))

	//已删除标记
	ssdb.New(use_ssdb_share_db).Set(fmt.Sprintf(SHARE_CACHE_DELTAG, id), "del")

	go stat.GetCounter(passport.MOD_NAME).DoC(s.Uid, -1, "notes")
	if s.ShareType&int(SHARE_KIND_VOD) == int(SHARE_KIND_VOD) {
		go stat.GetCounter(passport.MOD_NAME).DoC(s.Uid, -1, "vods")
	}
	return nil
}

func (n *Shares) Get(id int64) *Share {
	if id <= 0 {
		return nil
	}
	cache := utils.GetCache()
	s := &Share{}
	err := cache.Get(fmt.Sprintf(SHARE_CACHE_KEY_FMT, id), s)
	if err == nil {
		return s
	}

	//已删除标记
	exist, _ := ssdb.New(use_ssdb_share_db).Exists(fmt.Sprintf(SHARE_CACHE_DELTAG, id))
	if exist {
		return nil
	}

	o := dbs.NewOrm(share_db)
	tbl := n.Sel_tbl(id, share_tbl)
	if len(tbl) == 0 {
		return nil
	}

	err = o.Raw("select id,uid,st,t,source,geo,txt,create_time,ts,res,cmted_count,cmt_count,transferred_count,transfer_count,ref_uids from "+tbl+" where id=?", id).QueryRow(s)
	if err != nil {
		//已被删除标记
		ssdb.New(use_ssdb_share_db).Set(fmt.Sprintf(SHARE_CACHE_DELTAG, id), "del")
		return nil
	} else {
		cache.Set(fmt.Sprintf(SHARE_CACHE_KEY_FMT, id), *s, 96*time.Hour)
		return s
	}
	return nil
}

//某人发布的
func (n *Shares) Gets(uid int64, page int, size int, ts time.Time) (int, []*Share) {
	ss := []*Share{}
	if uid <= 0 {
		return 0, ss
	}
	c := stat.GetCounter(passport.MOD_NAME).GetC(uid, "notes")
	type Result struct {
		Id int64
		Ts int64
	}

	offset := page * size
	if offset <= mbox_share_length {
		col_name := fmt.Sprintf(SHARE_MEMBER_LIST_FMT, uid)
		kvs, _ := redis.ZRevRangeByScore(nil, col_name, reflect.TypeOf(int64(0)), utils.TimeMillisecond(ts)-1, 0, "LIMIT", 0, size)

		//redis数据可能被清空,重新从db中获取数据后导入
		if c > 0 && len(kvs) == 0 {
			o := dbs.NewOrm(share_db)
			tbl := n.Hash_tbl(uid, share_tbl)
			sql := fmt.Sprintf("select id,ts from %s where uid=? and ts <= ? order by ts desc limit %d", tbl, mbox_share_length)
			var res []Result
			o.Raw(sql, uid, utils.TimeMillisecond(time.Now())).QueryRows(&res)
			if len(res) > 0 {
				ids := []interface{}{}
				for _, r := range res {
					ids = append(ids, r.Ts)
					ids = append(ids, r.Id)
				}
				redis.ZMultiAdd(nil, col_name, ids...)
				kvs, _ = redis.ZRevRangeByScore(nil, col_name, reflect.TypeOf(int64(0)), utils.TimeMillisecond(ts)-1, 0, "LIMIT", 0, size)
			}
		}

		for _, key := range kvs {
			if id, ok := key.(*int64); ok {
				_s := n.Get(*id)
				if _s != nil {
					ss = append(ss, _s)
				}
			}
		}
	} else {
		o := dbs.NewOrm(share_db)
		tbl := n.Hash_tbl(uid, share_tbl)
		sql := fmt.Sprintf("select id,ts from %s where uid=? and ts <= ? order by ts desc limit %d", tbl, size)
		var res []Result
		o.Raw(sql, uid, utils.TimeMillisecond(ts)-1).QueryRows(&res)
		for _, r := range res {
			_s := n.Get(r.Id)
			if _s != nil {
				ss = append(ss, _s)
			}
		}
	}
	return c, ss
}

func (n *Shares) ShareOutside(uid int64) {
	if uid <= 0 {
		return
	}
	hook.Do("share_weixin", uid, 1)
}
