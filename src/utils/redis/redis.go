package redis

import (
	"errors"
	"fmt"
	"logs"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"utils/hash"

	"github.com/astaxie/beego"
	"github.com/gosexy/redis"
)

const (
	clientDb = 0
)

var consistent *hash.Consistent = hash.New()
var endPoints map[string]*redisEndPoint = make(map[string]*redisEndPoint)
var locker *sync.RWMutex = new(sync.RWMutex)

type redisEndPoint struct {
	Addr string
	Port uint
}

func AddSrv(name string, endPoint *redisEndPoint) {
	locker.Lock()
	defer locker.Unlock()
	consistent.Add(name)
	endPoints[name] = endPoint
}

func RemoveSrv(name string) {
	locker.Lock()
	defer locker.Unlock()
	consistent.Remove(name)
	delete(endPoints, name)
}

func GetSrv(key string) *redisEndPoint {
	locker.RLock()
	defer locker.RUnlock()
	server, err := consistent.Get(key)
	if err != nil {
		errstr := fmt.Sprintf("redis server not exist or fail:%v", err)
		logs.Errorf(errstr)
		panic(errstr)
	}
	return endPoints[server]
}

func init() {
	redishosts := beego.AppConfig.String("redis.hosts")
	redis_hosts := strings.Split(redishosts, ",")
	for _, host := range redis_hosts {
		uri := strings.Split(host, ":")
		port := 6379
		if len(uri) == 2 {
			port, _ = strconv.Atoi(uri[1])
		}
		ep := &redisEndPoint{uri[0], uint(port)}
		endPoints[host] = ep
		consistent.Add(host)
	}
	//redis_password := beego.AppConfig.String("redis.password")
	redis_client_maxpool, _ := beego.AppConfig.Int("redis.client.maxpool")
	if redis_client_maxpool == 0 {
		redis_client_maxpool = 50
	}
}

//普通
func Del(client *redis.Client, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, fmt.Errorf("keys not empty")
	}
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(keys[0])
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.Del(keys...)
}

func Set(client *redis.Client, key string, value interface{}) (ret string, err error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return "", err
		}
	}
	return client.Set(key, value)
}

func Get(client *redis.Client, key string) (string, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return "", err
		}
	}
	return client.Get(key)
}

func Exists(client *redis.Client, key string) (bool, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return false, err
		}
	}
	return client.Exists(key)
}

func Incr(client *redis.Client, key string) (int64, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.Incr(key)
}

//Hset
func HSet(client *redis.Client, key string, field string, value interface{}) (bool, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return false, err
		}
	}
	return client.HSet(key, field, value)
}

func HGet(client *redis.Client, key string, field string) (string, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return "", err
		}
	}
	return client.HGet(key, field)
}

func HKeys(client *redis.Client, key string) ([]string, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return nil, err
		}
	}
	return client.HKeys(key)
}

func HDel(client *redis.Client, key string, fields ...string) (int64, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.HDel(key, fields...)
}

func HGetAll(client *redis.Client, key string) ([]string, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return nil, err
		}
	}
	return client.HGetAll(key)
}

//Sort Set
func ZCard(client *redis.Client, key string) (int64, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZCard(key)
}

func ZAdd(client *redis.Client, key string, score interface{}, obj interface{}) (int64, error) {
	b, err := Serialize(obj)
	if err != nil {
		return 0, err
	}
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZAdd(key, score, b)
}

func ZMultiAdd(client *redis.Client, key string, arguments ...interface{}) (int64, error) {
	if len(arguments)%2 != 0 {
		return 0, errors.New("Failed to relate SCORE -> MEMBER using the given arguments.")
	}
	args := []interface{}{}
	for i := 0; i < len(arguments); i += 2 {
		b, err := Serialize(arguments[i+1])
		if err != nil {
			return 0, err
		}
		args = append(args, arguments[i])
		args = append(args, b)
	}
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZAdd(key, args...)
}

func ZRange(client *redis.Client, key string, valType reflect.Type, arguments ...interface{}) ([]interface{}, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return nil, err
		}
	}
	result, err := client.ZRange(key, arguments...)
	lst := []interface{}{}
	for _, byt := range result {
		obj := reflect.New(valType).Interface()
		err = Deserialize([]byte(byt), obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func ZRangeByScore(client *redis.Client, key string, valType reflect.Type, arguments ...interface{}) ([]interface{}, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return nil, err
		}
	}
	result, err := client.ZRangeByScore(key, arguments...)
	lst := []interface{}{}
	for _, byt := range result {
		obj := reflect.New(valType).Interface()
		err = Deserialize([]byte(byt), obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func ZRevRangeByScore(client *redis.Client, key string, valType reflect.Type, max interface{}, min interface{}, arguments ...interface{}) ([]interface{}, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return nil, err
		}
	}
	result, err := client.ZRevRangeByScore(key, max, min, arguments...)
	lst := []interface{}{}
	for _, byt := range result {
		obj := reflect.New(valType).Interface()
		err = Deserialize([]byte(byt), obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func ZRevRangeByScoreWithScores(client *redis.Client, key string, valType reflect.Type, max interface{}, min interface{}, arguments ...interface{}) ([]interface{}, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return nil, err
		}
	}
	args := []interface{}{}
	args = append(args, "WITHSCORES")
	args = append(args, arguments...)
	result, err := client.ZRevRangeByScore(key, max, min, args...)
	lst := []interface{}{}
	for i := 0; i < len(result); i += 2 {
		obj := reflect.New(valType).Interface()
		err = Deserialize([]byte(result[i]), obj)
		if err == nil {
			lst = append(lst, obj)
			score, _ := strconv.ParseFloat(result[i+1], 64)
			lst = append(lst, score)
		}
	}
	return lst, nil
}

func ZRemRangeByRank(client *redis.Client, key string, start int, stop int) (int64, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZRemRangeByRank(key, start, stop)
}

func ZRemRangeByScore(client *redis.Client, key string, min interface{}, max interface{}) (int64, error) {
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZRemRangeByScore(key, min, max)
}

func ZRem(client *redis.Client, key string, arguments ...interface{}) (int64, error) {
	bs := []interface{}{}
	for _, obj := range arguments {
		b, err := Serialize(obj)
		if err != nil {
			return 0, err
		}
		bs = append(bs, b)
	}
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZRem(key, bs...)
}

//bug score 不能很大
func ZScore(client *redis.Client, key string, member interface{}) (int64, error) {
	b, err := Serialize(member)
	if err != nil {
		return 0, err
	}
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return 0, err
		}
	}
	return client.ZScore(key, b)
}

func ZIncrBy(client *redis.Client, key string, increment int64, obj interface{}) (string, error) {
	b, err := Serialize(obj)
	if err != nil {
		return "", err
	}
	if client == nil {
		client = redis.New()
		endPoint := GetSrv(key)
		err := client.Connect(endPoint.Addr, endPoint.Port)
		defer client.Quit()
		if err != nil {
			return "", err
		}
	}
	return client.ZIncrBy(key, increment, b)
}
