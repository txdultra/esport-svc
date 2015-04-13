// @APIVersion 1.0.0
// @Title 电竞圈手机API
// @Description 根地址http://ip地址:端口/{version},version根据不同的app版本调用不同api
// @Contact tangxd@neotv.cn
package routers

import (
	"controllers"
	"controllers/admincp"

	"github.com/astaxie/beego"
)

func init() {
	//beego.EnableAdmin = true
	//beego.AdminHttpAddr = "localhost"
	//beego.AdminHttpPort = 9091

	//beego.Router("/", &controllers.WebController{}, "*:Home")
	beego.Router("/api_urls", &controllers.HomeController{}, "*:ApiUrls")

	beego.Include(
		&controllers.WebController{},
	)

	//file
	//if open_file {
	ns1 := beego.NewNamespace("v1",
		beego.NSNamespace("/com",
			beego.NSInclude(
				&controllers.CommonController{},
			),
		),
		beego.NSNamespace("/vod",
			beego.NSInclude(
				&controllers.VideoController{},
			),
		),
		beego.NSNamespace("/comment",
			beego.NSInclude(
				&controllers.CommentController{},
			),
		),
		beego.NSNamespace("/live",
			beego.NSInclude(
				&controllers.LiveController{},
			),
		),
		beego.NSNamespace("/collect",
			beego.NSInclude(
				&controllers.CollectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.MemberController{},
			),
		),
		beego.NSNamespace("/user_task",
			beego.NSInclude(
				&controllers.UserTaskController{},
			),
		),
		beego.NSNamespace("/shop",
			beego.NSInclude(
				&controllers.ShopController{},
			),
		),
		beego.NSNamespace("/openid",
			beego.NSInclude(
				&controllers.OpenIDController{},
			),
		),
		beego.NSNamespace("/friendships",
			beego.NSInclude(
				&controllers.FriendShipsController{},
			),
		),
		beego.NSNamespace("/share",
			beego.NSInclude(
				&controllers.ShareController{},
			),
		),
		beego.NSNamespace("/msg",
			beego.NSInclude(
				&controllers.MessageController{},
			),
		),
		beego.NSNamespace("/count",
			beego.NSInclude(
				&controllers.CountController{},
			),
		),
		beego.NSNamespace("/feedback",
			beego.NSInclude(
				&controllers.FeedbackController{},
			),
		),
		beego.NSNamespace("/file",
			beego.NSInclude(
				&controllers.FileController{},
			),
		),
		beego.NSNamespace("/img",
			beego.NSInclude(
				&controllers.ImageController{},
			),
		),
		////admin api
		//beego.NSNamespace("admin/cache",
		//	beego.NSInclude(
		//		&admincp.CacheCPController{},
		//	),
		//),
		//beego.NSNamespace("admin/common",
		//	beego.NSInclude(
		//		&admincp.CommonCPController{},
		//	),
		//),
		//beego.NSNamespace("admin/vod",
		//	beego.NSInclude(
		//		&admincp.VodCPController{},
		//	),
		//),
	)
	//admin api
	ns2 := beego.NewNamespace("admin",
		beego.NSNamespace("/auth_cp",
			beego.NSInclude(
				&admincp.AuthCPController{},
			),
		),
		beego.NSNamespace("/cache_cp",
			beego.NSInclude(
				&admincp.CacheCPController{},
			),
		),
		beego.NSNamespace("/com_cp",
			beego.NSInclude(
				&admincp.CommonCPController{},
			),
		),
		beego.NSNamespace("/vod_cp",
			beego.NSInclude(
				&admincp.VodCPController{},
			),
		),
		beego.NSNamespace("/live_cp",
			beego.NSInclude(
				&admincp.LiveCPController{},
			),
		),
		beego.NSNamespace("/member_cp",
			beego.NSInclude(
				&admincp.MemberCPController{},
			),
		),
		beego.NSNamespace("/share_cp",
			beego.NSInclude(
				&admincp.ShareCPController{},
			),
		),
		beego.NSNamespace("/feedback_cp",
			beego.NSInclude(
				&admincp.FeedbackCPController{},
			),
		),
		beego.NSNamespace("/utask_cp",
			beego.NSInclude(
				&admincp.UTaskCPController{},
			),
		),
		beego.NSNamespace("/ushop_cp",
			beego.NSInclude(
				&admincp.ShopCPController{},
			),
		),
	)
	beego.AddNamespace(ns1)
	beego.AddNamespace(ns2)
}
