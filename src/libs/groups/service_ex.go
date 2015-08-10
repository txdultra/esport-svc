package groups

import (
	"dbs"
	"fmt"
	"reflect"
	"time"
	"utils"
	"utils/ssdb"

	"github.com/astaxie/beego/orm"
)

const (
	//user_post_table_fmt      = "user_post_note_%d"
	//user_post_tid_existed    = "user_post_tid_existed_%d_%d"
	user_post_count_noted    = "group_post_count_noted_%d"
	user_post_count_cachekey = "group_user_post_count_%d"
	user_pthread_cachekey    = "group_user_pthread_%d"
	user_ppost_cachekey      = "group_user_ppost_%d"
)

func NewUserPostNoteService() *UserPostNotes {
	return &UserPostNotes{}
}

var userPostNoteTbls map[int64]string = make(map[int64]string)

type UserPostNotes struct{}

func (u *UserPostNotes) AddNote(ptr interface{}) {
	publishsAdds := 0
	replysAdds := 0
	var uid int64
	cclient := ssdb.New(use_ssdb_group_db)
	thread, ok_t := ptr.(*Thread)
	if ok_t {
		cclient.Zadd(fmt.Sprintf(user_pthread_cachekey, thread.AuthorId), thread.Id, thread.DateLine)
		publishsAdds = 1
		uid = thread.AuthorId
	}
	post, ok_p := ptr.(*Post)
	if ok_p {
		cclient.Zadd(fmt.Sprintf(user_ppost_cachekey, post.AuthorId), post.ThreadId, post.DateLine)
		replysAdds = 1
		uid = post.AuthorId
	}
	if !ok_t && !ok_p {
		return
	}
	cache := utils.GetCache()
	user_noted_key := fmt.Sprintf(user_post_count_cachekey, uid)
	exist, _ := cclient.Exists(user_noted_key)
	if !exist {
		o := dbs.NewOrm(db_aliasname)
		exist = o.QueryTable(&UserPostNoteCount{}).Filter("uid", uid).Exist()
	}
	if !exist {
		upc := &UserPostNoteCount{
			Uid:      uid,
			Publishs: publishsAdds,
			Replys:   replysAdds,
			LastTime: time.Now().Unix(),
		}
		o := dbs.NewOrm(db_aliasname)
		o.Insert(upc)
		cache.Set(fmt.Sprintf(user_post_count_cachekey, uid), *upc, 24*time.Hour)
		//已记录入数据库
		cclient.Set(fmt.Sprintf(user_post_count_noted, uid), 1)
	} else {
		ts := time.Now().Unix()
		o := dbs.NewOrm(db_aliasname)
		o.QueryTable(&UserPostNoteCount{}).Filter("uid", uid).Update(orm.Params{
			"publishs": orm.ColValue(orm.Col_Add, publishsAdds),
			"replys":   orm.ColValue(orm.Col_Add, replysAdds),
			"lasttime": ts,
		})
		upc := &UserPostNoteCount{}
		err := cache.Get(fmt.Sprintf(user_post_count_cachekey, uid), upc)
		if err == nil {
			upc.Publishs += publishsAdds
			upc.Replys += replysAdds
			upc.LastTime = ts
			cache.Replace(fmt.Sprintf(user_post_count_cachekey, uid), *upc, 24*time.Hour)
		}
	}
}

//func (u *UserPostNotes) getTableId(uid int64) string {
//	_r := uid % 100
//	if name, ok := userPostNoteTbls[_r]; ok {
//		return name
//	}
//	tbl_mutex.Lock()
//	defer tbl_mutex.Unlock()
//	if name, ok := userPostNoteTbls[_r]; ok {
//		return name
//	}
//	tbl := fmt.Sprintf(user_post_table_fmt, _r)

//	o := dbs.NewOrm(db_aliasname)
//	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
//	  id int(11) unsigned NOT NULL AUTO_INCREMENT,
//	  uid int(11) NOT NULL,
//	  tid int(11) NOT NULL,
//	  postid char(20) COLLATE utf8mb4_bin NOT NULL,
//	  subject char(30) COLLATE utf8mb4_bin NOT NULL,
//	  dateline int(11) NOT NULL,
//	  t tinyint(2) NOT NULL,
//      PRIMARY KEY (id),
//	  KEY idx_uid_time_t (uid,dateline,t) USING BTREE
//	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, tbl)
//	o.Raw(create_tbl_sql).Exec()
//	userPostNoteTbls[_r] = tbl

//	savePostTableToDb(int(_r), tbl, &UserPostNoteTable{
//		Id:      _r,
//		TblName: tbl,
//		Ts:      time.Now().Unix(),
//	})
//	return tbl
//}

//func (u *UserPostNotes) AddNote(ptr interface{}, subject string) error {
//	var up *UserPostNote
//	publishsAdds := 0
//	replysAdds := 0
//	thread, ok := ptr.(*Thread)
//	if ok {
//		up = &UserPostNote{
//			T:        USER_POST_TYPE_THREAD,
//			ThreadId: thread.Id,
//			Uid:      thread.AuthorId,
//			PostId:   thread.Lordpid,
//			Subject:  subject,
//		}
//		publishsAdds = 1
//	}
//	post, ok := ptr.(*Post)
//	if ok {
//		up = &UserPostNote{
//			T:        USER_POST_TYPE_POST,
//			ThreadId: post.ThreadId,
//			Uid:      post.AuthorId,
//			PostId:   post.Id,
//			Subject:  subject,
//		}
//		replysAdds = 1
//	}
//	if up != nil {
//		up.DateLine = time.Now().Unix()
//		u.AddOrUpdate(up)
//		//更新计数
//		o := dbs.NewOrm(db_aliasname)
//		exist := o.QueryTable(&UserPostNoteCount{}).Filter("uid", up.Uid).Exist()
//		if !exist {
//			o.Insert(&UserPostNoteCount{
//				Uid:      up.Uid,
//				Publishs: publishsAdds,
//				Replys:   replysAdds,
//				LastTime: time.Now().Unix(),
//			})
//		} else {
//			ts := time.Now().Unix()
//			o.QueryTable(&UserPostNoteCount{}).Filter("uid", up.Uid).Update(orm.Params{
//				"publishs": orm.ColValue(orm.Col_Add, publishsAdds),
//				"replys":   orm.ColValue(orm.Col_Add, replysAdds),
//				"lasttime": ts,
//			})
//			cache := utils.GetCache()
//			upc := &UserPostNoteCount{}
//			err := cache.Get(fmt.Sprintf(user_post_count_cachekey, up.Uid), upc)
//			if err == nil {
//				upc.Publishs += publishsAdds
//				upc.Replys += replysAdds
//				upc.LastTime = ts
//				cache.Replace(fmt.Sprintf(user_post_count_cachekey, up.Uid), *upc, 24*time.Hour)
//			}
//		}
//	}
//	return nil
//}

//func (u *UserPostNotes) AddOrUpdate(up *UserPostNote) {
//	tbl := u.getTableId(up.Uid)
//	cclient := ssdb.New(use_ssdb_group_db)
//	existkey := fmt.Sprintf(user_post_tid_existed, up.ThreadId, up.T)
//	has, _ := cclient.Exists(existkey)
//	if has {
//		sql := fmt.Sprintf("update %s set dateline=? where uid=? and tid=? and t=?", tbl)
//		o := dbs.NewOrm(db_aliasname)
//		o.Raw(sql, up.DateLine, up.Uid, up.ThreadId, up.T).Exec()
//	} else {
//		o := dbs.NewOrm(db_aliasname)
//		sql := fmt.Sprintf("select id from %s where uid=? and tid=? and t=?", tbl)
//		var list orm.ParamsList
//		num, _ := o.Raw(sql).ValuesFlat(&list)
//		if num == 0 {
//			sql := fmt.Sprintf("insert into %s(uid,tid,postid,subject,dateline,t) values(?,?,?,?,?,?)", tbl)
//			o.Raw(sql, up.Uid, up.ThreadId, up.PostId, up.Subject, up.DateLine, up.T).Exec()
//		} else {
//			sql := fmt.Sprintf("update %s set dateline=? where uid=? and tid=? and t=?", tbl)
//			o.Raw(sql, up.DateLine, up.Uid, up.ThreadId, up.T).Exec()
//		}
//	}
//}

func (u *UserPostNotes) GetNoteCount(uid int64) *UserPostNoteCount {
	cache := utils.GetCache()
	ckey := fmt.Sprintf(user_post_count_cachekey, uid)
	upc := &UserPostNoteCount{}
	err := cache.Get(ckey, upc)
	if err == nil {
		return upc
	}
	o := dbs.NewOrm(db_aliasname)
	err = o.QueryTable(upc).Filter("uid", uid).One(upc)
	if err == nil {
		cache.Set(ckey, *upc, 24*time.Hour)
		return upc
	}
	upc.Uid = uid
	cache.Set(ckey, *upc, 2*time.Hour)
	return nil
}

func (u *UserPostNotes) Gets(uid int64, t USER_POST_TYPE, ts int64, size int) (tids []int64, nextTs int64) {
	//	tbl := u.getTableId(uid)
	//	o := dbs.NewOrm(db_aliasname)
	//	offset := (page - 1) * size
	//	sql := fmt.Sprintf("select id,uid,tid,postid,subject,dateline,t from %s where uid=? and t = ? order by dateline desc limit %d,%d", tbl, offset, size)
	//	var ups []*UserPostNote
	//	o.Raw(sql, uid, t).QueryRows(&ups)
	//	return ups
	ckey := ""
	switch t {
	case USER_POST_TYPE_POST:
		ckey = fmt.Sprintf(user_ppost_cachekey, uid)
		break
	case USER_POST_TYPE_THREAD:
		ckey = fmt.Sprintf(user_pthread_cachekey, uid)
		break
	}
	cclient := ssdb.New(use_ssdb_group_db)
	kss, err := cclient.ZrscanKS(ckey, ts-1, 0, size, reflect.TypeOf(int64(0)), reflect.TypeOf(int64(0)))
	fmt.Println(kss)
	tids = []int64{}
	if err != nil {
		return tids, ts
	}
	for _, ks := range kss {
		tid := *(ks.Key.(*int64))
		nextTs = *(ks.Score.(*int64))
		tids = append(tids, tid)
	}
	return
}
