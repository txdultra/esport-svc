package admincp

import (
	"controllers"
	"fmt"
	"libs"
	"libs/passport"
	"outobjs"
	"utils"
)

// 后台权限模块 API
type AuthCPController struct {
	controllers.BaseController
}

const (
	MANAGE_PLAT_NEOTV_KEY = "34326H4L937zaR3nfE"
)

// @Title 获取管理者access_token
// @Description 获取管理者access_token
// @Param   plat_uid   path	string true  "管理平台uid"
// @Param   plat   path	string true  "平台编号"
// @Param   timestamp   path	string true  "时间戳"
// @Param   sign   path	string true  "sign"
// @Success 200 {object} outobjs.OutAccessToken
// @router /get_access_token [get]
func (c *AuthCPController) GetManagerAccessToken() {
	plat_uid, _ := c.GetInt64("plat_uid")
	timestamp := c.GetString("timestamp")
	sign := c.GetString("sign")
	plat := c.GetString("plat")
	ip := c.Ctx.Input.IP()

	vfySign := utils.Md5(fmt.Sprintf("%d_%s_%s_%s", plat_uid, plat, timestamp, MANAGE_PLAT_NEOTV_KEY))
	if vfySign != sign {
		c.Json(libs.NewError("admincp_manage_accesstoken_fail", "GM000_001", "sign验证失败", ""))
		return
	}
	pms := &passport.PlatManagers{}
	mm := pms.GetManagerMap(plat_uid, plat)
	if mm == nil {
		c.Json(libs.NewError("admincp_manage_accesstoken_fail", "GM000_002", "您的账号未被授权", ""))
		return
	}
	access_service := passport.NewAccessTokenService()
	acc_token, terr := access_service.GetAccessTokenNoCode(ip, "", "", mm.ToUid)
	if terr == nil {
		data := outobjs.OutAccessToken{
			AccessToken: acc_token.AccessToken,
			ExpiresIn:   acc_token.ExpiresIn,
			Uid:         acc_token.Uid,
		}
		c.Json(data)
		return
	}
	c.Json(terr)
}

// @Title access_token是否有效
// @Description access_token是否有效
// @Param   token   path	string true  "access_token"
// @Param   plat   path	string true  "平台编号"
// @Param   timestamp   path	string true  "时间戳"
// @Param   sign   path	string true  "sign"
// @Success 200
// @router /access_token/status [get]
func (c *AuthCPController) AccessTokenStatus() {
	access_token := c.GetString("token")
	timestamp := c.GetString("timestamp")
	sign := c.GetString("sign")
	plat := c.GetString("plat")
	vfySign := utils.Md5(fmt.Sprintf("%s_%s_%s_%s", access_token, plat, timestamp, MANAGE_PLAT_NEOTV_KEY))
	if vfySign != sign {
		c.Json(libs.NewError("admincp_manage_accesstoken_fail", "GM000_003", "sign验证失败", ""))
		return
	}
	access_service := passport.NewAccessTokenService()
	_, terr := access_service.GetTokenObj(access_token)
	if terr == nil {
		c.Ctx.WriteString("ok")
		return
	}
	c.Ctx.WriteString(terr.ErrorDescription)
}
