package controllers

import (
	//"fmt"
	"libs"
	"libs/passport"
)

type AuthorizeController struct {
	BaseController
	Access_Token *passport.AccessToken
}

func (c *AuthorizeController) Prepare() {
	c.BaseController.Prepare()

	token, err := c.ValidateAccessToken()
	if err != nil {
		out_err := libs.NewError("unauthorized", UNAUTHORIZED_CODE, "您未进行登录授权", "")
		c.Json(out_err)
		c.StopRun()
	}
	c.Access_Token = token
	//即将过期刷新过期时间
	dur := c.Access_Token.ExpireDurition()
	if int(dur.Seconds()) < passport.Authorization_access_token_expries_refresh {
		service := passport.NewAccessTokenService()
		service.RefreshAccessTokenExpries(c.Access_Token)
	}
}
