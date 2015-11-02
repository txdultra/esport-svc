package passport

import (
	"dbs"
	"fmt"
	"hash/crc32"
	"libs/hook"
	"libs/search"
	"libs/stat"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"
	"utils/ssdb"
	//"github.com/astaxie/beego/orm"
)

const (
	FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME = "friendship_new_follower_count_mod"
)

//删除关注事件处理
type UnFriendEvent interface {
	FSUnFriendDo(suid int64, tuid int64)
}

//添加关注事件处理
type FriendEvent interface {
	FSFriendDo(suid int64, tuid int64)
}

var friendEvens map[string]FriendEvent = make(map[string]FriendEvent)
var unFriendEvents map[string]UnFriendEvent = make(map[string]UnFriendEvent)
var friendTbls map[string]bool = make(map[string]bool)
var locker *sync.Mutex = new(sync.Mutex)
var memberProvier *MemberProvider = NewMemberProvider()

//注册关注事件
func RegisterFriendEvent(eventName string, event FriendEvent) {
	if _, ok := friendEvens[eventName]; !ok {
		friendEvens[eventName] = event
	}
}

//注册取消关注事件
func RegisterUnFriendEvent(eventName string, event UnFriendEvent) {
	if _, ok := unFriendEvents[eventName]; !ok {
		unFriendEvents[eventName] = event
	}
}

type RELATION_TYPE int

const (
	RELATION_TYPE_FRIEND RELATION_TYPE = 1
	RELATION_TYPE_FANS   RELATION_TYPE = 2
)

const (
	friendship_new_followers_sets = "friendship_new_followers_uid_v8:%d" //新粉丝
	friendship_friend_sets        = "friendship_friends_uid_v8:%d"       //关注列表
	friendship_follower_sets      = "friendship_follower_uid_v8:%d"      //粉丝列表
	friendship_friend_pylist      = "friendship_friends_pylist_uid_v8:%d"
	friendship_bothfriend_pylist  = "friendship_bothfriends_pylist_uid_v8:%d"
	friendship_friend_both_sets   = "friendship_friends_both_uid_v8:%d" //互粉
)

func MemberFriends(uid int64) int {
	count, _ := ssdb.New(use_ssdb_passport_db).Zcard(fmt.Sprintf(friendship_friend_sets, uid))
	return count
}

func MemberFollowers(uid int64) int {
	count, _ := ssdb.New(use_ssdb_passport_db).Zcard(fmt.Sprintf(friendship_follower_sets, uid))
	return count
}

//是否已关注某人
func IsFriend(source_uid int64, target_uid int64) bool {
	has_ns, _ := ssdb.New(use_ssdb_passport_db).Zexists(fmt.Sprintf(friendship_friend_sets, source_uid), target_uid)
	return has_ns
}

//是否互粉
func IsBothFriend(source_uid int64, target_uid int64) bool {
	//只需判断一个源里是否存在目标对象，就可断定是互粉状态
	has_ns, _ := ssdb.New(use_ssdb_passport_db).Zexists(fmt.Sprintf(friendship_friend_both_sets, source_uid), target_uid)
	return has_ns
}

type FriendShips struct{}

func (m *FriendShips) hash_tbl(uid int64, table string) string {
	str := strconv.FormatInt(uid, 10)
	hs := crc32.ChecksumIEEE([]byte(str))
	pfx := strconv.FormatUint(uint64(hs), 10)
	tbl := fmt.Sprintf("%s_%s", table, pfx[:2])
	if _, ok := friendTbls[tbl]; ok {
		return tbl
	}
	locker.Lock()
	defer locker.Unlock()
	if _, ok := friendTbls[tbl]; ok {
		return tbl
	}
	orm := dbs.NewOrm(relation_db)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(id int(11) unsigned NOT NULL AUTO_INCREMENT,s_uid int(11) NOT NULL,
									t_uid int(11) NOT NULL,post_time datetime NOT NULL,nano bigint(20) NOT NULL,PRIMARY KEY (id),rel_t tinyint(1) NOT NULL,
									UNIQUE KEY from_to_unique (s_uid,t_uid)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8`, tbl)
	_, err := orm.Raw(create_tbl_sql).Exec()
	if err == nil {
		friendTbls[tbl] = true
	}
	return tbl
}

func (m *FriendShips) FriendTo(source_uid int64, target_uid int64) error {
	if source_uid <= 0 || target_uid <= 0 {
		return fmt.Errorf("参数错误")
	}
	if source_uid == target_uid {
		return fmt.Errorf("不能关注自己")
	}
	mfs := MemberFriends(source_uid)
	if mfs >= friend_limit_counts {
		return fmt.Errorf("关注人数已超过限制")
	}

	if IsFriend(source_uid, target_uid) {
		return fmt.Errorf("已在关注列表中")
	}

	now := time.Now()
	ok, _ := ssdb.New(use_ssdb_passport_db).Zadd(fmt.Sprintf(friendship_friend_sets, source_uid), target_uid, utils.TimeMillisecond(now)) //添加入好友Set
	if !ok {
		return fmt.Errorf("关注失败")
	}
	ok, _ = ssdb.New(use_ssdb_passport_db).Zadd(fmt.Sprintf(friendship_follower_sets, target_uid), source_uid, utils.TimeMillisecond(now)) //添加入粉丝Set
	if !ok {
		ssdb.New(use_ssdb_passport_db).Zrem(fmt.Sprintf(friendship_friend_sets, source_uid), target_uid) //回滚删除好友Set
		return fmt.Errorf("关注失败")
	}
	//新粉丝
	ssdb.New(use_ssdb_passport_db).Zadd(fmt.Sprintf(friendship_new_followers_sets, target_uid), source_uid, utils.TimeMillisecond(now))
	go stat.UCIncrCount(target_uid, FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME)
	//互粉
	if IsFriend(target_uid, source_uid) {
		ssdb.New(use_ssdb_passport_db).Zadd(fmt.Sprintf(friendship_friend_both_sets, source_uid), target_uid, utils.TimeMillisecond(now))
		ssdb.New(use_ssdb_passport_db).Zadd(fmt.Sprintf(friendship_friend_both_sets, target_uid), source_uid, utils.TimeMillisecond(now))
	}
	//数据库更新
	stat.GetCounter(MOD_NAME).DoC(source_uid, 1, "friends")
	stat.GetCounter(MOD_NAME).DoC(target_uid, 1, "fans")

	//删除py关注列表缓存
	cache := utils.GetCache()
	cache.Delete(fmt.Sprintf(friendship_friend_pylist, source_uid))
	cache.Delete(fmt.Sprintf(friendship_bothfriend_pylist, source_uid))

	go func() {
		o := dbs.NewOrm(relation_db)
		tbl_s := m.hash_tbl(source_uid, relation_tbl)
		o.Raw("insert into "+tbl_s+"(s_uid,t_uid,post_time,nano,rel_t) values(?,?,?,?,?)", source_uid, target_uid, now, utils.TimeMillisecond(now), RELATION_TYPE_FRIEND).Exec()
		tbl_t := m.hash_tbl(target_uid, relation_tbl)
		o.Raw("insert into "+tbl_t+"(s_uid,t_uid,post_time,nano,rel_t) values(?,?,?,?,?)", target_uid, source_uid, now, utils.TimeMillisecond(now), RELATION_TYPE_FANS).Exec()
	}()

	//事件处理
	go func() {
		for _, event := range friendEvens {
			event.FSFriendDo(source_uid, target_uid)
		}
	}()

	//钩子事件
	go func() {
		friends := m.FriendCounts(source_uid)
		fans := m.FollowerCounts(target_uid)
		//		sBothFriends := m.BothFriendCounts(source_uid)
		//		tBothFriends := m.BothFriendCounts(target_uid)
		hook.Do("gz_count", source_uid, friends)
		hook.Do("fans_count", target_uid, fans)
		//		hook.Do("friend_count", source_uid, sBothFriends)
		//		hook.Do("friend_count", target_uid, tBothFriends)
		tMember := memberProvier.Get(target_uid)
		if tMember != nil && tMember.Certified {
			hook.Do("subscr_vuser", source_uid, 1)
		}
	}()
	return nil
}

func (m *FriendShips) DestroyFriend(source_uid int64, target_uid int64) error {
	if source_uid <= 0 || target_uid <= 0 {
		return fmt.Errorf("参数错误")
	}
	if source_uid == target_uid {
		return fmt.Errorf("不能取消关注自己")
	}

	if !IsFriend(source_uid, target_uid) {
		return fmt.Errorf("未在关注列表中")
	}

	ssdb.New(use_ssdb_passport_db).Zrem(fmt.Sprintf(friendship_friend_sets, source_uid), target_uid)   //删除好友Set
	ssdb.New(use_ssdb_passport_db).Zrem(fmt.Sprintf(friendship_follower_sets, target_uid), source_uid) //删除粉丝Set

	//删除新粉丝
	ssdb.New(use_ssdb_passport_db).Zrem(fmt.Sprintf(friendship_new_followers_sets, target_uid), source_uid)
	go stat.UCDecrCount(target_uid, FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME)

	//互粉
	ssdb.New(use_ssdb_passport_db).Zrem(fmt.Sprintf(friendship_friend_both_sets, source_uid), target_uid)
	ssdb.New(use_ssdb_passport_db).Zrem(fmt.Sprintf(friendship_friend_both_sets, target_uid), source_uid)

	//数据库更新
	stat.GetCounter(MOD_NAME).DoC(source_uid, -1, "friends")
	stat.GetCounter(MOD_NAME).DoC(target_uid, -1, "fans")

	//删除py关注列表缓存
	cache := utils.GetCache()
	cache.Delete(fmt.Sprintf(friendship_friend_pylist, source_uid))
	cache.Delete(fmt.Sprintf(friendship_bothfriend_pylist, source_uid))

	go func() {
		o := dbs.NewOrm(relation_db)
		tbl_s := m.hash_tbl(source_uid, relation_tbl)
		o.Raw("delete from "+tbl_s+" where s_uid=? and t_uid=?", source_uid, target_uid).Exec()
		tbl_t := m.hash_tbl(target_uid, relation_tbl)
		o.Raw("delete from "+tbl_t+" where s_uid=? and t_uid=?", target_uid, source_uid).Exec()
	}()
	//事件处理
	go func() {
		for _, event := range unFriendEvents {
			event.FSUnFriendDo(source_uid, target_uid)
		}
	}()
	return nil
}

func (m *FriendShips) NewEventCount(uid int64) int {
	//	key := fmt.Sprintf(friendship_new_followers_sets, uid)
	//	c, err := ssdb.New(use_ssdb_passport_db).Zcard(key)
	//	if err != nil {
	//		return 0
	//	}
	c := stat.UCGetCount(uid, FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME)
	return int(c)
}

//IEventCounter interface
func (m *FriendShips) ResetEventCount(uid int64) bool {
	key := fmt.Sprintf(friendship_new_followers_sets, uid)
	err := ssdb.New(use_ssdb_passport_db).Zclear(key)
	stat.UCResetCount(uid, FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME)
	return (err == nil)
}

func (m *FriendShips) RepairRelation(uid int64) {

}

//首字母拼音分组列表
func (m *FriendShips) FriendIdsByNickNameFirstPY(uid int64) map[string][]int64 {
	cache := utils.GetCache()
	ckey := fmt.Sprintf(friendship_friend_pylist, uid)
	var obj map[string][]int64
	err := cache.Get(ckey, &obj)
	if err == nil {
		return obj
	}

	uids := m.FriendIds(uid)
	pymap := m.py(uids)
	cache.Set(ckey, pymap, 48*time.Hour)
	return pymap
}

func (m *FriendShips) py(uids []int64) map[string][]int64 {
	pm := NewMemberProvider()
	idnks := make(map[int64]string)
	for _, _uid := range uids {
		nickName := pm.GetNicknameByUid(_uid)
		if len(nickName) == 0 {
			continue
		}
		idnks[_uid] = nickName
	}
	pymap := make(map[string][]int64)
	toF := func(char string) string {
		if len(char) == 0 {
			return "*"
		}
		if char[0] < 65 || char[0] > 90 {
			return "*"
		}
		return char
	}
	for k, v := range idnks {
		py := utils.PYConvert(v, 1)
		py = strings.ToUpper(py)
		fpy := ""
		if len(py) > 0 {
			fpy = py[:1]
		}
		fpy = toF(fpy)
		if ids, ok := pymap[fpy]; ok {
			ids = append(ids, k)
			pymap[fpy] = ids
		} else {
			ids := []int64{k}
			pymap[fpy] = ids
		}
	}
	return pymap
}

//推存关注名单(根据关注游戏和主播粉丝数)
//func (m *FriendShips) RecmdFriendUids(uid int64, tops int) []int64 {
//	ms := NewMemberProvider()
//	mgs := ms.MemberGames(uid)
//	intts := []string{}
//	for _, mg := range mgs {
//		intts = append(intts, strconv.Itoa(mg.GameId))
//	}
//	gameids := strings.Join(intts, ",")
//	ckey := fmt.Sprintf("friendship_recmd_gameids_%s_tops_%d", gameids, tops)
//	var uids []int64
//	cache := utils.GetCache()
//	err := cache.Get(ckey, &uids)
//	if err == nil {
//		return uids
//	}
//	sql := fmt.Sprintf("SELECT DISTINCT(cmg.uid) FROM common_member_games cmg join common_member_states cms on cmg.uid = cms.uid where cmg.gid in (%s) and cmg.uid <> %d order by cms.fans DESC limit %d;", gameids, uid, tops)
//	var maps []orm.Params
//	uids = []int64{}
//	o := dbs.NewDefaultOrm()
//	num, err := o.Raw(sql).Values(&maps)
//	if err == nil && num > 0 {
//		fmt.Println(maps)
//		for _, row := range maps {
//			_uid := row["uid"].(string)
//			uid, _ := strconv.ParseInt(_uid, 10, 64)
//			uids = append(uids, uid)
//		}
//	}
//	cache.Set(ckey, uids, 3*time.Hour)
//	return uids
//}

//推存关注名单(根据关注游戏和主播粉丝数,随机)
func (m *FriendShips) RecmdFriendUids(uid int64, tops int, certified bool, showFriended bool) []int64 {
	ms := NewMemberProvider()
	mgs := ms.MemberGames(uid)
	uids := []int64{}

	//搜索配置
	fgameIds := []uint64{}
	for _, mg := range mgs {
		fgameIds = append(fgameIds, uint64(mg.GameId))
	}
	var filters []search.SearchFilter
	if len(fgameIds) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "gids",
			Values:  fgameIds,
			Exclude: false,
		})
	}
	if certified {
		filters = append(filters, search.SearchFilter{
			Attr:    "certified",
			Values:  []uint64{1},
			Exclude: false,
		})
	}

	if !showFriended {
		ffsuids := []uint64{}
		fsuids := m.FriendIds(uid)
		for _, fuid := range fsuids {
			ffsuids = append(ffsuids, uint64(fuid))
		}
		filters = append(filters, search.SearchFilter{
			Attr:    "uid",
			Values:  ffsuids,
			Exclude: true,
		})
	}
	_, uids = ms.Query("", 1, tops, "extended", "@random", filters, nil)
	return uids
}

//关注列表(全部)
func (m *FriendShips) FriendIds(uid int64) []int64 {
	ids := []int64{}
	lst, err := ssdb.New(use_ssdb_passport_db).Zrevrange(fmt.Sprintf(friendship_friend_sets, uid), 0, -1, reflect.TypeOf(int64(0)))
	if err != nil {
		return ids
	}
	for _, id := range lst {
		_id := *(id.(*int64))
		ids = append(ids, _id)
	}
	return ids
}

func (m *FriendShips) FriendCounts(uid int64) int {
	i, _ := ssdb.New(use_ssdb_passport_db).Zcard(fmt.Sprintf(friendship_friend_sets, uid))
	return i
}

//关注列表(分页型)
func (m *FriendShips) FriendIdsP(uid int64, size int, page int, ts time.Time) (int, []int64, int64) {
	friends_col := fmt.Sprintf(friendship_friend_sets, uid)
	total, _ := ssdb.New(use_ssdb_passport_db).Zcard(friends_col)
	uids, _ := ssdb.New(use_ssdb_passport_db).Zrscan(friends_col, utils.TimeMillisecond(ts)-1, 0, size, reflect.TypeOf(int64(0)))
	if len(uids) == 0 {
		return 0, []int64{}, utils.TimeMillisecond(ts)
	}
	ids := []int64{}
	for _, _uid := range uids {
		_id := *(_uid.(*int64))
		ids = append(ids, _id)
	}
	last_t := utils.TimeMillisecond(ts)
	if len(ids) > 0 {
		_last_t, _ := ssdb.New(use_ssdb_passport_db).Zscore(friends_col, ids[len(ids)-1])
		last_t = _last_t
	}
	return total, ids, last_t
}

//粉丝列表(分页型)
func (m *FriendShips) FollowerIdsP(uid int64, size int, page int, ts time.Time) (int, []int64, int64) {

	follower_col := fmt.Sprintf(friendship_follower_sets, uid)
	total, _ := ssdb.New(use_ssdb_passport_db).Zcard(follower_col)
	uids, _ := ssdb.New(use_ssdb_passport_db).Zrscan(follower_col, utils.TimeMillisecond(ts)-1, 0, size, reflect.TypeOf(int64(0)))
	if len(uids) == 0 {
		return 0, []int64{}, utils.TimeMillisecond(ts)
	}
	ids := []int64{}
	for _, _uid := range uids {
		_id := *(_uid.(*int64))
		ids = append(ids, _id)
	}
	last_t := utils.TimeMillisecond(ts)
	if len(ids) > 0 {
		_last_t, _ := ssdb.New(use_ssdb_passport_db).Zscore(follower_col, ids[len(ids)-1])
		last_t = _last_t
	}
	return total, ids, last_t
}

func (m *FriendShips) FollowerCounts(uid int64) int {
	follower_col := fmt.Sprintf(friendship_follower_sets, uid)
	total, _ := ssdb.New(use_ssdb_passport_db).Zcard(follower_col)
	return total
}

//互粉的用户
func (m *FriendShips) BothFriendUids(uid int64) []int64 {
	key := fmt.Sprintf(friendship_friend_both_sets, uid)
	uids, _ := ssdb.New(use_ssdb_passport_db).Zrevrange(key, 0, -1, reflect.TypeOf(int64(0)))
	ids := []int64{}
	for _, _uid := range uids {
		_id := *(_uid.(*int64))
		ids = append(ids, _id)
	}
	return ids
}

func (m *FriendShips) BothFriendCounts(uid int64) int {
	key := fmt.Sprintf(friendship_friend_both_sets, uid)
	total, _ := ssdb.New(use_ssdb_passport_db).Zcard(key)
	return total
}

//首字母拼音分组列表
func (m *FriendShips) BothFriendUidsPy(uid int64) map[string][]int64 {
	cache := utils.GetCache()
	ckey := fmt.Sprintf(friendship_bothfriend_pylist, uid)
	var obj map[string][]int64
	err := cache.Get(ckey, &obj)
	if err == nil {
		return obj
	}

	key := fmt.Sprintf(friendship_friend_both_sets, uid)
	uids, _ := ssdb.New(use_ssdb_passport_db).Zrevrange(key, 0, -1, reflect.TypeOf(int64(0)))
	ids := []int64{}
	for _, _uid := range uids {
		_id := *(_uid.(*int64))
		ids = append(ids, _id)
	}
	pymap := m.py(ids)
	cache.Set(ckey, pymap, 48*time.Hour)
	return pymap
}

//关系
func (m *FriendShips) Relation(source_uid int64, target_uid int64) STMemberRelation {
	//o := dbs.NewDefaultOrm()
	//ns, _ := o.QueryTable(&MemberFriend{}).Filter("source_uid", source_uid).Filter("target_uid", target_uid).Count()
	//nt, _ := o.QueryTable(&MemberFriend{}).Filter("source_uid", target_uid).Filter("target_uid", source_uid).Count()

	has_ns, err := ssdb.New(use_ssdb_passport_db).Zexists(fmt.Sprintf(friendship_friend_sets, source_uid), target_uid) //源的好友Set
	has_nt, err := ssdb.New(use_ssdb_passport_db).Zexists(fmt.Sprintf(friendship_friend_sets, target_uid), source_uid) //目标好友Set
	if err != nil {
		fmt.Println("-------------------------friendship_fail:", err)
	}

	mrs := MemberRelation{}
	mrt := MemberRelation{}
	mrs.Uid = source_uid
	mrt.Uid = target_uid
	mrs.NotificationsEnabled = false
	mrt.NotificationsEnabled = false
	if has_ns {
		mrs.Following = true
		mrt.FollowedBy = true
	}
	if has_nt {
		mrt.Following = true
		mrs.FollowedBy = true
	}
	st := STMemberRelation{
		Source: mrs,
		Target: mrt,
	}
	return st
}
