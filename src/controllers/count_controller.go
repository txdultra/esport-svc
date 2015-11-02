package controllers

import (
	"libs"
	"libs/groups"
	"libs/lives"
	"libs/message"
	"libs/passport"
	"libs/share"
	"libs/stat"
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

//var fsc libs.IEventCounter = &passport.FriendShips{}

//var sc libs.IEventCounter = &share.ShareVodSubcurs{}
//var pmc libs.IEventCounter = lives.NewProgramNoticeService()
//var sn libs.IEventCounter = share.NewShareNoticeService()
//var sm libs.IEventCounter = share.NewShareMsgService()
//var sharec libs.IEventCounter = &share.Shares{}
//var vc libs.IEventCounter = &vod.Vods{}
//var gms libs.IEventCounter = groups.NewGroupService(groups.GetDefaultCfg())

var count_mods = []string{
	passport.FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME,
	share.SHARE_MY_VOD_SUBSCRIPTIONS_COUNT_MODNAME,
	lives.MEMBER_PROGRAM_NEWNOTICE_COUNT_MODNAME,
	share.SHARE_MY_NEW_NOTICES_COUNT_MODNAME,
	share.SHARE_LASTNEW_MSG_MODNAME,
	share.SHARE_MEMBER_ATMSG_NEW_COUNT_MODNAME,
	vod.VOD_MSG_BOX_COUNT_MODNAME,
	groups.GROUP_MSG_NEW_COUNT_MODNAME,
	message.SYS_MSG_NEWS_COUNT_MODENAME,
}

func countLimitDigit(count int64) int {
	if count > 99 {
		return 99
	}
	return int(count)
}

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

	ucs := stat.UCGetCounts(uid, count_mods)
	//	follwers := fsc.NewEventCount(uid)
	//	scs := sc.NewEventCount(uid)
	//	//需要修改
	//	msgs := sharec.NewEventCount(uid)
	//	vods := vc.NewEventCount(uid)
	//	pms := pmc.NewEventCount(uid)
	//	sns := sn.NewEventCount(uid)
	//	lsm := sm.NewEventCount(uid)
	//	gcm := gms.NewEventCount(uid)
	//	sysm := message.NewEventCount(uid, vars.MSG_TYPE_SYS)

	cs := &outobjs.OutMemberNewCount{}

	for module, c := range ucs {
		switch module {
		case passport.FRIENDSHIP_NEW_FOLLOWER_COUNT_MODNAME:
			cs.NewFollowers = countLimitDigit(c)
			break
		case share.SHARE_MY_VOD_SUBSCRIPTIONS_COUNT_MODNAME:
			cs.NewSubscrs = countLimitDigit(c)
			break
		case share.SHARE_MY_NEW_NOTICES_COUNT_MODNAME:
			cs.NewShareNotices = countLimitDigit(c)
			break
		case lives.MEMBER_PROGRAM_NEWNOTICE_COUNT_MODNAME:
			cs.NewLiveSubscrs = countLimitDigit(c)
			break
		case share.SHARE_LASTNEW_MSG_MODNAME:
			cs.LastNewShareMsgs = countLimitDigit(c)
			break
		case share.SHARE_MEMBER_ATMSG_NEW_COUNT_MODNAME:
			cs.ShareMsgs = countLimitDigit(c)
			break
		case vod.VOD_MSG_NEW_COUNT_MODNAME:
			cs.VodMsgs = countLimitDigit(c)
			break
		case groups.GROUP_MSG_NEW_COUNT_MODNAME:
			cs.NewGroupMsgs = countLimitDigit(c)
			break
		case message.SYS_MSG_NEWS_COUNT_MODENAME:
			cs.NewMsgs = countLimitDigit(c)
			break
		}
	}

	//	if follwers > 99 {
	//		follwers = 99
	//	}
	//	if scs > 99 {
	//		scs = 99
	//	}
	//	if msgs > 99 {
	//		msgs = 99
	//	}
	//	if vods > 99 {
	//		vods = 99
	//	}
	//	if pms > 99 {
	//		pms = 99
	//	}
	//	if sns > 99 {
	//		sns = 99
	//	}
	//	if lsm > 99 {
	//		lsm = 99
	//	}
	//	if gcm > 99 {
	//		gcm = 99
	//	}

	//	cs := &outobjs.OutMemberNewCount{
	//		NewFollowers:     follwers,
	//		NewSubscrs:       scs,
	//		NewMsgs:          sysm,
	//		ShareMsgs:        msgs,
	//		VodMsgs:          vods,
	//		NewLiveSubscrs:   pms,
	//		NewShareNotices:  sns,
	//		LastNewShareMsgs: lsm,
	//		NewGroupMsgs:     gcm,
	//	}
	c.Json(cs)
}
