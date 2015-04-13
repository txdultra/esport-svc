package controllers

import (
	"libs/passport"
	"time"
	"utils"
)

type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Uid         int64  `json:"uid"`
}

type UcController struct {
	BaseController
	atkService *passport.AccessTokenService
}

func (c *UcController) Prepare() {
	c.BaseController.Prepare()
	c.atkService = passport.NewAccessTokenService()
}

// @router /create [post]
func (c *UcController) Create() {
	username := c.GetString("username")
	password := c.GetString("password")
	email := c.GetString("email")
	client_secret := c.GetString("client_secret")
	client_id := c.GetString("client_id")
	ip := c.Ctx.Input.IP()

	member := passport.Member{}
	member.UserName = username
	//member.NickName = username
	member.Email = email
	member.Password = password
	member.Salt, member.MobileIdentifier, member.MemberIdentifier = "", "", ""
	member.CreateTime = time.Now().Unix()
	member.CreateIP = utils.IpToInt(ip)

	uid, err := passport.NewMemberProvider().Create(member, 0, 0) //c.uProvider.Create(member, 0, 0)
	if err != nil {
		c.Json(err)
		return
	}
	if uid > 0 {
		code := c.atkService.NewCode()
		at, err := c.atkService.GetAccessToken(code, ip, client_id, client_secret, uid)
		if err == nil {
			result := &AccessTokenResult{at.AccessToken, at.ExpiresIn, uid}
			c.Json(result)
			return
		}
	}
	c.Json(err)
}
