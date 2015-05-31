package dlock

import (
	"strings"

	"github.com/astaxie/beego"
)

var zookeeper_addrs []string

func init() {
	//zookeeper 加载节点集合
	zkaddrs := beego.AppConfig.String("zookeeper.addrs")
	//	if len(zkaddrs) == 0 {
	//		panic("未配置zookeeper.addrs")
	//	}
	zookeeper_addr_strs := strings.Split(zkaddrs, ",")
	zookeeper_addrs = []string{}
	for _, zk := range zookeeper_addr_strs {
		if len(zk) > 0 {
			zookeeper_addrs = append(zookeeper_addrs, zk)
		}
	}
}
