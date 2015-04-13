package controllers

import (
	//dbs "../dbs"
	//. "../models"
	//"fmt"
	"outobjs"
)

var out_api_mod_hosts []*outobjs.OutApiModHost = nil

type HomeController struct {
	BaseController
}

func (c *HomeController) URLMapping() {
	c.Mapping("Get", c.Get)
	c.Mapping("ApiUrls", c.ApiUrls)
}

func (c *HomeController) Get() {
	c.Ctx.WriteString("ntv mobile api v1.0 ")
}

func (c *HomeController) ApiUrls() {
	if out_api_mod_hosts == nil {
		out_api_mod_hosts = []*outobjs.OutApiModHost{
			&outobjs.OutApiModHost{
				ModName: "com",
				BaseUrl: host_maps["com"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "vod",
				BaseUrl: host_maps["vod"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "live",
				BaseUrl: host_maps["live"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "comment",
				BaseUrl: host_maps["comment"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "collect",
				BaseUrl: host_maps["collect"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "user",
				BaseUrl: host_maps["user"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "friendships",
				BaseUrl: host_maps["friendships"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "share",
				BaseUrl: host_maps["share"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "msg",
				BaseUrl: host_maps["msg"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "count",
				BaseUrl: host_maps["count"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "other",
				BaseUrl: host_maps["other"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "file",
				BaseUrl: host_maps["file"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
			&outobjs.OutApiModHost{
				ModName: "img",
				BaseUrl: host_maps["img"], //"http://192.168.0.33:8080",
				Version: "v1",
			},
		}
	}
	c.Json(out_api_mod_hosts)
}
