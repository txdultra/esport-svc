package controllers

import (
	"sync"
	"time"
	"utils/bitmap"
	//"fmt"
	"libs"
	"libs/hook"
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

	dayLoginEvent(token.Uid)

	c.Access_Token = token
	//即将过期刷新过期时间
	dur := c.Access_Token.ExpireDurition()
	if int(dur.Seconds()) < passport.Authorization_access_token_expries_refresh {
		service := passport.NewAccessTokenService()
		service.RefreshAccessTokenExpries(c.Access_Token)
	}
}

type elBitmap struct {
	lock *sync.Mutex
	bms  []*bitmap.RoaringBitmap
}

var el = &elBitmap{new(sync.Mutex), make([]*bitmap.RoaringBitmap, 2, 2)}

func dayLoginEvent(uid int64) {
	day := time.Now().Day()
	mod := day % 2
	if el.bms[mod] == nil {
		el.lock.Lock()
		if el.bms[mod] == nil {
			el.bms[mod] = bitmap.New()
			el.bms[(day+1)%2] = nil
		}
		el.lock.Unlock()
	}
	if !el.bms[mod].Contains(uint32(uid)) {
		el.bms[mod].Add(uint32(uid))
		hook.Do("everyday_login", uid, 1)
	}
}
