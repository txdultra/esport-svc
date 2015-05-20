package groups

import (
	"dbs"
	"fmt"
	"libs/dlock"
	"logs"
	"reflect"
	"time"
	"utils"
	"utils/ssdb"
)

func runGroupCountUpdateService() {
	go gpService()
}

func gpService() {
	for {
		//分布式同步锁控制更新lock
		lock := dlock.NewDistributedLock()
		locker, err := lock.Lock("group_update_count_lock")
		if err != nil {
			logs.Errorf("group service update count get lock fail:%s", err.Error())
			return
		}

		gmaps := make(map[int64]map[GP_PROPERTY]int)
		o := dbs.NewOrm(db_aliasname)
		cclient := ssdb.New(use_ssdb_group_db)
		for _, gp := range GP_PROPERTY_ALL {
			_key := fmt.Sprintf(count_action_set, string(gp))
			kss, _ := cclient.ZscanKS(_key, -1<<10, 1<<10, 1<<10, reflect.TypeOf(int64(0)), reflect.TypeOf(int64(0)))
			for _, ks := range kss {
				groupid := *(ks.Key.(*int64))
				incrs := *(ks.Score.(*int64))
				if pv, ok := gmaps[groupid]; ok {
					pv[gp] = int(incrs)
				} else {
					gmaps[groupid] = make(map[GP_PROPERTY]int)
					gmaps[groupid][gp] = int(incrs)
				}
			}
			cclient.Zclear(_key)
		}

		if len(gmaps) > 0 {
			updates := []*Group{}
			g := &Group{}
			gs := NewGroupService(GetDefaultCfg())
			p, _ := o.Raw("UPDATE " + g.TableName() + " SET members = members + ?,threads = threads + ? WHERE id = ?").Prepare()
			for groupid, gps := range gmaps {
				group := gs.Get(groupid)
				if group == nil {
					continue
				}
				addMembers := 0
				addThreads := 0
				for k, v := range gps {
					switch k {
					case GP_PROPERTY_MEMBERS:
						group.Members += v
						addMembers = v
					case GP_PROPERTY_THREADS:
						group.Threads += v
						addThreads = v
					default:
					}
				}
				p.Exec(addMembers, addThreads, groupid)
				updates = append(updates, group)
			}
			p.Close()
			gs.UpdateSearchEngineAttr(updates)
			cache := utils.GetCache()
			for _, g = range updates {
				cache.Replace(gs.GetCacheKey(g.Id), *g, 1*time.Hour)
			}
		}
		//解锁
		locker.Unlock()
		//15秒间隔
		time.Sleep(15 * time.Second)
	}
}
