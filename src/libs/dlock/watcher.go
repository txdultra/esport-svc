package dlock

import (
	"logs"
	"time"

	"github.com/go-zookeeper/zk"
)

type Watcher struct{}

func NewWatcher() *Watcher {
	return &Watcher{}
}

func (w *Watcher) Write(path string, data []byte) error {
	conn, _, err := zk.Connect(zookeeper_addrs, 5*time.Second)
	if err != nil {
		logs.Errorf("connect zookeeper server fail:%+v", err)
		return err
	}
	defer conn.Close()
	existed, _, err := conn.Exists(path)
	if err != nil {
		logs.Errorf("Exists returned fail:%+v", err)
		return err
	}
	if !existed {
		if _path, err := conn.Create(path, []byte(data), 0, zk.WorldACL(zk.PermAll)); err != nil {
			logs.Errorf("Create returned error: %+v", err)
			return err
		} else if _path != path {
			logs.Errorf("Create returned different path '%s' != '%s'", _path, path)
			return err
		}
	} else {
		if _, err := conn.Set(path, []byte(data), -1); err != nil {
			logs.Errorf("Set value returned error: %+v", err)
			return err
		}
	}
	return nil
}

func (w *Watcher) RegisterWatcher(path string, callback func(data []byte)) {
	conn, _, err := zk.Connect(zookeeper_addrs, 5*time.Second)
	if err != nil {
		logs.Errorf("connect zookeeper server fail:%+v", err)
		return
	}

	go func() {
		defer conn.Close()
		for {
			data, _, event, err := conn.GetW(path)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			select {
			case ev := <-event:
				logs.Errorf("Received event: %+v", ev)
				callback(data)
			case <-time.After(5 * time.Second):
				continue
			}
		}
	}()
}
