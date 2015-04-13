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

func (this *AuthorizeController) Prepare() {
	this.BaseController.Prepare()

	token, err := this.ValidateAccessToken()
	if err != nil {
		out_err := libs.NewError("unauthorized", UNAUTHORIZED_CODE, "您未进行登录授权", "")
		this.Json(out_err)
		this.StopRun()
	}
	this.Access_Token = token
	//即将过期刷新过期时间
	dur := this.Access_Token.ExpireDurition()
	if int(dur.Seconds()) < passport.Authorization_access_token_expries_refresh {
		service := passport.NewAccessTokenService()
		service.RefreshAccessTokenExpries(this.Access_Token)
	}
}
