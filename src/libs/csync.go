package libs

import (
	"sync"
	"time"
	"utils"
)

var singnalLocker *sync.RWMutex = new(sync.RWMutex)
var singnals map[string]chan string = make(map[string]chan string)

//线程同步等待
//完成是发送completed
func ChanSync(key string, run func()) (bool, <-chan string) {
	singnalLocker.Lock()
	defer singnalLocker.Unlock()
	c, ok := singnals[key]
	if !ok {
		c = make(chan string)
		singnals[key] = c
		go func() {
			log := func(interface{}) {}
			finally := func() {
				c <- "completed"
				singnalLocker.Lock()
				defer singnalLocker.Unlock()
				delete(singnals, key)
			}
			utils.Try(run, log, finally)
		}()
	}
	return ok, c
}

func ChanSyncExist(key string) bool {
	singnalLocker.RLock()
	defer singnalLocker.RUnlock()
	_, ok := singnals[key]
	return ok
}

func ChanSyncWait(key string) {
	singnalLocker.RLock()
	c, ok := singnals[key]
	singnalLocker.RUnlock()
	if ok {
		select {
		case <-time.After(30 * time.Second): //超过2分钟自动退出
		case <-c:
		}
	}
}
