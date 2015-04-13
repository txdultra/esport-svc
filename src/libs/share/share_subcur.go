package share

import (
	"dbs"
	"fmt"
	"libs/passport"
	"logs"
	"reflect"
	"time"
	"utils"
	"utils/redis"
	"utils/ssdb"
)

const (
	SHARE_MY_VOD_SUBSCRIPTIONS       = "share.my.subscriptions:%d"
	SHARE_MY_VOD_SUBSCRIPTIONS_COUNT = "share.my.subscriptions.counts:%d"
)

type ShareVodSubcurs struct{}

//只支持视频
func (n *ShareVodSubcurs) Gets(uid int64, page int, size int, ts time.Time) (int, []*Share) {
	ss := []*Share{}
	if uid <= 0 {
		return 0, ss
	}
	vod_col := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS, uid)
	col_size, err := redis.ZCard(nil, vod_col)
	if err != nil {
		//log
	}
	if col_size == 0 && err == nil {
		return 0, ss
	}
	list, err := redis.ZRevRangeByScore(nil, vod_col, reflect.TypeOf(int64(0)), utils.TimeMillisecond(ts)-1, 0, "LIMIT", 0, size)
	if err != nil {
		logs.Errorf("my subscription fail:%v", err)
	}
	s := &Shares{}
	for _, obj := range list {
		if id, ok := obj.(*int64); ok {
			share := s.Get(*id)
			if share != nil {
				ss = append(ss, share)
			}
		}
	}
	return int(col_size), ss
}

func (n *ShareVodSubcurs) distributeToUserBox(s *Share, uids []int64) {
	for _, uid := range uids {
		vod_follower_colname := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS, uid)
		//加入订阅列表
		redis.ZAdd(nil, vod_follower_colname, utils.TimeMillisecond(s.CreateTime), s.Id)
		redis.ZRemRangeByRank(nil, vod_follower_colname, 0, -mbox_subscription_length) //限制box大小
		n.IncrEventCount(uid)
	}
}

func (n *ShareVodSubcurs) deleteOnUserBox(s *Share, uids []int64) {
	for _, uid := range uids {
		vod_follower_colname := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS, uid)
		//删除订阅列表中的share
		redis.ZRem(nil, vod_follower_colname, s.Id)
		n.DecrEventCount(uid)
	}
}

func (n ShareVodSubcurs) NotifyFans(s *Share) {
	if s == nil {
		return
	}
	//现在只需要支持视频
	if s.ShareType&int(SHARE_KIND_VOD) != int(SHARE_KIND_VOD) {
		return
	}
	pst := &passport.FriendShips{}
	page := 1
	t := time.Now()
	for {
		_, follower_uids, last_t := pst.FollowerIdsP(s.Uid, 500, page, t)
		if len(follower_uids) == 0 {
			break
		}
		go n.distributeToUserBox(s, follower_uids)
		page++
		t = utils.MsToTime(int64(last_t))
	}
}

func (n ShareVodSubcurs) RecoverFans(s *Share) {
	if s == nil {
		return
	}
	//现在只需要支持视频
	if s.ShareType&int(SHARE_KIND_VOD) != int(SHARE_KIND_VOD) {
		return
	}
	pst := &passport.FriendShips{}
	page := 1
	t := time.Now()
	for {
		_, follower_uids, last_t := pst.FollowerIdsP(s.Uid, 500, page, t)
		if len(follower_uids) == 0 {
			break
		}
		go n.deleteOnUserBox(s, follower_uids)
		page++
		t = utils.MsToTime(int64(last_t))
	}
}

//IEventInDecrCounter interface
func (n ShareVodSubcurs) IncrEventCount(uid int64) int {
	vod_follower_count_col := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS_COUNT, uid)
	c, _ := ssdb.New(use_ssdb_share_db).Incr(vod_follower_count_col)
	return int(c)
}

//IEventInDecrCounter interface
func (n ShareVodSubcurs) DecrEventCount(uid int64) int {
	vod_follower_count_col := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS_COUNT, uid)
	c, _ := ssdb.New(use_ssdb_share_db).Decr(vod_follower_count_col)
	if c < 0 {
		ssdb.New(use_ssdb_share_db).Incrby(vod_follower_count_col, -c)
		return 0
	}
	return int(c)
}

//IEventCounter interface
func (n ShareVodSubcurs) ResetEventCount(uid int64) bool {
	follower_count_col := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS_COUNT, uid)
	ok, _ := ssdb.New(use_ssdb_share_db).Del(follower_count_col)
	return ok
}

//IEventCounter interface
func (n ShareVodSubcurs) NewEventCount(uid int64) int {
	c := 0
	follower_count_col := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS_COUNT, uid)
	err := ssdb.New(use_ssdb_share_db).Get(follower_count_col, &c)
	if err != nil {
		return 0
	}
	return c
}

//friendship interface
func (n ShareVodSubcurs) FSFriendDo(suid int64, tuid int64) {
	//订阅,把被关注对象的分享视频装入源对象的box中
	type Result struct {
		Id int64
		Ts int64
	}
	my_colname := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS, suid)
	o := dbs.NewOrm(share_db)
	s := &Shares{}
	tbl := s.Hash_tbl(tuid, share_tbl)
	sql := fmt.Sprintf("select id,ts from %s where uid=%d and st=%d order by ts desc limit %d", tbl, tuid, int(SHARE_KIND_VOD), mbox_subscription_length)
	var res []Result
	o.Raw(sql).QueryRows(&res)

	if len(res) == 0 {
		return
	}

	adds := []interface{}{}
	for _, r := range res {
		adds = append(adds, r.Ts)
		adds = append(adds, r.Id)
	}
	redis.ZMultiAdd(nil, my_colname, adds...)
	redis.ZRemRangeByRank(nil, my_colname, 0, -mbox_subscription_length) //限制box大小
}

func (n ShareVodSubcurs) FSUnFriendDo(suid int64, tuid int64) {
	//订阅,重新筛选内容装入源对象的box中
	type Result struct {
		Id int64
	}
	o := dbs.NewOrm(share_db)
	s := &Shares{}
	tbl_b := s.Hash_tbl(tuid, share_tbl)
	sql := fmt.Sprintf("select id from %s where uid=%d and st=%d order by ts desc limit %d", tbl_b, tuid, int(SHARE_KIND_VOD), mbox_subscription_length)
	var res []Result
	o.Raw(sql).QueryRows(&res)

	if len(res) == 0 {
		return
	}

	my_colname := fmt.Sprintf(SHARE_MY_VOD_SUBSCRIPTIONS, suid)
	adds := []interface{}{}
	for _, r := range res {
		adds = append(adds, r.Id)
	}
	redis.ZRem(nil, my_colname, adds...)
}
