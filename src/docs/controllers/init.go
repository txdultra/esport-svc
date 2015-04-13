package controllers

import (
	"libs"

	"github.com/astaxie/beego"
)

var host_maps map[string]string = make(map[string]string)
var file_storage libs.IFileStorage

func init() {
	file_storage = libs.NewFileStorage()
	//api host
	host_maps["com"] = beego.AppConfig.String("api.host.com")
	host_maps["vod"] = beego.AppConfig.String("api.host.vod")
	host_maps["live"] = beego.AppConfig.String("api.host.live")
	host_maps["comment"] = beego.AppConfig.String("api.host.comment")
	host_maps["collect"] = beego.AppConfig.String("api.host.collect")
	host_maps["user"] = beego.AppConfig.String("api.host.user")
	host_maps["friendships"] = beego.AppConfig.String("api.host.friendships")
	host_maps["share"] = beego.AppConfig.String("api.host.share")
	host_maps["msg"] = beego.AppConfig.String("api.host.msg")
	host_maps["count"] = beego.AppConfig.String("api.host.count")
	host_maps["other"] = beego.AppConfig.String("api.host.other")
	host_maps["file"] = beego.AppConfig.String("api.host.file")
	host_maps["img"] = beego.AppConfig.String("api.host.img")
}
