package reptile

import "github.com/astaxie/beego"

var douyu_api_host string

func init() {
	douyu_api_host = beego.AppConfig.String("reptile.douyu.api.host")
	if len(douyu_api_host) == 0 {
		panic("斗鱼远程服务地址未设置")
	}
}
