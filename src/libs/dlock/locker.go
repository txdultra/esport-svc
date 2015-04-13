package dlock

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"logs"
	"strings"
	"time"
)

const (
	DISTRIBUTED_LOCK_PATH = "/distributed_locks"
)

type ILocker interface {
	Lock(src string) (ILocker, error)
	Unlock() error
}

type DistributedLock struct {
	zkL *zk.Lock
	zkC *zk.Conn
}

func NewDistributedLock() *DistributedLock {
	return &DistributedLock{}
}

func (d DistributedLock) Lock(src string) (ILocker, error) {
	if len(src) == 0 {
		return nil, fmt.Errorf("未提供锁资源")
	}
	conn, _, err := zk.Connect(zookeeper_addrs, 5*time.Second)
	if err != nil {
		logs.Errorf("connect zookeeper server fail:%s", err.Error())
		return nil, err
	}

	var lockPath string
	if strings.HasPrefix(src, "/") {
		lockPath = DISTRIBUTED_LOCK_PATH + "/" + src
	} else {
		lockPath = DISTRIBUTED_LOCK_PATH + src
	}
	acl := zk.WorldACL(zk.PermAll)
	zklock := zk.NewLock(conn, lockPath, acl)

	if err = zklock.Lock(); err != nil {
		defer conn.Close()
		logs.Errorf("zookeeper service open lock fail:%s", err.Error())
		return nil, err
	}

	d.zkL = zklock
	d.zkC = conn
	return d, nil
}

func (d DistributedLock) Unlock() error {
	if d.zkL == nil || d.zkC == nil {
		return fmt.Errorf("未进行lock操作,不能解锁")
	}
	err := d.zkL.Unlock()
	if err != nil {
		return err
	}
	d.zkC.Close()
	return nil
}
