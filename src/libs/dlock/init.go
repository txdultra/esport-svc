package dlock

import (
	"github.com/astaxie/beego"
	"strings"
)

var zookeeper_addrs []string

func init() {
	//zookeeper 加载节点集合
	zookeeper_addr_strs := strings.Split(beego.AppConfig.String("zookeeper.addrs"), ",")
	zookeeper_addrs = []string{}
	for _, zk := range zookeeper_addr_strs {
		if len(zk) > 0 {
			zookeeper_addrs = append(zookeeper_addrs, zk)
		}
	}
}
