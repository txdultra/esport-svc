package hook

import "sync"

type IHook interface {
	Do(event string, args ...interface{})
}

var hooks map[string]map[string]IHook = make(map[string]map[string]IHook)
var lock *sync.RWMutex = new(sync.RWMutex)

func RegisterHook(event string, name string, hook IHook) {
	if hook == nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	var mHooks map[string]IHook
	if mhs, ok := hooks[event]; !ok {
		mHooks = make(map[string]IHook)
		hooks[event] = mHooks
	} else {
		mHooks = mhs
	}
	if _, ok := mHooks[name]; ok {
		panic(event + "事件下已存在相同名称的hook")
		return
	}
	mHooks[name] = hook
}

func Do(event string, args ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	if subhooks, ok := hooks[event]; ok {
		for _, hook := range subhooks {
			hook.Do(event, args...)
		}
	}
}
