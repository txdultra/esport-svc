package cache

import (
	//"errors"
	//"fmt"
	"github.com/valyala/ybc/memcache"
	//"sync"
	"time"
)

type MemcachedCacheV2 struct {
	_defaultExpiration time.Duration
	_timeout           time.Duration
	rc                 *memcache.DistributedClient
	hosts              []string
}

func NewMemcachedCacheV2(hostList []string, defaultExpiration time.Duration, timeout time.Duration) *MemcachedCacheV2 {
	c := &MemcachedCacheV2{
		hosts:              hostList,
		_defaultExpiration: defaultExpiration,
		_timeout:           timeout,
	}
	c.init()
	return c
}

func (c *MemcachedCacheV2) init() {
	if c.rc != nil {
		c.rc.Stop()
	}
	dc := &memcache.DistributedClient{}
	//dc.StartStatic(c.hosts)
	dc.Start()
	for _, addr := range c.hosts {
		dc.AddServer(addr)
	}
	c.rc = dc
}

func (c *MemcachedCacheV2) Stop() {
	c.rc.Stop()
}

func (c *MemcachedCacheV2) Set(key string, value interface{}, expires time.Duration) error {
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	item := memcache.Item{
		Key:        []byte(key),
		Value:      b,
		Expiration: expires,
	}
	if err := c.rc.Set(&item); err != nil {
		return err
	}
	return nil
}

func (c *MemcachedCacheV2) Add(key string, value interface{}, expires time.Duration) error {
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	item := memcache.Item{
		Key:        []byte(key),
		Value:      b,
		Expiration: expires,
	}
	if err := c.rc.Add(&item); err != nil {
		return err
	}
	return nil
}

func (c *MemcachedCacheV2) Replace(key string, value interface{}, expires time.Duration) error {
	return c.Set(key, value, expires)
}

func (c *MemcachedCacheV2) Get(key string, ptrValue interface{}) error {
	item := &memcache.Item{
		Key: []byte(key),
	}
	err := c.rc.Get(item)
	if err != nil {
		return err
	}
	return Deserialize(item.Value, ptrValue)
}

func (c *MemcachedCacheV2) GetMulti(keys ...string) (Getter, error) {
	var items []memcache.Item
	for _, key := range keys {
		items = append(items, memcache.Item{
			Key: []byte(key),
		})
	}
	err := c.rc.GetMulti(items)
	if err != nil {
		return nil, err
	}
	return ItemMapGetterV2(items), nil
	//panic("Not Implemented")
}

func (c *MemcachedCacheV2) Delete(key string) error {
	return c.rc.Delete([]byte(key))
}

func (c *MemcachedCacheV2) Increment(key string, delta uint64) (newValue uint64, err error) {
	//return 0, nil
	panic("Not Implemented")
}

func (c *MemcachedCacheV2) Decrement(key string, delta uint64) (newValue uint64, err error) {
	//return 0, nil
	panic("NotImplemented")
}

func (c *MemcachedCacheV2) Flush() error {
	return c.rc.FlushAll()
}

// Implement a Getter on top of the returned item map.
type ItemMapGetterV2 []memcache.Item

func (g ItemMapGetterV2) Get(key string, ptrValue interface{}) error {
	for _, it := range g {
		if string(it.Key) == key {
			return Deserialize(it.Value, ptrValue)
		}
	}
	return ErrCacheMiss
}
