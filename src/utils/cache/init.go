package cache

import (
	//"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

func init() {
	clientTimeout := 500 * time.Millisecond
	// Set the default expiration time.
	defaultExpiration := time.Hour // The default for the default is one hour.
	if expireStr := beego.AppConfig.String("cache.expires"); expireStr != "" {
		if _defaultExpiration, err := time.ParseDuration(expireStr); err != nil {
			panic("Could not parse default cache expiration duration " + expireStr + ": " + err.Error())
		} else {
			defaultExpiration = _defaultExpiration
		}
	}
	if timeout := beego.AppConfig.String("cache.timeout"); timeout != "" {
		if _clientTimeout, err := time.ParseDuration(timeout); err != nil {
			panic("Could not parse default cache timeout duration " + timeout + ": " + err.Error())
		} else {
			clientTimeout = _clientTimeout
		}
	}

	// make sure you aren't trying to use both memcached and redis
	ismc, _ := beego.AppConfig.Bool("cache.memcached")
	isre, _ := beego.AppConfig.Bool("cache.redis")
	if ismc && isre {
		panic("You've configured both memcached and redis, please only include configuration for one cache!")
	}
	//default instance local memory
	Local = NewInMemoryCache(defaultExpiration)

	// Use memcached?
	if ismc {
		hosts := strings.Split(beego.AppConfig.String("cache.hosts"), ",")
		if len(hosts) == 0 {
			panic("Memcache enabled but no memcached hosts specified!")
		}

		Instance = NewMemcachedCache(hosts, defaultExpiration, clientTimeout)
		return
	}

	// Use Redis (share same config as memcached)?
	if isre {
		hosts := strings.Split(beego.AppConfig.String("cache.hosts"), ",")
		if len(hosts) == 0 {
			panic("Redis enabled but no Redis hosts specified!")
		}
		if len(hosts) > 1 {
			panic("Redis currently only supports one host!")
		}
		password := beego.AppConfig.String("cache.redis.password")
		Instance = NewRedisCache(hosts[0], password, defaultExpiration, clientTimeout)
		return
	}

	// By default, use the in-memory cache.
	Instance = NewInMemoryCache(defaultExpiration)
}
