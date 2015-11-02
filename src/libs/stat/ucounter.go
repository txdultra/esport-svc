package stat

import (
	"fmt"
	"reflect"
	"utils/ssdb"

	"github.com/astaxie/beego"
)

const (
	USER_COUNT_LIST_KEY = "user_count_list_%d"
)

var ucdb = beego.AppConfig.String("ssdb.counter.db")
var ucounters map[string]func(uid int64) string = make(map[string]func(uid int64) string)

func RegisterUCountKey(module string, uFunc func(uid int64) string) {
	if _, ok := ucounters[module]; !ok {
		ucounters[module] = uFunc
	}
}

func getUCountKey(module string) func(uid int64) string {
	if _func, ok := ucounters[module]; ok {
		return _func
	}
	panic("未注册" + module + "用户计数器")
}

func UCIncrCount(uid int64, module string) int64 {
	key := getUCountKey(module)(uid)
	ukey := fmt.Sprintf(USER_COUNT_LIST_KEY, uid)
	c, _ := ssdb.New(ucdb).Zincrby(ukey, key, 1)
	return c
}

func UCDecrCount(uid int64, module string) int64 {
	key := getUCountKey(module)(uid)
	ukey := fmt.Sprintf(USER_COUNT_LIST_KEY, uid)
	c, _ := ssdb.New(ucdb).Zincrby(ukey, key, -1)
	if c < 0 {
		ssdb.New(ucdb).Zincrby(ukey, key, -c)
		return 0
	}
	return c
}

func UCResetCount(uid int64, module string) bool {
	key := getUCountKey(module)(uid)
	ukey := fmt.Sprintf(USER_COUNT_LIST_KEY, uid)
	ok, _ := ssdb.New(ucdb).Zrem(ukey, key)
	return ok
}

func UCGetCount(uid int64, module string) int64 {
	key := getUCountKey(module)(uid)
	ukey := fmt.Sprintf(USER_COUNT_LIST_KEY, uid)
	i, _ := ssdb.New(ucdb).Zscore(ukey, key)
	return i
}

func UCSetCount(uid int64, module string, score int64) (bool, error) {
	key := getUCountKey(module)(uid)
	ukey := fmt.Sprintf(USER_COUNT_LIST_KEY, uid)
	return ssdb.New(ucdb).Zadd(ukey, key, score)
}

func UCGetCounts(uid int64, modules []string) map[string]int64 {
	ukey := fmt.Sprintf(USER_COUNT_LIST_KEY, uid)
	kvs, err := ssdb.New(ucdb).ZscanKS(ukey, 0, 1<<10, 1000, reflect.TypeOf(""), reflect.TypeOf(int64(0)))
	kcs := make(map[string]int64)
	for _, module := range modules {
		key := getUCountKey(module)(uid)
		if err != nil {
			kcs[module] = 0
		} else {
			has := false
			for _, ks := range kvs {
				if key == *(ks.Key.(*string)) {
					kcs[module] = *(ks.Score.(*int64))
					has = true
					break
				}
			}
			if !has {
				kcs[module] = 0
			}
		}
	}
	return kcs
}
