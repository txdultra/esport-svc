package controllers

import (

	//"fmt"
	"libs"
	"libs/passport"
	"utils"

	"github.com/astaxie/beego"
	//log "github.com/cihub/seelog"
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

func (c *BaseController) Prepare() {
	//log.Tracef("ip:%s url:%s", c.Ctx.Input.IP(), this.Ctx.Input.Uri())
}

func (c *BaseController) Json(data interface{}) {
	c.Data["json"] = data
	c.ServeJson()
	//c.StopRun()
}

func (c *BaseController) WriteString(str string) {
	c.Ctx.WriteString(str)
}

func (c *BaseController) CurrentUid() int64 {
	access_token := c.GetString("access_token")
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

func (c *BaseController) CurrentUser() *passport.Member {
	currentUid := c.CurrentUid()
	if currentUid <= 0 {
		return nil
	}
	mp := passport.NewMemberProvider()
	member := mp.Get(currentUid)
	return member
}

func (c *BaseController) ValidateAccessToken() (*passport.AccessToken, *libs.Error) {
	access_token := c.GetString("access_token")
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
