package cache

import (
	//"../utils"
	//"bytes"
	//"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"strconv"
	"time"
)

var (
	// the collection name of redis for cache adapter.
	DefaultKey string = "beecacheRedis"
)

// Redis cache adapter.
type RedisCache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	key      string
}

// create new redis cache with default collection name.
func NewRedisCache() *RedisCache {
	return &RedisCache{key: DefaultKey}
}

// actually do the redis cmds
func (rc *RedisCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Get cache from redis.
// slice input,clazz use &[]object{}
func (rc *RedisCache) Get(key string, clazz interface{}) interface{} {
	v, err := rc.do("HGET", rc.key, key)
	if err != nil {
		return nil
	}
	if clazz != nil {
		kind := reflect.TypeOf(clazz).Kind()
		switch kind {
		case reflect.Ptr:
			elemKind := reflect.TypeOf(clazz).Elem().Kind()
			if elemKind == reflect.Struct || elemKind == reflect.Slice || elemKind == reflect.Array || elemKind == reflect.Map {
				b, err := redis.Bytes(v, err)
				if err != nil {
					return nil
				}
				err = json.Unmarshal(b, clazz)
				if err != nil {
					fmt.Println("ptr get unmarshal:" + err.Error())
					return nil
				}
				return clazz
			}
			return nil
		case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
			b, err := redis.Bytes(v, err)
			if err != nil {
				return nil
			}
			err = json.Unmarshal(b, clazz)
			if err != nil {
				fmt.Println("get unmarshal:" + err.Error())
				return nil
			}
			fmt.Println(b)
			return clazz
		case reflect.Bool:
			tf, _ := redis.Bool(v, err)
			return tf
		case reflect.String:
			s, err := redis.String(v, err)
			if err != nil {
				return nil
			}
			return s
		case reflect.Int:
			n, err := redis.Int(v, err)
			if err != nil {
				return nil
			}
			return n
		case reflect.Int64:
			n64, err := redis.Int64(v, err)
			if err != nil {
				return nil
			}
			return n64
		case reflect.Float64:
			f64, err := redis.Float64(v, err)
			if err != nil {
				return nil
			}
			return f64
		case reflect.TypeOf([]byte{}).Kind():
			bs, err := redis.Bytes(v, err)
			if err != nil {
				return nil
			}
			return bs
		default:
			return v
		}
	}
	return v
}

// put cache to redis.
// timeout is ignored.
func (rc *RedisCache) Put(key string, val interface{}, timeout int64) error {
	kind := reflect.TypeOf(val).Kind()
	if kind == reflect.Struct || (kind == reflect.Ptr && reflect.TypeOf(val).Elem().Kind() == reflect.Struct) || kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map {
		b, err := json.Marshal(val)
		if err != nil {
			fmt.Println("put marshal:" + err.Error())
			return err
		}
		_, err = rc.do("HSET", rc.key, key, b)
		_, err = rc.do("EXPIRE", key, strconv.FormatInt(timeout, 10))
		//fmt.Println(res, err)
		return err
	} else {
		_, err := rc.do("HSET", rc.key, key, val)
		_, err = rc.do("EXPIRE", key, strconv.FormatInt(timeout, 10))
		return err
	}
}

// delete cache in redis.
func (rc *RedisCache) Delete(key string) error {
	_, err := rc.do("HDEL", rc.key, key)
	return err
}

// check cache exist in redis.
func (rc *RedisCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("HEXISTS", rc.key, key))
	if err != nil {
		return false
	}

	return v
}

// increase counter in redis.
func (rc *RedisCache) Incr(key string) error {
	_, err := redis.Bool(rc.do("HINCRBY", rc.key, key, 1))
	return err
}

// decrease counter in redis.
func (rc *RedisCache) Decr(key string) error {
	_, err := redis.Bool(rc.do("HINCRBY", rc.key, key, -1))
	return err
}

// clean all cache in redis. delete this redis collection.
func (rc *RedisCache) ClearAll() error {
	_, err := rc.do("DEL", rc.key)
	return err
}

// start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *RedisCache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}

	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()
	if err := c.Err(); err != nil {
		return err
	}

	return nil
}

// connect to redis.
func (rc *RedisCache) connectInit() {
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rc.conninfo)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func init() {
	Register("redis", NewRedisCache())
}
