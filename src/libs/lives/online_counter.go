package lives

import (
	"fmt"
	"time"
	"utils"
	"utils/redis"
)

type OnlineCounterMUD struct {
	Uid      int64
	UUID     string
	JoinTime time.Time
}

type OnlineCounter struct{}

func (c *OnlineCounter) init() {

}

func (c *OnlineCounter) channelKey(lt LIVE_TYPE, cid int) string {
	return fmt.Sprintf("live_online_%d_%d", lt, cid)
}

func (c *OnlineCounter) NewUUID() string {
	return utils.RandomStrings(32)
}

func (c *OnlineCounter) AddRandomViewers(liveType LIVE_TYPE, channelId int, min int, max int) {
	if min <= 0 && max <= 0 {
		return
	}
	rands, _ := utils.IntRange(min, max)
	args := []interface{}{}
	ckey := c.channelKey(liveType, channelId)
	for i := 1; i <= rands; i++ {
		args = append(args, time.Now().Unix()+int64(i))
		args = append(args, c.NewUUID())
		if i%500 == 0 || rands == i {
			redis.ZMultiAdd(nil, ckey, args...)
			args = []interface{}{}
		}
	}
}

func (c *OnlineCounter) Reset(liveType LIVE_TYPE, channelId int) error {
	ckey := c.channelKey(liveType, channelId)
	_, err := redis.Del(nil, ckey)
	return err
}

func (c *OnlineCounter) GetChannelCounts(liveType LIVE_TYPE, channelId int) int {
	ckey := c.channelKey(liveType, channelId)
	nums, err := redis.ZCard(nil, ckey)
	if err != nil {
		return 0
	}
	return int(nums)
}

func (c *OnlineCounter) JoinChannel(liveType LIVE_TYPE, channelId int, uid int64, uuid string) error {
	ckey := c.channelKey(liveType, channelId)
	_t, _ := redis.ZScore(nil, ckey, uuid)
	if _t > 0 {
		_now_t := time.Now().Unix()
		_incr := _now_t - _t
		redis.ZIncrBy(nil, ckey, _incr, uuid)
	} else {
		redis.ZAdd(nil, ckey, time.Now().Unix(), uuid)
	}
	md := OnlineCounterMUD{
		Uid:      uid,
		UUID:     uuid,
		JoinTime: time.Now(),
	}
	cache := utils.GetCache()
	cache.Set(uuid, md, 24*time.Hour)
	return nil
}

func (c *OnlineCounter) LeaveChannel(liveType LIVE_TYPE, channelId int, uuid string) error {
	ckey := c.channelKey(liveType, channelId)
	redis.ZRem(nil, ckey, uuid)
	cache := utils.GetCache()
	cache.Delete(uuid)
	return nil
}
