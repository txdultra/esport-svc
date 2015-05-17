package controllers

import (
	"libs"
	"libs/lives"
	"libs/passport"
	"libs/share"
	"libs/vod"
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

var fsc libs.IEventCounter = &passport.FriendShips{}
var sc libs.IEventCounter = &share.ShareVodSubcurs{}
var pmc libs.IEventCounter = lives.NewProgramNoticeService()
var sn libs.IEventCounter = share.NewShareNoticeService()
var sm libs.IEventCounter = share.NewShareMsgService()
var sharec libs.IEventCounter = &share.Shares{}
var vc libs.IEventCounter = &vod.Vods{}

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
	follwers := fsc.NewEventCount(uid)
	scs := sc.NewEventCount(uid)
	//需要修改
	msgs := sharec.NewEventCount(uid)
	vods := vc.NewEventCount(uid)
	pms := pmc.NewEventCount(uid)
	sns := sn.NewEventCount(uid)
	lsm := sm.NewEventCount(uid)

	if follwers > 99 {
		follwers = 99
	}
	if scs > 99 {
		scs = 99
	}
	if msgs > 99 {
		msgs = 99
	}
	if vods > 99 {
		vods = 99
	}
	if pms > 99 {
		pms = 99
	}
	if sns > 99 {
		sns = 99
	}
	if lsm > 99 {
		lsm = 99
	}

	cs := &outobjs.OutMemberNewCount{
		NewFollowers:     follwers,
		NewSubscrs:       scs,
		NewMsgs:          msgs,
		ShareMsgs:        msgs,
		VodMsgs:          vods,
		NewLiveSubscrs:   pms,
		NewShareNotices:  sns,
		LastNewShareMsgs: lsm,
	}
	c.Json(cs)
}
