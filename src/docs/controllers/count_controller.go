package controllers

import (
	"libs"
	"libs/lives"
	"libs/message"
	"libs/passport"
	"libs/share"
	"outobjs"
)

// 计数 API
type CountController struct {
	BaseController
}

func (c *CountController) Prepare() {
	c.BaseController.Prepare()
}

func (c *CountController) URLMapping() {
	c.Mapping("MemberCounts", c.MemberCounts)
}

var fsc *passport.FriendShips = &passport.FriendShips{}
var sc *share.ShareVodSubcurs = &share.ShareVodSubcurs{}
var pmc *lives.ProgramNoticer = lives.NewProgramNoticeService()

// @Title 用户所有计数器
// @Description 用户所有计数器
// @Param   access_token   path  string  true  "access_token"
// @Success 200  {object} outobjs.OutMemberNewCount
// @router /all [get]
func (c *CountController) MemberCounts() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_count_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询", ""))
		return
	}
	follwers := fsc.NewFollowers(uid)
	scs := sc.NewEventCounts(uid)
	msgs := message.NewMsgCount(uid)
	pms := pmc.NewMemberNoticeCount(uid)

	if follwers > 99 {
		follwers = 99
	}
	if scs > 99 {
		scs = 99
	}
	if msgs > 99 {
		msgs = 99
	}
	if pms > 99 {
		pms = 99
	}

	cs := &outobjs.OutMemberNewCount{
		NewFollowers:   follwers,
		NewSubscrs:     scs,
		NewMsg:         msgs,
		NewLiveSubscrs: pms,
	}
	c.Json(cs)
}
