package controllers

import (

	//"fmt"
	"libs"
	"libs/passport"
	"utils"

	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
)

const (
	RESPONSE_SUCCESS       = "REP000"
	REPSONSE_FAIL          = "REP001"
	UNAUTHORIZED_CODE      = "SYS101" //未授权
	PREMISSION_DENIED_CODE = "SYS102" //权限不足
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	log.Tracef("ip:%s url:%s", this.Ctx.Input.IP(), this.Ctx.Input.Uri())
}

func (this *BaseController) Json(data interface{}) {
	this.Data["json"] = data
	this.ServeJson()
	//this.StopRun()
}

func (this *BaseController) WriteString(str string) {
	this.Ctx.WriteString(str)
}

func (this *BaseController) CurrentUid() int64 {
	access_token := this.GetString("access_token")
	if len(access_token) == 0 {
		return 0
	}
	_acc_tkn := utils.StripSQLInjection(access_token)
	if access_token != _acc_tkn {
		return 0
	}
	acc := passport.NewAccessTokenService()
	token, err := acc.GetTokenObj(access_token)
	if err == nil {
		return token.Uid
	}
	return 0
}

func (this *BaseController) CurrentUser() *passport.Member {
	currentUid := this.CurrentUid()
	if currentUid <= 0 {
		return nil
	}
	mp := passport.NewMemberProvider()
	member := mp.Get(currentUid)
	return member
}

func (this *BaseController) ValidateAccessToken() (*passport.AccessToken, *libs.Error) {
	access_token := this.GetString("access_token")
	if len(access_token) == 0 {
		return nil, libs.NewError("access_token_not_empty", "A1098", "access_token不能为空", "")
	}
	_acc_tkn := utils.StripSQLInjection(access_token)
	if access_token != _acc_tkn {
		return nil, libs.NewError("access_token_illegality", "A1091", "access_token非法", "")
	}
	service := passport.NewAccessTokenService()
	token, err := service.GetTokenObj(access_token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
