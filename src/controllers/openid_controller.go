//error prefix "P3"
package controllers

import (
	"libs"
	"libs/passport"
	"outobjs"
	"time"
	"utils"
)

// OpenID api
type OpenIDController struct {
	BaseController
}

func (c *OpenIDController) Prepare() {
	c.BaseController.Prepare()
}

func (c *OpenIDController) URLMapping() {
	c.Mapping("Login", c.Login)
}

// @Title OpenID登录
// @Description OpenID登录(成功登录后请调用set_pushid接口绑定第三方推送账号)
// @Param   open_token   path   string  true  "open_token,第三方登录时获取的token"
// @Param   openid   path  string  true  "用户第三方的登录的唯一标识码"
// @Param   op   path  int  true  "第三方标识(腾讯=1,微信=2)"
// @Param   mobile_iden   path  string  false  "用户手机唯一序列号"
// @Param   longi   path  float  false  "经度"
// @Param   lati   path  float  false  "纬度"
// @Success 200
// @router /login [post]
func (oc *OpenIDController) Login() {
	open_token := oc.GetString("open_token")
	openid := oc.GetString("openid")
	op, err := oc.GetInt("op")
	mobile_iden := oc.GetString("mobile_iden")
	ip := oc.Ctx.Input.IP()
	longiTude, _ := oc.GetFloat("longi")
	latiTude, _ := oc.GetFloat("lati")

	if err != nil {
		oc.Json(libs.NewError("openidlogin_op_empty", "P3001", "op cant empty", ""))
		return
	}

	longi := float32(longiTude)
	lati := float32(latiTude)

	var third passport.IThirdAuth
	//opp := byte(op)
	switch passport.OPENID_MARK(op) {
	case passport.OPENID_MARK_QQ:
		third = new(passport.QQThridOpenID)
	case passport.OPENID_MARK_WEIXIN:
		third = new(passport.WeixinThridOpenID)
	default:
		third = nil
	}

	if third == nil {
		oc.Json(libs.NewError("openidlogin_op_provider_notexist", "P3002", "op not exist", ""))
		return
	}
	//third_openid := openid
	third_info, err := third.GetUserInfo(open_token, openid)
	if err != nil {
		oc.Json(libs.NewError("openidlogin_sys_fail", "P3010", "system fail", ""))
		return
	}

	//查找绑定用户
	open_provider := passport.NewOpenIDAuth()
	auth, err := open_provider.Get(openid, passport.OPENID_MARK(op))
	//新用户情况
	if auth == nil {
		member := passport.Member{}
		//openid模式自动创建用户名
		_username := utils.RandomStrings(6) //third_info.NickName + "_" + utils.RandomStrings(3)
		if len(third_info.NickName) == 0 {
			_username = utils.RandomStrings(6)
		}
		member.Password = openid
		member.UserName = _username
		member.RegMode = passport.MEMBER_REGISTER_OPENID
		if len(mobile_iden) == 0 {
			member.MobileIdentifier = ""
		} else {
			member.MobileIdentifier = mobile_iden
		}
		//member.MemberIdentifier = utils.RandomStrings(32)
		//member.CreateTime = time.Now().Unix()
		member.CreateIP = utils.IpToInt(ip)
		member_provider := passport.NewMemberProvider()
		uid, err := member_provider.Create(member, longi, lati)
		if err != nil {
			oc.Json(err)
			return
		}
		auth = new(passport.OpenIDOAuth)
		auth.AuthType = byte(op)
		auth.Uid = uid
		auth.AuthIdentifier = openid
		auth.AuthEmail, auth.AuthEmailVerified = "", ""
		auth.AuthPreferredUserName, auth.AuthProviderName = "", ""
		auth.AuthToken = open_token
		auth.AuthDate = time.Now()
		auth.OpenIDUid = openid
		auth_provider := passport.NewOpenIDAuth()
		_, auth_err := auth_provider.Create(auth)
		if auth_err != nil {
			oc.Json(libs.NewError("openidlogin_openauth_create_fail", "P3004", "internal error", ""))
			return
		}

		access_service := passport.NewAccessTokenService()
		acc_token, terr := access_service.GetAccessTokenNoCode(ip, "", "", uid)
		if terr == nil {
			//oc.Json(oc.returnAccessToken(acc_token))
			oc.Json(outobjs.GetAccessToken(acc_token))
			return
		}
		oc.Json(terr)
		return

	} else { //已有用户
		uid := auth.Uid
		access_service := passport.NewAccessTokenService()
		acc_token, terr := access_service.GetAccessTokenNoCode(ip, "", "", uid)
		if terr == nil {
			//oc.Json(oc.returnAccessToken(acc_token))
			oc.Json(outobjs.GetAccessToken(acc_token))
			return
		}
		oc.Json(terr)
		return
	}
}
