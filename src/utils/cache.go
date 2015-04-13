package utils

import (
	"utils/cache"
	//"github.com/astaxie/beego"
	//"sync"
	"fmt"
	"time"
)

func GetCache() cache.Cache {
	return cache.Instance
}

func GetLocalCache() cache.Cache {
	return cache.Local
}

////////////////////////////////////////////////////////////////////////////////
//时间间隔缓存
////////////////////////////////////////////////////////////////////////////////
type CACHE_INTERVAL_TIME_TYPE int

const (
	CACHE_INTERVAL_TIME_TYPE_HOUR   CACHE_INTERVAL_TIME_TYPE = 1
	CACHE_INTERVAL_TIME_TYPE_MINUTE CACHE_INTERVAL_TIME_TYPE = 2
	CACHE_INTERVAL_TIME_TYPE_SECOND CACHE_INTERVAL_TIME_TYPE = 3
)

func createTD(t time.Time, interval CACHE_INTERVAL_TIME_TYPE) (string, time.Duration) {
	y := t.Year()
	m := t.Month()
	d := t.Day()
	h := t.Hour()
	mm := t.Minute()
	ss := t.Second()
	var query_cache_time_key string
	var cache_duration time.Duration

	switch interval {
	case CACHE_INTERVAL_TIME_TYPE_SECOND:
		query_cache_time_key = fmt.Sprintf("y:%dm:%d:d:%d:h:%d:mm:%d:ss:%d", y, m, d, h, mm, ss)
		cache_duration = 1 * time.Second
	case CACHE_INTERVAL_TIME_TYPE_HOUR:
		query_cache_time_key = fmt.Sprintf("y:%dm:%d:d:%d:h:%d", y, m, d, h)
		cache_duration = 1 * time.Hour
	case CACHE_INTERVAL_TIME_TYPE_MINUTE:
		query_cache_time_key = fmt.Sprintf("y:%dm:%d:d:%d:h:%d:mm:%d", y, m, d, h, mm)
		cache_duration = 1 * time.Minute
	}
	return query_cache_time_key, cache_duration
}

func SetLocalFastTimePartCache(t time.Time, key string, interval CACHE_INTERVAL_TIME_TYPE, data interface{}) {
	qkey, dur := createTD(t, interval)
	cache_key := fmt.Sprintf("%s_%s", key, qkey)
	GetLocalCache().Set(cache_key, data, dur)
}

func SetLocalFastExpriesTimePartCache(duration time.Duration, key string, data interface{}) {
	GetLocalCache().Set(key, data, duration)
}

func GetLocalFastTimePartCache(t time.Time, key string, interval CACHE_INTERVAL_TIME_TYPE) interface{} {
	qkey, _ := createTD(t, interval)
	cache_key := fmt.Sprintf("%s_%s", key, qkey)
	var obj interface{}
	err := GetLocalCache().Get(cache_key, &obj)
	if err == nil {
		return obj
	}
	return nil
}

func GetLocalFastExpriesTimePartCache(key string) interface{} {
	var obj interface{}
	err := GetLocalCache().Get(key, &obj)
	if err == nil {
		return obj
	}
	return nil
}

//cache v1
//var wcache cache.Cache
//var wmutex *sync.Mutex = new(sync.Mutex)

//var lcache cache.Cache
//var lmutex *sync.Mutex = new(sync.Mutex)

//func GetCache() cache.Cache {
//	if wcache == nil {
//		wmutex.Lock()
//		defer func() {
//			wmutex.Unlock()
//		}()
//		provider := beego.AppConfig.String("cache_provider")
//		var err error
//		if wcache == nil {
//			switch provider {
//			case "redis":
//				wcache, err = cache.NewCache("redis", `{"conn":"192.168.233.128:6379"}`)
//			case "memcache":
//				wcache, err = cache.NewCache("memcache", `127.0.0.1:11211`)
//			default:
//				wcache, err = cache.NewCache("memory", `{"interval":60}`)
//			}
//			if err != nil {
//				panic(err)
//			}
//		}
//	}
//	return wcache
//}

//func LocalCache() cache.Cache {
//	if lcache == nil {
//		lmutex.Lock()
//		defer func() {
//			lmutex.Unlock()
//		}()
//		if lcache == nil {
//			provider := beego.AppConfig.String("cache_provider")
//			if provider != "memory" {
//				lcache, _ = cache.NewCache("memory", `{"interval":60}`)
//			} else {
//				_wc := GetCache()
//				lcache = _wc
//			}
//		}
//	}
//	return lcache
//}
