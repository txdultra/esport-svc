package groupcache

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/golang/groupcache"
	"strings"
	"sync"
)

type groupCfg struct {
	Name    string
	Size    int64
	Existed bool
}

var groups map[string]*groupCfg = make(map[string]*groupCfg)
var locker *sync.Mutex = new(sync.Mutex)

func init() {
	me := beego.AppConfig.String("groupcache.me")
	grouphosts := beego.AppConfig.String("groupcache.groups")
	groupnames := beego.AppConfig.String("groupcache.groupnames")

	group_hosts := strings.Split(grouphosts, ",")
	group_names := strings.Split(groupnames, ",")

	if len(me) == 0 {
		return
	}
	peers := groupcache.NewHTTPPool(me)
	peers.Set(group_hosts...)

	for _, name := range group_names {
		if len(name) == 0 {
			continue
		}
		size := beego.AppConfig.DefaultInt64("groupcache."+name+".size", 268435456) //256mb
		if _, ok := groups[name]; ok {
			continue
		}
		groups[name] = &groupCfg{
			Name:    name,
			Size:    size,
			Existed: false,
		}
	}
}

func RegisterGroupCache(groupName string, getFunc groupcache.GetterFunc) error {
	locker.Lock()
	defer locker.Unlock()
	if cfg, ok := groups[groupName]; ok {
		if cfg.Existed {
			return fmt.Errorf(groupName + " registered...")
		}
		groupcache.NewGroup(cfg.Name, cfg.Size, getFunc)
		cfg.Existed = true
	} else {
		return fmt.Errorf(groupName + " config not exist")
	}
	return nil
}

func GetGroupCache(groupName string) *groupcache.Group {
	return groupcache.GetGroup(groupName)
}
