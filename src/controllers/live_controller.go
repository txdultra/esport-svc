package controllers

import (
	"fmt"
	"libs"
	//"libs/collect"
	"libs/lives"
	"libs/passport"
	"libs/reptile"
	"libs/search"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"

	//"github.com/sortutil"
)

// 直播 API
type LiveController struct {
	BaseController
}

func (c *LiveController) Prepare() {
	c.BaseController.Prepare()
}

func (c *LiveController) URLMapping() {
	c.Mapping("PerGet", c.PerGet)
	c.Mapping("PerGets", c.PerGets)
	c.Mapping("ChannelGet", c.ChannelGet)
	c.Mapping("ChannelStreams", c.ChannelStreams)
	c.Mapping("Programs", c.Programs)
	c.Mapping("RemindSingle", c.RemindSingle)
	c.Mapping("RemoveRemindSingle", c.RemoveRemindSingle)
	c.Mapping("NewUUID", c.NewUUID)
	c.Mapping("JoinLive", c.JoinLive)
	c.Mapping("SearchProgram", c.SearchProgram)
	c.Mapping("SubsReminded", c.SubsReminded)
	c.Mapping("LiveStreamCallback", c.LiveStreamCallback)
	c.Mapping("SubProgramsByChannel", c.SubProgramsByChannel)
}

// @Title 获取个人直播信息
// @Description 个人直播信息
// @Success 200 {object} outobjs.OutPersonalLive
// @router /personal/:id([0-9]+) [get]
func (c *LiveController) PerGet() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("live_parameter", "L1101", "参数错误", ""))
		return
	}
	live := &lives.LivePers{}
	per := live.Get(id)
	if per == nil {
		c.Json(libs.NewError("live_personal_notexist", "L1102", "不存在个人直播频道", ""))
		return
	}
	c.Json(outobjs.GetOutPersonalLive(per))
}

// @Title 获取个人直播列表
// @Description 个人直播列表(返回数组)
// @Param   access_token path    string  false  "access_token(验证通过,将进行用户关注过滤,搜索时无效)"
// @Param   query 		 path	 string  false  "查询字符串"
// @Param   game_ids     path    string  false  "游戏分类列表(示例:1,2,3),空为全部"
// @Param   size     path    int  true  "一页数据(默认20)"
// @Param   page     path    int  false  "页数(默认1)"
// @Param   exclude_ids path   string  false  "忽略指定直播ids(示例:1,2,3)"
// @Param   match_mode path  string   false  "搜索模式(all,any,phrase,boolean,extended,fullscan,extended2),默认为any"
// @Param   only_c path  bool   false  "只显示关注过的个人直播"
// @Success 200 {object} outobjs.OutPersonalLiveList
// @router /personal/list [get]
func (c *LiveController) PerGets() {
	uid := c.CurrentUid()
	query, _ := utils.UrlDecode(c.GetString("query"))
	size, _ := c.GetInt("size")
	page, _ := c.GetInt("page")
	game_ids := c.GetString("game_ids")
	match_mode := c.GetString("match_mode")
	exclude_id_str := c.GetString("exclude_ids")
	only_c, _ := c.GetBool("only_c")
	//timestamp, _ := c.GetInt64("t")

	//t := time.Now()
	//if timestamp > 0 {
	//	t = time.Unix(timestamp, 0)
	//}
	//2分钟间隔缓存
	query_cache_key := fmt.Sprintf("front_fast_cache.personallive.query:words_%s_gameids_%s_p_%d_s_%d_mode_%s_ex_%s_uid_%d_onlyc_%t",
		query, game_ids, page, size, match_mode, exclude_id_str, uid, only_c)
	c_obj := utils.GetLocalFastExpriesTimePartCache(query_cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}

	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	spgids := strings.Split(game_ids, ",")
	exclude_id_sps := strings.Split(exclude_id_str, ",")
	fgameIds := []uint64{}
	for _, id := range spgids {
		_id, err := strconv.Atoi(id)
		if err == nil && _id > 0 {
			fgameIds = append(fgameIds, uint64(_id))
		}
	}

	exclude_ids := []uint64{}
	c_livings := []*lives.LivePerson{}
	onlyc_livings := []*lives.LivePerson{}
	lps := &lives.LivePers{}
	if len(query) == 0 && uid > 0 {
		lives := lps.Livings()
		for _, lp := range lives {
			if !lp.Enabled {
				continue
			}
			ofGameIds := lps.GetOfGames(lp.Id)
			//比较方法
			isIn := func(src []int, tar []uint64) bool {
				for _, s := range src {
					for _, t := range tar {
						if int64(s) == int64(t) {
							return true
						}
					}
				}
				return false
			}
			if passport.IsFriend(uid, lp.Uid) {
				if isIn(ofGameIds, fgameIds) {
					c_livings = append(c_livings, lp)
					exclude_ids = append(exclude_ids, uint64(lp.Id))
				}
				onlyc_livings = append(onlyc_livings, lp)
			}
		}
	}
	//只显示关注过的
	if only_c {
		outs := []*outobjs.OutPersonalLive{}
		for _, v := range onlyc_livings {
			outs = append(outs, outobjs.GetOutPersonalLive(v))
		}
		//sortutil.AscByField(outs, "SortField")
		var totals = len(outs)
		out_oc := &outobjs.OutPersonalLiveList{
			Total:       totals,
			Lives:       outs,
			TotalPage:   utils.TotalPages(int(totals), int(size)),
			CurrentPage: int(page),
			Size:        int(size),
			Time:        time.Now().Unix(),
		}
		c.Json(out_oc)
		return
	}

	//参数id过滤
	for _, id := range exclude_id_sps {
		_id, err := strconv.Atoi(id)
		if err == nil && _id > 0 {
			exclude_ids = append(exclude_ids, uint64(_id))
		}
	}
	var filters []search.SearchFilter
	if len(fgameIds) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "game_ids",
			Values:  fgameIds,
			Exclude: false,
		})
	}
	if len(exclude_ids) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "pid",
			Values:  exclude_ids,
			Exclude: true,
		})
	}

	if len(match_mode) == 0 {
		match_mode = "any"
	}
	if len(query) == 0 {
		match_mode = "all"
	}
	total, list := lps.Query(query, int(page), int(size), match_mode, filters, nil)
	outs := []*outobjs.OutPersonalLive{}
	scs := 0
	if page == 1 { //第一页时加载收藏的个人直播
		c_out_livings := make([]*outobjs.OutPersonalLive, len(c_livings), len(c_livings))
		for i, v := range c_livings {
			c_out_livings[i] = outobjs.GetOutPersonalLive(v)
		}
		//sortutil.DescByField(c_out_livings, "Onlines")
		outs = append(outs, c_out_livings...)
		scs = len(c_livings) + len(list)
	} else {
		scs = size
	}

	for _, v := range list {
		outs = append(outs, outobjs.GetOutPersonalLive(v))
	}
	//sortutil.AscByField(outs, "SortField")
	//再按在线人数排序

	var totals = total + len(c_livings)
	out_p := outobjs.OutPersonalLiveList{
		Total:       totals,
		Lives:       outs,
		TotalPage:   utils.TotalPages(int(totals), int(scs)),
		CurrentPage: int(page),
		Size:        scs,
		Time:        time.Now().Unix(),
	}
	utils.SetLocalFastExpriesTimePartCache(1*time.Minute, query_cache_key, out_p)
	c.Json(out_p)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//机构直播
///////////////////////////////////////////////////////////////////////////////////////////////////

// @Title 获取机构直播信息
// @Description 机构直播信息
// @Param access_token  path   string  false  "access_token"
// @Success 200 {object} outobjs.OutLiveChannel
// @router /channel/:id([0-9]+) [get]
func (c *LiveController) ChannelGet() {
	uid := c.CurrentUid()
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("live_parameter", "L1201", "参数错误", ""))
		return
	}
	live := &lives.LiveOrgs{}
	channel := live.GetChannel(id)
	if channel == nil {
		c.Json(libs.NewError("live_org_notexist", "L1202", "不存在机构直播频道", ""))
		return
	}

	c.Json(outobjs.GetOutLiveChannel(channel, uid))
}

// @Title 获取机构直播流地址
// @Description 机构直播流地址(默认流放在第一个)
// @Param access_token  path   string  false  "access_token"
// @Success 200 {object} outobjs.OutChannelV1
// @router /channel/streams/:id([0-9]+) [get]
func (c *LiveController) ChannelStreams() {
	uid := c.CurrentUid()
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("live_parameter", "L1301", "参数错误", ""))
		return
	}
	live := &lives.LiveOrgs{}
	channel := live.GetChannel(id)
	if channel == nil {
		c.Json(libs.NewError("live_org_notexist", "L1302", "不存在机构直播频道", ""))
		return
	}
	streams := live.GetStreams(id)
	out_streams := []*outobjs.OutChannelStream{}
	out_channel := outobjs.GetOutLiveChannel(channel, uid)
	for _, stream := range streams {
		oos := &outobjs.OutChannelStream{
			Id:        stream.Id,
			Rep:       c.streamName(stream.Rep),
			StreamUrl: stream.StreamUrl,
			IsDefault: stream.Default,
			RepMethod: reptile.LiveRepMethod(stream.Rep),
		}
		out_streams = append(out_streams, oos)
	}
	c.Json(&outobjs.OutChannelV1{
		Channel: out_channel,
		Streams: out_streams,
	})
}

func (c *LiveController) streamName(rep reptile.REP_SUPPORT) string {
	switch rep {
	case reptile.REP_SUPPORT_17173:
		return "直播流2"
	case reptile.REP_SUPPORT_DOUYU:
		return "直播流3"
	case reptile.REP_SUPPORT_QQ:
		return "直播流1"
	case reptile.REP_SUPPORT_ZHANQI:
		return "直播流4"
	case reptile.REP_SUPPORT_HUOMAO:
		return "直播流5"
	case reptile.REP_SUPPORT_TWITCH:
		return "直播流6"
	case reptile.REP_SUPPORT_QQOPEN:
		return "直播流7"
	case reptile.REP_SUPPORT_CC163:
		return "直播流8"
	case reptile.REP_SUPPORT_PANDATV:
		return "直播流9"
	case reptile.REP_SUPPORT_MLGTV:
		return "直播流10"
	case reptile.REP_SUPPORT_HITBOX:
		return "直播流11"
	default:
		return "未知流"
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//节目单
///////////////////////////////////////////////////////////////////////////////////////////////////
// @Title 搜索节目单
// @Description 搜索节目单,返回数组
// @Param access_token  path   string  false  "access_token"
// @Param query path string  false  "关键字"
// @Param size path string  false  "返回数量(默认10)"
// @Success 200 {object} outobjs.OutProgramObj
// @router /channel/programs/search [get]
func (c *LiveController) SearchProgram() {
	uid := c.CurrentUid()
	query, _ := utils.UrlDecode(c.GetString("query"))
	size, _ := c.GetInt("size")
	match_mode := c.GetString("match_mode")
	if size <= 0 {
		size = 10
	}
	if len(match_mode) == 0 {
		match_mode = "any"
	}
	if len(query) == 0 {
		match_mode = "all"
	}
	t := time.Now()
	year_hours := time.Duration(24 * 365)
	filterRanges := []search.FilterRangeInt{
		//search.FilterRangeInt{
		//	Attr:    "stime",
		//	Min:     uint64(t.Unix()),
		//	Max:     uint64(t.Add(year_hours * time.Hour).Unix()),
		//	Exclude: false,
		//},
		search.FilterRangeInt{
			Attr:    "etime",
			Min:     uint64(t.Add(5 * time.Minute).Unix()),
			Max:     uint64(t.Add(year_hours * time.Hour).Unix()),
			Exclude: false,
		},
	}
	lps := &lives.LivePrograms{}
	lsps := &lives.LiveSubPrograms{}
	_, ps := lps.Query(query, 1, int(size), match_mode, nil, filterRanges)
	today_lives := make(map[int64]*lives.LiveProgram)
	for _, p := range ps {
		//查询正在直播的二级节目单
		if p.Date.Month() == t.Month() && p.Date.Day() == t.Day() {
			if lsps.IsLiving(p.Id) {
				today_lives[p.Id] = &p
			}
		} else {
			//否则直接退出,后面的节目单都不是当天的
			break
		}
	}
	out_ps := []*outobjs.OutProgramObj{}
	for _, p := range today_lives {
		out_pm := outobjs.ConvertOutProgramObj(p, uid)
		out_ps = append(out_ps, out_pm)
	}
	for _, p := range ps {
		if _, ok := today_lives[p.Id]; !ok {
			out_pm := outobjs.ConvertOutProgramObj(&p, uid)
			out_ps = append(out_ps, out_pm)
		}
	}
	c.Json(out_ps)
}

// @Title 节目单信息
// @Description 节目单信息
// @Param access_token  path   string  false  "access_token"
// @Param id path int  false  "节目单id(主节目单)"
// @Success 200 {object} outobjs.OutProgramObj
// @router /channel/program [get]
func (c *LiveController) Program() {
	uid := c.CurrentUid()
	program_id, _ := c.GetInt64("id")
	if program_id <= 0 {
		c.Json(libs.NewError("live_program_id_zero", "L1401", "id参数不能小于或等于0", ""))
		return
	}
	lps := &lives.LivePrograms{}
	program := lps.Get(program_id)
	if program == nil {
		c.Json(libs.NewError("live_program_notexist", "L1402", "节目单不存在", ""))
		return
	}
	out_p := outobjs.ConvertOutProgramObj(program, uid)
	c.Json(*out_p)
}

// @Title 节目单列表
// @Description 当前日期和后n天 (n > 0 && n < 10), 游戏分类过滤由客户端完成,返回数组
// @Param 	access_token  path   string  false  "access_token"
// @Param   n path int  false  "后几天,默认4"
// @Success 200 {object} outobjs.OutProgramV1
// @router /channel/programs [get]
func (c *LiveController) Programs() {
	uid := c.CurrentUid()
	n, _ := c.GetInt("n")
	if n <= 0 {
		n = 4
	}
	if n > 10 {
		n = 10
	}

	dates := []time.Time{time.Now()}
	for i := 1; i <= int(n); i++ {
		dur := utils.StrToDuration(strconv.Itoa(i*24) + "h")
		dates = append(dates, time.Now().Add(dur))
	}
	pms := &lives.LivePrograms{}
	out_ps := []*outobjs.OutProgramV1{}
	//组装输出数据
	for _, d := range dates {
		dps := pms.Gets(d)
		out_pms := []*outobjs.OutProgramObj{}
		for _, dp := range dps {
			out_pm := outobjs.ConvertOutProgramObj(dp, uid)
			out_pms = append(out_pms, out_pm)
		}
		//直播->待播->点播 按顺序
		out_order_pms := []*outobjs.OutProgramObj{}
		//直播
		for _, pm := range out_pms {
			if pm.IsLiving {
				out_order_pms = append(out_order_pms, pm)
			}
		}
		//点播
		for _, pm := range out_pms {
			if !pm.IsExpired && !pm.IsLiving {
				out_order_pms = append(out_order_pms, pm)
			}
		}
		//过期
		for _, pm := range out_pms {
			if pm.IsExpired {
				out_order_pms = append(out_order_pms, pm)
			}
		}

		out_ps = append(out_ps, &outobjs.OutProgramV1{
			Date:     d,
			Year:     d.Year(),
			YearStr:  utils.IntToCh(d.Year()),
			Month:    int(d.Month()),
			MonthStr: utils.MonthToCh(d),
			Day:      d.Day(),
			DayStr:   fmt.Sprintf("%.2d", d.Day()),
			Week:     utils.WeekName2(d, "zh"),
			Programs: out_order_pms,
		})
	}
	c.Json(out_ps)
}

// @Title 频道子节目单列表
// @Description 频道子节目单列表(返回数组)
// @Param  access_token  path   string  false  "access_token"
// @Param  cid   path    int  true  "频道id(channel_id)"
// @Success 200 {object} outobjs.OutSubProgramObj
// @router /channel/subprograms_by_channel [get]
func (c *LiveController) SubProgramsByChannel() {
	uid := c.CurrentUid()
	cid, _ := c.GetInt64("cid")
	lps := &lives.LivePrograms{}
	lsps := &lives.LiveSubPrograms{}

	out_subps := []*outobjs.OutSubProgramObj{}
	if cid <= 0 {
		c.Json(out_subps)
		return
	}

	var pids []int64
	t := time.Now()
	//1分钟间隔缓存
	query_cache_key := fmt.Sprintf("front_fast_cache.live.subprogram_by_channel:cid:%d", cid)
	c_obj := utils.GetLocalFastTimePartCache(t, query_cache_key, utils.CACHE_INTERVAL_TIME_TYPE_MINUTE)
	if c_obj != nil {
		pids = c_obj.([]int64)
	} else {
		pids = lps.GetsByDate(">=", t, cid)
		utils.SetLocalFastTimePartCache(t, query_cache_key, utils.CACHE_INTERVAL_TIME_TYPE_MINUTE, pids)
	}

	for _, pid := range pids {
		p := lps.Get(pid)
		if p == nil {
			continue
		}
		subs := lsps.Gets(pid)
		for _, sub := range subs {
			out_s := outobjs.ConvertOutSubProgramObj(sub, p, uid)
			if out_s != nil {
				out_subps = append(out_subps, out_s)
			}
		}
	}
	//sortutil.AscByField(out_subps, "StartTime")
	c.Json(out_subps)
}

//// @Title 节目单提醒(用于一组)
//// @Description 提交节目单提醒(用于一组)
//// @Param  access_token  path    string  true  "access_token"
//// @Param  pid   path    int  true  "一级节目单id"
//// @Param  spids  path string  true  "选中的二级节目单id数组(示例:1,2,3,4)"
//// @Success 200 成功编号error_code:REP000
//// @router /channel/programs/remind [post]
//func (c *LiveController) Remind() {
//	uid := c.CurrentUid()
//	if uid <= 0 {
//		c.Json(libs.NewError("member_setremind_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能设置提醒", ""))
//		return
//	}

//	pid, _ := c.GetInt("pid")
//	spidStrs := c.GetString("spids")
//	if pid <= 0 {
//		c.Json(libs.NewError("live_remind_pid_zero", "L1501", "pid参数不能为0", ""))
//		return
//	}
//	if len(strings.Trim(spidStrs, " ")) == 0 {
//		c.Json(libs.NewError("live_remind_spids_not_empty", "L1502", "spids不能为空", ""))
//		return
//	}
//	spids := []int64{}
//	for _, spid := range strings.Split(spidStrs, ",") {
//		_id, err := strconv.Atoi(spid)
//		if err != nil {
//			c.Json(libs.NewError("live_remind_illicit_parameter", "L1503", "非法参数", ""))
//			return
//		}
//		spids = append(spids, int64(_id))
//	}
//	noticer := lives.NewProgramNoticeService()
//	//保持用户提交的原始数据
//	err := noticer.SaveOrginalNotices(uid, pid, spids)
//	if err != nil {
//		c.Json(libs.NewError("live_remind_save_orginal", "L1506", "保存提交的原始数据时失败", ""))
//		return
//	}
//	filterStrategy := lives.FilterProgramNoticeStrategy(lives.PROGRAM_NOTICE_STRATEGY_SEQUENCE)
//	noticeIds := filterStrategy.CalculateEffectiveSubscriptionIds(pid, spids)
//	err = noticer.SubscribeNotice(uid, pid, noticeIds)
//	if err == nil {
//		c.Json(libs.NewError("live_remind_success", RESPONSE_SUCCESS, "已成功设置提醒", ""))
//		return
//	}
//	c.Json(libs.NewError("live_remind_fail", "L1505", "失败:"+err.Error(), ""))
//}

// @Title 节目单提醒(用于单个)
// @Description 提交节目单提醒(用于单个)
// @Param  access_token  path    string  true  "access_token"
// @Param  spid  path 	 int  true  "选中的二级节目单id"
// @Success 200 成功编号error_code:REP000
// @router /channel/programs/remind_single [post]
func (c *LiveController) RemindSingle() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_setremind_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能设置提醒", ""))
		return
	}

	spid, _ := c.GetInt64("spid")
	if spid <= 0 {
		c.Json(libs.NewError("live_remind_spid_illegality", "L1520", "spid参数非法", ""))
		return
	}

	lsp := &lives.LiveSubPrograms{}
	_sp := lsp.Get(spid)
	if _sp == nil {
		c.Json(libs.NewError("live_remind_spid_illegality", "L1521", "不存在对应的子菜单id", ""))
		return
	}

	noticer := lives.NewProgramNoticeService()
	//保持用户提交的原始数据
	err := noticer.SaveOrginalNoticeSingle(uid, _sp.ProgramId, spid)
	if err != nil {
		c.Json(libs.NewError("live_remind_save_orginal", "L1522", "保存提交的原始数据时失败", ""))
		return
	}
	err = noticer.SubscribeNoticeSingle(uid, _sp.ProgramId, spid)
	if err == nil {
		c.Json(libs.NewError("live_remind_success", RESPONSE_SUCCESS, "已成功设置提醒", ""))
		return
	}
	c.Json(libs.NewError("live_remind_fail", "L1523", "失败:"+err.Error(), ""))
}

// @Title 移除节目单提醒(用于单个)
// @Description 移除节目单提醒(用于单个)
// @Param  access_token  path    string  true  "access_token"
// @Param  spid  path 	 int  true  "选中的二级节目单id"
// @Success 200 成功编号error_code:REP000
// @router /channel/programs/remove_remind_single [delete]
func (c *LiveController) RemoveRemindSingle() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_rmvremind_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能设置提醒", ""))
		return
	}

	spid, _ := c.GetInt64("spid")
	if spid <= 0 {
		c.Json(libs.NewError("live_remind_spid_illegality", "L1530", "spid参数非法", ""))
		return
	}

	lsp := &lives.LiveSubPrograms{}
	_sp := lsp.Get(spid)
	if _sp == nil {
		c.Json(libs.NewError("live_remind_spid_illegality", "L1531", "不存在对应的子菜单id", ""))
		return
	}
	noticer := lives.NewProgramNoticeService()
	//保持用户提交的原始数据
	err := noticer.RemoveOrginalNoticeSingle(uid, _sp.ProgramId, spid)
	if err != nil {
		c.Json(libs.NewError("live_remind_save_orginal", "L1532", "保存提交的原始数据时失败", ""))
		return
	}
	err = noticer.RemoveSubscribeNoticeSingle(uid, _sp.ProgramId, spid)
	if err == nil {
		c.Json(libs.NewError("live_remind_success", RESPONSE_SUCCESS, "已成功移除提醒", ""))
		return
	}
	c.Json(libs.NewError("live_remind_fail", "L1535", "失败:"+err.Error(), ""))
}

// @Title 移除一组节目单提醒(批量)
// @Description 移除一组节目单提醒(批量)
// @Param  access_token  path    string  true  "access_token"
// @Param  spids  path 	 string  true  "选中的二级节目单ids,用,隔开"
// @Success 200 成功编号error_code:REP000
// @router /channel/programs/remove_reminds [delete]
func (c *LiveController) RemoveReminds() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_rmvremind_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能设置提醒", ""))
		return
	}

	spids_str := c.GetString("spids")
	spids_sp := strings.Split(spids_str, ",")
	spids := []int64{}
	for _, spidstr := range spids_sp {
		_id, err := strconv.ParseInt(spidstr, 10, 64)
		if err == nil {
			spids = append(spids, _id)
		}
	}
	if len(spids) == 0 {
		c.Json(libs.NewError("live_remind_spids_illegality", "L1540", "提交的参数格式错误", ""))
		return
	}
	lsp := &lives.LiveSubPrograms{}
	noticer := lives.NewProgramNoticeService()
	for _, spid := range spids {
		_sp := lsp.Get(spid)
		if _sp == nil {
			continue
		}
		//保持用户提交的原始数据
		noticer.RemoveOrginalNoticeSingle(uid, _sp.ProgramId, spid)
		noticer.RemoveSubscribeNoticeSingle(uid, _sp.ProgramId, spid)
	}
	c.Json(libs.NewError("live_remind_success", RESPONSE_SUCCESS, "已成功移除提醒", ""))
}

//// @Title 已订阅节目单提醒
//// @Description 已订阅节目单提醒,只查询订阅的某个节目单下的子节目单 返回[]int
//// @Param  access_token  path    string  true  "access_token"
//// @Param  pid   path    int  true  "一级节目单id"
//// @Success 200
//// @router /channel/programs/reminded [get]
//func (c *LiveController) Reminded() {
//	uid := c.CurrentUid()
//	if uid <= 0 {
//		c.Json(libs.NewError("member_getremind_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能获取提醒信息", ""))
//		return
//	}

//	pid, _ := c.GetInt("pid")
//	if pid <= 0 {
//		c.Json(libs.NewError("live_remind_pid_zero", "L1601", "pid参数不能为0", ""))
//		return
//	}
//	pns := lives.NewProgramNoticeService()
//	//data := pns.GetOrginalNotices(uid, pid)
//	//if data == nil {
//	//	c.Json(libs.NewError("live_remind_empty", "L1602", "未预定提醒", ""))
//	//}

//	//lsp := &lives.LiveSubPrograms{}
//	//spns := []outobjs.OutUserSubProgramNotice{}
//	//for _, refId := range data.RefIds {
//	//	spns = append(spns, outobjs.OutUserSubProgramNotice{
//	//		SubId:           refId,
//	//		IsLiving:        lsp.IsLiving(refId),
//	//		SubScribeLocked: lsp.IsLocked(refId),
//	//	})
//	//}
//	//out := outobjs.OutUserProgramNotice{
//	//	ProgramId:   data.EventId,
//	//	LastTime:    data.LastTime,
//	//	SubPrograms: spns,
//	//}

//	out := []int64{}
//	lsp := &lives.LiveSubPrograms{}
//	sps := lsp.Gets(pid)
//	for _, sp := range sps {
//		if pns.IsSubsuribed(uid, sp.Id) {
//			out = append(out, sp.Id)
//		}
//	}

//	c.Json(out)
//}

// @Title 已订阅节目单提醒
// @Description 查询订阅的将要开始的子节目单(自动清空计数;返回数组)
// @Param  access_token  path  string  true  "access_token"
// @Success 200 {object} outobjs.OutSubProgramObj
// @router /channel/programs/subs_reminded [get]
func (c *LiveController) SubsReminded() {
	uid := c.CurrentUid()
	if uid <= 0 {
		c.Json(libs.NewError("member_subremined_premission_denied", UNAUTHORIZED_CODE, "未登录状态不能获取节目单提醒信息", ""))
		return
	}
	noticer := lives.NewProgramNoticeService()
	_, subs := noticer.GetSubscribeNotices(uid, time.Now(), 1, 100)
	out_spms := []*outobjs.OutSubProgramObj{}
	lps := &lives.LivePrograms{}
	for _, sub := range subs {
		program := lps.Get(sub.ProgramId)
		if program == nil {
			continue
		}
		out_spm := outobjs.ConvertOutSubProgramObj(sub, program, uid)
		if out_spm != nil && !out_spm.SubScribeLocked {
			out_spms = append(out_spms, out_spm)
		}
	}
	//清空计数器
	noticer.ResetEventCount(uid)
	c.Json(out_spms)
}

// @Title 获取在线唯一标识符
// @Description 获取在线唯一标识符(error_code=REP000表示获取成功,error_content为获取的标识符)
// @Success 200 error_code=REP000表示获取成功,error_content为获取的标识符
// @router /online/new_uuid [get]
func (c *LiveController) NewUUID() {
	oc := &lives.OnlineCounter{}
	uuid := oc.NewUUID()
	c.Json(libs.NewError("live_online_new_uuid", RESPONSE_SUCCESS, uuid, ""))
}

// @Title 加入观看频道人员列表中
// @Description 加入到观看频道人员列表中(每5分钟轮调一次,离开时调用leave接口,error_code=REP000表示获取成功)
// @Param  access_token  path    string  false  "access_token"
// @Param  uuid  path   string  true  "uuid(通过接口获得)"
// @Param  live_type   path    int  true  "直播类型(机构=1,个人=2)"
// @Param  cid   path    int  true  "直播频道编号"
// @Success 200 error_code=REP000表示获取成功
// @router /online/join [post]
func (c *LiveController) JoinLive() {
	uid := c.CurrentUid()
	uuid := c.GetString("uuid")
	cid, _ := c.GetInt("cid")
	lt, _ := c.GetInt("live_type")
	if len(uuid) != 32 {
		c.Json(libs.NewError("live_online_join_fail", "L1701", "uuid格式错误", ""))
		return
	}
	if lt <= 0 || lt > 2 {
		c.Json(libs.NewError("live_online_join_fail", "L1701", "live_type参数错误", ""))
		return
	}
	if cid <= 0 {
		c.Json(libs.NewError("live_online_join_fail", "L1701", "cid参数错误", ""))
		return
	}
	live_type := lives.LIVE_TYPE(lt)
	oc := &lives.OnlineCounter{}
	err := oc.JoinChannel(live_type, int(cid), uid, uuid)
	if err == nil {
		c.Json(libs.NewError("live_online_join_succ", RESPONSE_SUCCESS, "加入成功", ""))
		return
	}
	c.Json(libs.NewError("live_online_join_fail", "L1702", err.Error(), ""))
	return
}

// @Title 离开观看频道
// @Description 离开观看频道(error_code=REP000表示获取成功)
// @Param  uuid  path   string  true  "uuid(通过接口获得)"
// @Param  live_type   path    int  true  "直播类型(机构=1,个人=2)"
// @Param  cid   path    int  true  "直播频道编号"
// @Success 200 error_code=REP000表示获取成功
// @router /online/leave [post]
func (c *LiveController) LeaveLive() {
	uuid := c.GetString("uuid")
	cid, _ := c.GetInt("cid")
	lt, _ := c.GetInt("live_type")
	if len(uuid) != 32 {
		c.Json(libs.NewError("live_online_leave_fail", "L1701", "uuid格式错误", ""))
		return
	}
	if lt <= 0 || lt > 2 {
		c.Json(libs.NewError("live_online_leave_fail", "L1701", "live_type参数错误", ""))
		return
	}
	if cid <= 0 {
		c.Json(libs.NewError("live_online_leave_fail", "L1701", "cid参数错误", ""))
		return
	}
	live_type := lives.LIVE_TYPE(lt)
	oc := &lives.OnlineCounter{}
	err := oc.LeaveChannel(live_type, int(cid), uuid)
	if err == nil {
		c.Json(libs.NewError("live_online_leave_succ", RESPONSE_SUCCESS, "已离开", ""))
		return
	}
	c.Json(libs.NewError("live_online_leave_fail", "L1702", err.Error(), ""))
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//抓取帮助方法
///////////////////////////////////////////////////////////////////////////////////////////////////

// @Title 抓取流程接口
// @Description 抓取流程接口(completed完成抓取,exit无法播放,reply请重新发起流程)
// @Param  id  path   int  true  "频道编号(个人为频道id,机构为流id)"
// @Param  live_type   path    int  true  "直播类型(机构=1,个人=2)"
// @Param  content   path    string  false  "客户端获取内容(首次提交时无须填写)"
// @Param  state   path    string  false  "状态(首次提交时无须填写)"
// @Success 200 {object} outobjs.OutClientProxyRepCallback
// @router /stream/callback [post]
func (c *LiveController) LiveStreamCallback() {
	pid, _ := c.GetInt64("id")
	ltype, _ := c.GetInt("live_type")
	content, _ := utils.UrlDecode(c.GetString("content"))
	state := c.GetString("state")
	if pid <= 0 {
		c.Json(libs.NewError("live_stream_rep_fail", "L1801", "id格式错误", ""))
		return
	}
	if ltype != 1 && ltype != 2 {
		c.Json(libs.NewError("live_stream_rep_fail", "L1802", "未指定直播类型", ""))
		return
	}
	if ltype == 2 {
		live := &lives.LivePers{}
		per := live.Get(pid)
		if per == nil {
			c.Json(libs.NewError("live_stream_rep_fail", "L1803", "频道不存在", ""))
			return
		}
		var parameter string
		if len(state) == 0 {
			parameter = per.ReptileUrl
		} else {
			parameter = content
		}
		_url, _state, err := lives.ClientProxyReptile(per.Rep, parameter, state)
		out_err := ""
		if err != nil {
			out_err = err.Error()
		}
		out_callback := &outobjs.OutClientProxyRepCallback{
			ClientUrl: _url,
			State:     _state,
			Error:     out_err,
		}
		c.Json(out_callback)
		return
	}
	if ltype == 1 {
		live := &lives.LiveOrgs{}
		stream := live.GetStream(pid)
		if stream == nil {
			c.Json(libs.NewError("live_stream_rep_fail", "L1805", "频道流不存在", ""))
			return
		}
		var parameter string
		if len(state) == 0 {
			parameter = stream.ReptileUrl
		} else {
			parameter = content
		}
		_url, _state, err := lives.ClientProxyReptile(stream.Rep, parameter, state)
		out_err := ""
		if err != nil {
			out_err = err.Error()
		}
		out_callback := &outobjs.OutClientProxyRepCallback{
			ClientUrl: _url,
			State:     _state,
			Error:     out_err,
		}
		c.Json(out_callback)
		return
	}
	c.Json(libs.NewError("live_stream_rep_fail", "L1809", "无支持的方法", ""))
}
