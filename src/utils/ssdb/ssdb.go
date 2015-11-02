package ssdb

import (
	"fmt"
	"logs"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/hoisie/redis"
	//"github.com/gosexy/redis"
	"utils/hash"

	"github.com/ssdb"
)

const (
	consistentHashReplicasCount = 100
	consistentHashBucketsCount  = 1024
	clientDb                    = 0
	ssdb_chs                    = "abcdefghijklmnopqrstuvwxyz"
)

var clusters map[string]*Cluster = make(map[string]*Cluster)
var locker *sync.RWMutex = new(sync.RWMutex)

type endPoint struct {
	Addr string
	Port int
}

type KeyScore struct {
	Key   interface{}
	Score interface{}
}

type Cluster struct {
	consistent *hash.Consistent
	clients    map[string]*ssdbClient
}

func (c *Cluster) SetConsistentClients(consistent *hash.Consistent, clients map[string]*ssdbClient) {
	c.consistent = consistent
	c.clients = clients
}

func (c *Cluster) AddSrv(name string, client *ssdbClient) {
	locker.Lock()
	defer locker.Unlock()
	c.consistent.Add(name)
	c.clients[name] = client
}

func (c *Cluster) RemoveSrv(name string) {
	locker.Lock()
	defer locker.Unlock()
	c.consistent.Remove(name)
	delete(c.clients, name)
}

func (c *Cluster) GetSrv(key string) *ssdbClient {
	locker.RLock()
	defer locker.RUnlock()
	server, err := c.consistent.Get(key)
	if err != nil {
		errstr := fmt.Sprintf("ssdb server not exist or fail:%v", err)
		logs.Errorf(errstr)
		panic(errstr)
	}
	return c.clients[server]
}

func init() {
	clts := make(map[string][]*endPoint)
	for _, c := range ssdb_chs {
		hosts := beego.AppConfig.String("ssdb.addr." + string(c))
		if len(hosts) == 0 {
			continue
		}
		addrs := strings.Split(hosts, ",")
		eps := []*endPoint{}
		for _, addr := range addrs {
			ap := strings.Split(addr, ":")
			port := 6379
			if len(ap) == 2 {
				port, _ = strconv.Atoi(ap[1])
			}
			eps = append(eps, &endPoint{ap[0], port})
		}
		clts[string(c)] = eps
	}
	ssdb_password := beego.AppConfig.String("ssdb.password")
	ssdb_client_maxpool, _ := beego.AppConfig.Int("ssdb.client.maxpool")
	if ssdb_client_maxpool == 0 {
		ssdb_client_maxpool = 50
	}

	for k, v := range clts {
		consistent := hash.New()
		clients := make(map[string]*ssdbClient)
		for _, ep := range v {
			addr := fmt.Sprintf("%s:%d", ep.Addr, ep.Port)
			rc := &ssdbClient{
				&redis.Client{
					Addr:        addr,
					Db:          clientDb,
					Password:    ssdb_password,
					MaxPoolSize: ssdb_client_maxpool,
				},
				false,
				ep.Addr,
				ep.Port,
			}
			consistent.Add(addr)
			clients[addr] = rc
		}
		c := &Cluster{consistent, clients}
		clusters[k] = c
	}
}

func New(c string) *Cluster {
	if c, ok := clusters[c]; ok {
		return c
	}
	panic(fmt.Sprintf("ssdb group %s not setting", c))
}

type ssdbClient struct {
	*redis.Client
	readonly  bool
	conn_addr string
	conn_port int
}

//
//普通操作
//
func (c *Cluster) Exists(key string) (bool, error) {
	return c.GetSrv(key).Exists(key)
}

func (c *Cluster) Del(key string) (bool, error) {
	return c.GetSrv(key).Del(key)
}

func (c *Cluster) Set(key string, val interface{}) error {
	b, err := Serialize(val)
	if err != nil {
		return err
	}
	return c.GetSrv(key).Set(key, b)
}

func (c *Cluster) MultiGet(keys []string, valType reflect.Type) []interface{} {
	if len(keys) == 0 {
		return []interface{}{}
	}
	client := c.GetSrv(keys[0])
	byts, err := client.Mget(keys...)
	if err != nil {
		return []interface{}{}
	}
	lst := []interface{}{}
	for _, byt := range byts {
		obj := reflect.New(valType).Interface()
		err = Deserialize(byt, obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst
}

func (c *Cluster) MultiDel(keys []string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	client := c.GetSrv(keys[0])
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return 0, err
	}
	args := []interface{}{}
	args = append(args, "multi_del")
	for _, key := range keys {
		args = append(args, key)
	}
	resp, err := db.Do(args...)
	if err == nil {
		if len(resp) == 2 && resp[0] == "ok" {
			count, _ := strconv.Atoi(resp[1])
			return count, nil
		}
	}
	return 0, err
}

//
//strings
//
func (c *Cluster) Incr(key string) (int64, error) {
	return c.GetSrv(key).Incr(key)
}

func (c *Cluster) Incrby(key string, val int64) (int64, error) {
	return c.GetSrv(key).Incrby(key, val)
}

func (c *Cluster) Decr(key string) (int64, error) {
	return c.GetSrv(key).Decr(key)
}

func (c *Cluster) Decrby(key string, val int64) (int64, error) {
	return c.GetSrv(key).Decrby(key, val)
}

func (c *Cluster) Get(key string, ptr interface{}) error {
	byt, err := c.GetSrv(key).Get(key)
	if err != nil {
		return err
	}
	err = Deserialize(byt, ptr)
	if err != nil {
		return err
	}
	return nil
}

//
//hash
//

func (c *Cluster) Hset(key string, field string, val interface{}) (bool, error) {
	b, err := Serialize(val)
	if err != nil {
		return false, err
	}
	return c.GetSrv(key).Hset(key, field, b)
}

func (c *Cluster) Hget(key string, field string, ptr interface{}) error {
	byt, err := c.GetSrv(key).Hget(key, field)
	if err != nil {
		return err
	}
	err = Deserialize(byt, ptr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) Hlen(key string) (int, error) {
	return c.GetSrv(key).Hlen(key)
}

func (c *Cluster) Hexists(key string, field string) (bool, error) {
	return c.GetSrv(key).Hexists(key, field)
}

func (c *Cluster) Hdel(key string, field string) (bool, error) {
	return c.GetSrv(key).Hdel(key, field)
}

func (c *Cluster) Hkeys(key string) ([]string, error) {
	return c.GetSrv(key).Hkeys(key)
}

//ssdb use,删除整个hash
func (c *Cluster) Hclear(key string) error {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil
	}
	_, err = db.Do("hclear", key)
	return err
}

func (c *Cluster) Hmset(key string, mapping interface{}) error {
	return c.GetSrv(key).Hmset(key, mapping)
}

func (c *Cluster) Hmget(key string, valType reflect.Type, fields ...string) ([]interface{}, error) {
	byts, err := c.GetSrv(key).Hmget(key, fields...)
	if err != nil {
		return nil, err
	}
	lst := []interface{}{}
	for _, byt := range byts {
		obj := reflect.New(valType).Interface()
		err = Deserialize(byt, obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func (c *Cluster) Hgetall(key string) ([]string, error) {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return []string{}, nil
	}
	resp, err := db.Do("hgetall", key)
	if err == nil {
		if resp[0] == "ok" {
			return resp[1:], nil
		} else {
			return []string{}, fmt.Errorf(resp[0])
		}
	}
	return []string{}, err
}

//
//列表操作
//
func (c *Cluster) Lpush(key string, val interface{}) error {
	b, err := Serialize(val)
	if err != nil {
		return err
	}
	return c.GetSrv(key).Lpush(key, b)
}

//ssdb use //删除整个list
func (c *Cluster) Lclear(key string) error {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil
	}
	_, err = db.Do("qclear", key)
	return err
}

func (c *Cluster) Ltrim(key string, start int, end int) error {
	return c.GetSrv(key).Ltrim(key, start, end)
}

func (c *Cluster) Lrem(key string, count int, val interface{}) (int, error) {
	b, err := Serialize(val)
	if err != nil {
		return 0, err
	}
	return c.GetSrv(key).Lrem(key, count, b)
}

func (c *Cluster) Llen(key string) (int, error) {
	return c.GetSrv(key).Llen(key)
}

func (c *Cluster) Lrange(key string, start int, end int, valType reflect.Type) ([]interface{}, error) {
	byts, err := c.GetSrv(key).Lrange(key, start, end)
	if err != nil {
		return nil, err
	}
	lst := []interface{}{}
	for _, byt := range byts {
		obj := reflect.New(valType).Interface()
		err = Deserialize(byt, obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func (c *Cluster) Lindex(key string, index int, ptr interface{}) error {
	byt, err := c.GetSrv(key).Lindex(key, index)
	if err != nil {
		return err
	}
	err = Deserialize(byt, ptr)
	if err != nil {
		return err
	}
	return nil
}

// sorted set commands
func (c *Cluster) Zadd(key string, val interface{}, score int64) (bool, error) {
	b, err := Serialize(val)
	if err != nil {
		return false, err
	}
	return c.GetSrv(key).Zadd(key, b, score)
}

func (c *Cluster) MultiZadd(key string, vals []interface{}, scores []int64) (int, error) {
	if len(vals) != len(scores) {
		return 0, fmt.Errorf("vals lens must equal scores lens")
	}
	if len(vals) == 0 {
		return 0, nil
	}
	args := []interface{}{}
	args = append(args, "multi_zset")
	args = append(args, key) //key
	i := 0
	for _, val := range vals {
		b, err := Serialize(val)
		if err != nil {
			return 0, err
		}
		score := scores[i]
		args = append(args, b)
		args = append(args, score)
		i++
	}
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	resp, err := db.Do2(args)
	if len(resp) == 2 && resp[0] == "ok" {
		count, _ := strconv.Atoi(resp[1])
		return count, nil
	}
	return 0, err
}

func (c *Cluster) MultiZget(key string, vals []interface{}, valType reflect.Type) (map[interface{}]int64, error) {
	args := []interface{}{}
	args = append(args, "multi_zget")
	args = append(args, key) //key
	for _, val := range vals {
		b, err := Serialize(val)
		if err != nil {
			return nil, err
		}
		args = append(args, b)
	}
	client := c.GetSrv(key)
	db, _ := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	resp, err := db.Do(args...)
	maps := make(map[interface{}]int64)
	if len(resp) > 0 && resp[0] == "ok" {
		for i := 1; i < len(resp); i += 2 {
			val := resp[i]
			obj := reflect.New(valType).Interface()
			err = Deserialize([]byte(val), obj)
			if err != nil {
				continue
			}
			score := resp[i+1]
			maps[obj], _ = strconv.ParseInt(score, 10, 64)
		}
	}
	return maps, err
}

func (c *Cluster) MultiZdel(key string, vals []interface{}) (int, error) {
	if len(vals) == 0 {
		return 0, nil
	}
	args := []interface{}{}
	args = append(args, "multi_zdel")
	args = append(args, key) //key
	for _, val := range vals {
		b, err := Serialize(val)
		if err != nil {
			return 0, err
		}
		args = append(args, b)
	}
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	resp, err := db.Do2(args)
	if len(resp) == 2 && resp[0] == "ok" {
		count, _ := strconv.Atoi(resp[1])
		return count, nil
	}
	return 0, err
}

//ssdb use 删除整个sort set
func (c *Cluster) Zclear(key string) error {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil
	}
	_, err = db.Do("zclear", key)
	return err
}

func (c *Cluster) Zrem(key string, val interface{}) (bool, error) {
	b, err := Serialize(val)
	if err != nil {
		return false, err
	}
	return c.GetSrv(key).Zrem(key, b)
}

func (c *Cluster) Zcard(key string) (int, error) {
	return c.GetSrv(key).Zcard(key)
}

func (c *Cluster) Zremrangebyrank(key string, start int, end int) (int, error) {
	return c.GetSrv(key).Zremrangebyrank(key, start, end)
}

//[]interface{} =>[]ptr 指针对象
func (c *Cluster) Zrange(key string, start int, end int, valType reflect.Type) ([]interface{}, error) {
	byts, err := c.GetSrv(key).Zrange(key, start, end)
	if err != nil {
		return nil, err
	}
	lst := []interface{}{}
	for _, byt := range byts {
		obj := reflect.New(valType).Interface()
		err = Deserialize(byt, obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

//[]interface{} =>[]ptr 指针对象
func (c *Cluster) Zrevrange(key string, start int, end int, valType reflect.Type) ([]interface{}, error) {
	byts, err := c.GetSrv(key).Zrevrange(key, start, end)
	if err != nil {
		return nil, err
	}
	lst := []interface{}{}
	for _, byt := range byts {
		obj := reflect.New(valType).Interface()
		err = Deserialize(byt, obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func (c *Cluster) Zremrangebyscore(key string, start int64, end int64) (int, error) {
	return c.GetSrv(key).Zremrangebyscore(key, start, end)
}

func (c *Cluster) Zcount(key string, min int64, max int64) (int, error) {
	return c.GetSrv(key).Zcount(key, min, max)
}

func (c *Cluster) Zrangebyscore(key string, start int64, end int64, valType reflect.Type) ([]interface{}, error) {
	byts, err := c.GetSrv(key).Zrangebyscore(key, start, end)
	if err != nil {
		return nil, err
	}
	lst := []interface{}{}
	for _, byt := range byts {
		obj := reflect.New(valType).Interface()
		err = Deserialize(byt, obj)
		if err == nil {
			lst = append(lst, obj)
		}
	}
	return lst, nil
}

func (c *Cluster) Zscan(key string, min int64, max int64, limit int, valType reflect.Type) ([]interface{}, error) {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	resp, err := db.Do("zscan", key, "", min, max, limit)
	if err != nil {
		return nil, err
	}
	if resp[0] == "ok" {
		lst := []interface{}{}
		for i := 1; i < len(resp); i += 2 {
			obj := reflect.New(valType).Interface()
			byts := []byte(resp[i])
			err = Deserialize(byts, obj)
			if err == nil {
				lst = append(lst, obj)
			}
		}
		return lst, nil
	}
	return nil, fmt.Errorf("resp fail:", resp)
}

func (c *Cluster) ZscanKS(key string, min int64, max int64, limit int, keyType reflect.Type, scoreType reflect.Type) ([]*KeyScore, error) {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	resp, err := db.Do("zscan", key, "", min, max, limit)
	if err != nil {
		return nil, err
	}
	if resp[0] == "ok" {
		lst := []*KeyScore{}
		for i := 1; i < len(resp); i += 2 {
			key := reflect.New(keyType).Interface()
			val := reflect.New(scoreType).Interface()
			b1 := []byte(resp[i])
			err1 := Deserialize(b1, key)
			b2 := []byte(resp[i+1])
			err2 := Deserialize(b2, val)
			if err1 == nil && err2 == nil {
				lst = append(lst, &KeyScore{
					key,
					val,
				})
			}
		}
		return lst, nil
	}
	return nil, fmt.Errorf("resp fail:", resp)
}

func (c *Cluster) Zrscan(key string, max int64, min int64, limit int, valType reflect.Type) ([]interface{}, error) {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	resp, err := db.Do("zrscan", key, "", max, min, limit)
	if err != nil {
		return nil, err
	}
	if resp[0] == "ok" {
		lst := []interface{}{}
		for i := 1; i < len(resp); i += 2 {
			obj := reflect.New(valType).Interface()
			byts := []byte(resp[i])
			err = Deserialize(byts, obj)
			if err == nil {
				lst = append(lst, obj)
			}
		}
		return lst, nil
	}
	return nil, fmt.Errorf("resp fail:", resp)
}

func (c *Cluster) ZrscanKS(key string, max int64, min int64, limit int, keyType reflect.Type, scoreType reflect.Type) ([]*KeyScore, error) {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	resp, err := db.Do("zrscan", key, "", max, min, limit)
	if err != nil {
		return nil, err
	}
	if resp[0] == "ok" {
		lst := []*KeyScore{}
		for i := 1; i < len(resp); i += 2 {
			key := reflect.New(keyType).Interface()
			val := reflect.New(scoreType).Interface()
			b1 := []byte(resp[i])
			err1 := Deserialize(b1, key)
			b2 := []byte(resp[i+1])
			err2 := Deserialize(b2, val)
			if err1 == nil && err2 == nil {
				lst = append(lst, &KeyScore{
					key,
					val,
				})
			}
		}
		return lst, nil
	}
	return nil, fmt.Errorf("resp fail:", resp)
}

func (c *Cluster) Zrank(key string, val interface{}) (int, error) {
	b, err := Serialize(val)
	if err != nil {
		return -1, err
	}
	return c.GetSrv(key).Zrank(key, b)
}

func (c *Cluster) Zincrby(key string, val interface{}, score int64) (int64, error) {
	b, err := Serialize(val)
	if err != nil {
		return 0, err
	}
	return c.GetSrv(key).Zincrby(key, b, score)
}

func (c *Cluster) Zscore(key string, val interface{}) (int64, error) {
	b, err := Serialize(val)
	if err != nil {
		return 0, err
	}
	return c.GetSrv(key).Zscore(key, b)
}

func (c *Cluster) Zexists(key string, val interface{}) (bool, error) {
	client := c.GetSrv(key)
	db, err := ssdb.Connect(client.conn_addr, client.conn_port)
	defer db.Close()
	if err != nil {
		return false, err
	}
	b, err := Serialize(val)
	if err != nil {
		return false, err
	}
	resp, err := db.Do("zexists", key, b)
	if err != nil {
		return false, err
	}
	if resp[0] == "ok" && resp[1] == "1" {
		return true, nil
	}
	return false, nil
}
