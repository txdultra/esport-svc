package collect

import (
	"dbs"
	"fmt"
	"strconv"
	"time"
	"utils/ssdb"

	"logs"

	"labix.org/v2/mgo/bson"
)

type MemberCollect struct {
	RelId   string
	RelType string
}

func member_collect_box(uid int64) string {
	return fmt.Sprintf("member_collect_box:%d", uid)
}

func Create(uid int64, relId string, relType string) (error, string) {
	if len(relType) == 0 {
		return fmt.Errorf("RelType don't empty"), ""
	}
	compler := getCompler(relType)
	if compler == nil {
		return fmt.Errorf(relType + " provider not exist"), ""
	}
	collect := Collectible{
		Id:         bson.NewObjectId(),
		Uid:        uid,
		RelId:      relId,
		RelType:    relType,
		CreateTime: time.Now(),
	}
	err := compler.CompleCollectible(&collect)
	if err != nil {
		logs.Error(err.Error())
		return err, ""
	}
	mc := MemberCollect{
		RelId:   relId,
		RelType: relType,
	}
	session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	defer session.Close()
	var fav *Collectible
	err = col.Find(bson.M{"uid": uid, "rel_id": relId, "rel_type": relType}).One(&fav)
	if err == nil {
		go func() {
			existed, _ := ssdb.New(use_ssdb_collect_db).Zexists(member_collect_box(uid), mc)
			if !existed {
				ssdb.New(use_ssdb_collect_db).Zadd(member_collect_box(uid), mc, collect.CreateTime.Unix())
			}
		}()
		return nil, fav.Id.Hex()
	}
	err = col.Insert(collect)
	if err != nil {
		logs.Error(err.Error())
		return fmt.Errorf("插入错误:%s", err), ""
	}
	go ssdb.New(use_ssdb_collect_db).Zadd(member_collect_box(uid), mc, collect.CreateTime.Unix())
	return nil, collect.Id.Hex()
}

func IsCollcetd(uid int64, relId string, relType string) (bool, error) {
	mc := MemberCollect{
		RelId:   relId,
		RelType: relType,
	}
	exist, err := ssdb.New(use_ssdb_collect_db).Zexists(member_collect_box(uid), mc)
	return exist, err
}

func Get(uid int64, relId string, relType string) *Collectible {
	session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	defer session.Close()
	var fav *Collectible
	err := col.Find(bson.M{"uid": uid, "rel_id": relId, "rel_type": relType}).One(&fav)
	if err == nil {
		return fav
	}
	return nil
}

func Delete(uid int64, id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("id格式错误")
	}
	session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	defer session.Close()
	var fav *Collectible
	err := col.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&fav)
	if err != nil {
		return fmt.Errorf("删除失败:COLLECT_DM_001")
	}
	err = col.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return fmt.Errorf("删除失败:COLLECT_DM_002")
	}
	mc := MemberCollect{
		RelId:   fav.RelId,
		RelType: fav.RelType,
	}
	ssdb.New(use_ssdb_collect_db).Zrem(member_collect_box(uid), mc)
	return nil
}

func Deletes(uid int64, ids []string) {
	session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	defer session.Close()
	for _, id := range ids {
		if len(id) == 0 || !bson.IsObjectIdHex(id) {
			continue
		}
		var fav *Collectible
		err := col.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&fav)
		if err != nil {
			continue
		}
		mc := MemberCollect{
			RelId:   fav.RelId,
			RelType: fav.RelType,
		}
		col.RemoveId(bson.ObjectIdHex(id))
		ssdb.New(use_ssdb_collect_db).Zrem(member_collect_box(uid), mc)
	}
}

func DeleteAll(relId, relType string) error {
	alldbs := alldb_cols(collect_collection)
	for _, col := range alldbs {
		session, col := dbs.MgoC(collect_db, col)
		col.RemoveAll(bson.M{"rel_id": relId, "rel_type": relType})
		session.Close()
	}
	//session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	//defer session.Close()
	//_, err := col.RemoveAll(bson.M{"rel_id": relId, "rel_type": relType})
	//if err != nil {
	//	return err
	//}
	return nil
}

func Gets(uid int64, p int, s int, last_time time.Time) (int, []*Collectible) {
	session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	defer session.Close()
	var collects []*Collectible
	fm := bson.M{}
	fm["uid"] = uid
	qs := col.Find(fm)
	counts, _ := qs.Count()
	fm["create_time"] = bson.M{"$lt": last_time}
	err := col.Find(fm).Sort("-create_time").Limit(s).All(&collects)
	if err != nil {
		print(err.Error())
	}
	return counts, collects
}

func GetsAll(relType string, uid int64, p int, s int, last_time time.Time) (int, []*Collectible) {
	session, col := dbs.MgoC(collect_db, sdb_col(uid, collect_collection))
	defer session.Close()
	var collects []*Collectible
	fm := bson.M{}
	fm["uid"] = uid
	fm["rel_type"] = relType
	qs := col.Find(fm)
	counts, _ := qs.Count()
	fm["create_time"] = bson.M{"$lt": last_time}
	err := col.Find(fm).Sort("-create_time").Limit(s).All(&collects)
	if err != nil {
		print(err.Error())
	}
	return counts, collects
}

var allpfxs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

func alldb_cols(colname string) []string {
	all := make([]string, 9)
	for i, j := range allpfxs {
		all[i] = fmt.Sprintf("%s_%d", colname, j)
	}
	return all
}

func sdb_col(uid int64, colname string) string {
	str := strconv.FormatInt(uid, 10)
	return fmt.Sprintf("%s_%s", colname, string(str[0]))
}
