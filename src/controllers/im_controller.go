package controllers

import (
	"fmt"
	"outobjs"
	"strings"
	"time"
	"utils"
)

// IM API
type IMController struct {
	BaseController
}

func (c *IMController) Prepare() {
	c.BaseController.Prepare()
}

func (c *IMController) URLMapping() {
}

const (
	IM_CKEY   = "im_%s_%s"
	IM_DOMAIN = "neotv.cn"
)

// @Title 加入聊天室所需配置
// @Description 加入聊天室所需配置
// @Param   access_token  path  string  false  "access_token"
// @Success 200  {object} outobjs.ChatUserCfg
// @router /join_usr_cfg [get]
func (c *IMController) JoinUserCfg() {
	usr := c.CurrentUser()
	rstr := utils.RandomStrings(6)
	cuc := &outobjs.ChatUserCfg{
		NickName: fmt.Sprintf("匿名%s", rstr),
	}
	if usr != nil {
		if len(usr.NickName) > 0 {
			cuc.NickName = usr.NickName
		}
		cuc.IsVip = usr.Certified
		cuc.IsOC = usr.OfficialCertified
	}
	cuc.Pwd = utils.RandomStrings(10)
	cuc.UserName = rstr
	cuc.Domain = IM_DOMAIN
	cuc.Port = 5222
	cuc.Host = "10.10.20.199"
	cuc.AllowMsg = true
	cuc.FontColor = 0xffffff

	cache := utils.GetCache()
	ckey := strings.ToLower(fmt.Sprintf(IM_CKEY, cuc.Domain, rstr))
	cache.Set(ckey, cuc.Pwd, 5*time.Minute)
	c.Json(cuc)
}

// @Title ejabberd调用接口
// @Description ejabberd调用接口
// @Success 200
// @router /get_password [get]
func (c *IMController) EjabberdGetPassword() {
	user_name := c.GetString("user")
	server := c.GetString("server")
	ckey := strings.ToLower(fmt.Sprintf(IM_CKEY, server, user_name))
	cache := utils.GetCache()
	pwd := utils.RandomStrings(10)
	cache.Get(ckey, &pwd)
	c.WriteString(pwd)
}
