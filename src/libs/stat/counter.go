package stat

import (
	//"fmt"
	"sync"
)

//计数器接口
type ICounter interface {
	DoC(id int64, n int, term string)
	GetC(id int64, term string) int
}

var counter_mods map[string]ICounter = make(map[string]ICounter)

func RegisterCounter(mod string, counter ICounter) {
	if _, ok := counter_mods[mod]; !ok {
		counter_mods[mod] = counter
	}
}

func GetCounter(mod string) ICounter {
	if _c, ok := counter_mods[mod]; ok {
		return _c
	}
	panic("未注册" + mod + "计数器")
}

//计数骰子
var s_counters map[string]int = make(map[string]int)
var s_counterLocker *sync.Mutex = new(sync.Mutex)

func NewCountHplper() *CountHplper {
	return &CountHplper{}
}

type CountHplper struct{}

func (c *CountHplper) IncreaseCount(key string, n int, dice func(i int) bool, succ_callback func(i int) bool) {
	s_counterLocker.Lock()
	defer s_counterLocker.Unlock()
	nums, ok := s_counters[key]
	if !ok {
		s_counters[key] = n
		return
	}
	nums += n
	s_counters[key] = nums
	if dice(nums) {
		_n := nums
		go func() {
			succ_callback(_n)
		}()
		s_counters[key] = 0
	}
}

func (c *CountHplper) ResetCount(key string) {
	s_counterLocker.Lock()
	defer s_counterLocker.Unlock()
	s_counters[key] = 0
}
