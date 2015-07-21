package admincp

import (
	"controllers"
	//"encoding/json"
	//"fmt"
	"libs"
	"libs/lives"
	"libs/reptile"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 直播管理 API
type LiveCPController struct {
	AdminController
	storage libs.IFileStorage
}

func (c *LiveCPController) Prepare() {
	c.AdminController.Prepare()
	c.storage = libs.NewFileStorage()
}

// @Title 抓取支持的平台
// @Success 200
// @router /live/rep_plats [get]
func (c *LiveCPController) SupportPlatforms() {
	outs := reptile.SupportReptilePlatforms()
	c.Json(outs)
}

// @Title 添加新个人直播
// @Description 添加新个人直播
// @Param   name   path	string true  "名称"
// @Param   img_id   path	int true  "图片id"
// @Param   uid   path	int true  "所属播主"
// @Param   desc   path	string true  "描述"
// @Param   rep_url   path	string true  "抓取地址"
// @Param   game_ids   path	string true  "所属游戏"
// @Param   do_rep   path	bool true  "立即抓取"
// @Param   show_online_min   path	int false  ""
// @Param   show_online_max   path	int false  ""
// @Success 200 {object} libs.Error
// @router /personal/add [post]
func (c *LiveCPController) LivePersonalAdd() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	img_id, _ := c.GetInt64("img_id")
	uid, _ := c.GetInt64("uid")
	desc, _ := utils.UrlDecode(c.GetString("desc"))
	rep_url, _ := utils.UrlDecode(c.GetString("rep_url"))
	do_rep, _ := c.GetBool("do_rep")
	online_min, _ := c.GetInt("show_online_min")
	online_max, _ := c.GetInt("show_online_max")
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_personal_live_add_fail", "GM015_001", "name不能为空", ""))
		return
	}
	if img_id <= 0 {
		c.Json(libs.NewError("admincp_personal_live_add_fail", "GM015_002", "必须提供图片", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("admincp_personal_live_add_fail", "GM015_003", "必须设置所属播主", ""))
		return
	}
	if len(rep_url) == 0 {
		c.Json(libs.NewError("admincp_personal_live_add_fail", "GM015_004", "必须提供抓取地址", ""))
		return
	}
	ofgames := strings.Split(c.GetString("game_ids"), ",")
	ofgameIds := []int{}
	for _, ofg := range ofgames {
		_id, err := strconv.Atoi(ofg)
		if err == nil {
			ofgameIds = append(ofgameIds, _id)
		}
	}

	var per = lives.LivePerson{}
	per.Name = name
	per.Img = img_id
	per.Uid = uid
	per.Des = desc
	per.ReptileUrl = rep_url
	per.ShowOnlineMin = online_min
	per.ShowOnlineMax = online_max
	liveSrv := &lives.LivePers{}
	_, err := liveSrv.Create(per, ofgameIds, do_rep)
	if err == nil {
		c.Json(libs.NewError("admincp_personal_live_add_succ", controllers.RESPONSE_SUCCESS, "个人直播添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_personal_live_add_fail", "GM015_005", "其他错误:"+err.Error(), ""))
}

// @Title 更新个人直播
// @Description 更新个人直播
// @Param   id   path	int true  "个人直播id"
// @Param   name   path	string true  "名称"
// @Param   img_id   path	int true  "图片id"
// @Param   uid   path	int true  "所属播主"
// @Param   desc   path	string true  "描述"
// @Param   rep_url   path	string true  "抓取地址"
// @Param   game_ids   path	string true  "所属游戏"
// @Param   enabled   path	bool false  "是否开启"
// @Param   show_online_min   path	int false  ""
// @Param   show_online_max   path	int false  ""
// @Success 200 {object} libs.Error
// @router /personal/update [post]
func (c *LiveCPController) LivePersonalUpate() {
	id, _ := c.GetInt64("id")
	name, _ := utils.UrlDecode(c.GetString("name"))
	img_id, _ := c.GetInt64("img_id")
	uid, _ := c.GetInt64("uid")
	desc, _ := utils.UrlDecode(c.GetString("desc"))
	rep_url, _ := utils.UrlDecode(c.GetString("rep_url"))
	online_min, _ := c.GetInt("show_online_min")
	online_max, _ := c.GetInt("show_online_max")

	if id <= 0 {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_010", "id非法", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_011", "name不能为空", ""))
		return
	}
	if img_id <= 0 {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_012", "必须提供图片", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_013", "必须设置所属播主", ""))
		return
	}
	if len(rep_url) == 0 {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_014", "必须提供抓取地址", ""))
		return
	}
	ofgames := strings.Split(c.GetString("game_ids"), ",")
	ofgameIds := []int{}
	for _, ofg := range ofgames {
		_id, err := strconv.Atoi(ofg)
		if err == nil {
			ofgameIds = append(ofgameIds, _id)
		}
	}

	liveSrv := &lives.LivePers{}
	per := liveSrv.Get(id)
	if per == nil {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_015", "id对应的个人直播不存在", ""))
		return
	}
	per.Name = name
	per.Img = img_id
	per.Uid = uid
	per.Des = desc
	per.ReptileUrl = rep_url
	per.ShowOnlineMin = online_min
	per.ShowOnlineMax = online_max
	enabled, err := c.GetBool("enabled")
	if err == nil {
		per.Enabled = enabled
	}

	err = liveSrv.Update(*per)
	if err != nil {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_016", "其他错误:"+err.Error(), ""))
		return
	}
	err = liveSrv.UpdateGames(id, ofgameIds)
	if err != nil {
		c.Json(libs.NewError("admincp_personal_live_update_fail", "GM015_017", "其他错误:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_personal_live_update_succ", controllers.RESPONSE_SUCCESS, "个人直播更新成功", ""))
}

// @Title 删除个人直播
// @Description 删除个人直播
// @Param   id   path	int true  "个人直播id"
// @Success 200 {object} libs.Error
// @router /personal/del [delete]
func (c *LiveCPController) LivePersonalDel() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_personal_live_del_fail", "GM015_020", "id非法", ""))
		return
	}
	liveSrv := &lives.LivePers{}
	err := liveSrv.Delete(id)
	if err != nil {
		c.Json(libs.NewError("admincp_personal_live_del_fail", "GM015_021", "其他错误:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_personal_live_del_succ", controllers.RESPONSE_SUCCESS, "个人直播已被删除", ""))
}

// @Title 个人直播列表
// @Description 个人直播列表
// @Param   query   path	string true  "title关键字"
// @Param   game_id   path	int true  "所属游戏"
// @Param   page   path	int true  "分页"
// @Param   size   path	int true  "页数量"
// @Success 200 {object} outobjs.OutPersonalLiveListForAdmin
// @router /personal/list [get]
func (c *LiveCPController) LivePersonalList() {
	query, _ := utils.UrlDecode(c.GetString("query"))
	game_id, _ := c.GetInt("game_id")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	liveSrv := &lives.LivePers{}
	total, pers := liveSrv.ListForAdmin(query, int(game_id), int(page), int(size))
	out_lives := []*outobjs.OutPersonalLiveForAdmin{}
	for _, p := range pers {
		out_l := c.toOutPersonalLiveAdmin(p)
		out_lives = append(out_lives, out_l)
	}
	c.Json(outobjs.OutPersonalLiveListForAdmin{
		Total:       total,
		TotalPage:   utils.TotalPages(int(total), int(size)),
		CurrentPage: int(page),
		Size:        int(size),
		Time:        time.Now().Unix(),
		Lives:       out_lives,
	})
}

// @Title 个人直播
// @Description 个人直播
// @Param   id   path	int true  "直播id""
// @Success 200 {object} outobjs.OutPersonalLiveForAdmin
// @router /personal/:id([0-9]+) [get]
func (c *LiveCPController) LivePersonalGet() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_personal_live_get_fail", "GM015_030", "id非法", ""))
		return
	}
	lps := &lives.LivePers{}
	pl := lps.Get(id)
	if pl == nil {
		c.Json(libs.NewError("admincp_personal_live_get_fail", "GM015_031", "频道不存在", ""))
		return
	}
	c.Json(c.toOutPersonalLiveAdmin(pl))
}

func (c *LiveCPController) toOutPersonalLiveAdmin(per *lives.LivePerson) *outobjs.OutPersonalLiveForAdmin {
	lps := &lives.LivePers{}
	lp_games := lps.GetOfGames(per.Id)
	ofGames := []*outobjs.OutGame{}
	for _, ofgid := range lp_games {
		ofg := outobjs.GetOutGameById(ofgid)
		if ofg != nil {
			ofGames = append(ofGames, ofg)
		}
	}
	oc := &lives.OnlineCounter{}
	return &outobjs.OutPersonalLiveForAdmin{
		Id:            per.Id,
		Title:         per.Name,
		Des:           per.Des,
		ImgId:         per.Img,
		ImgUrl:        c.storage.GetFileUrl(per.Img),
		Uid:           per.Uid,
		StreamUrl:     per.StreamUrl,
		Status:        per.LiveStatus,
		Member:        outobjs.GetOutMember(per.Uid, 0),
		Onlines:       oc.GetChannelCounts(lives.LIVE_TYPE_PERSONAL, int(per.Id)),
		OfGames:       ofGames,
		RepMethod:     reptile.LiveRepMethod(per.Rep), //RepMethod:  per.RepMethod,
		ReptileUrl:    per.ReptileUrl,
		ReptileDes:    per.ReptileDes,
		PostTime:      per.PostTime,
		Enabled:       per.Enabled,
		ShowOnlineMin: per.ShowOnlineMin,
		ShowOnlineMax: per.ShowOnlineMax,
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//机构直播
///////////////////////////////////////////////////////////////////////////////////////////////////

// @Title 机构直播频道列表
// @Description 机构直播频道列表
// @Success 200 {object} outobjs.OutLiveChannelForAdmin
// @router /org/channel/list [get]
func (c *LiveCPController) LiveChannelList() {
	lvs := &lives.LiveOrgs{}
	channels := lvs.GetChannels()
	outs := []*outobjs.OutLiveChannelForAdmin{}
	for _, channel := range channels {
		outs = append(outs, outobjs.GetOutLiveChannelForAdmin(channel, 0))
	}
	c.Json(outs)
}

// @Title 机构直播频道列表
// @Description 机构直播频道列表
// @Param   id   path	int true  "id"
// @Success 200 {object} outobjs.OutLiveChannelForAdmin
// @router /org/channel/:id([0-9]+) [get]
func (c *LiveCPController) LiveChannelGet() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_channel_get_fail", "GM016_201", "id非法", ""))
		return
	}
	lvs := &lives.LiveOrgs{}
	channel := lvs.GetChannel(id)
	if channel == nil {
		c.Json(libs.NewError("admincp_live_channel_update_fail", "GM016_202", "频道不存在", ""))
		return
	}
	c.Json(outobjs.GetOutLiveChannelForAdmin(channel, 0))
}

// @Title 添加机构直播频道
// @Description 添加机构直播频道列表
// @Param   title   path	string true  "标题"
// @Param   img_id  path	int true  "所属游戏"
// @Param   uid   path	int true  "所属主播"
// @Param   enabled   path	int false  "是否开启(默认开启)"
// @Success 200  {object} libs.Error
// @router /org/channel/add [post]
func (c *LiveCPController) LiveChannelAdd() {
	title, _ := utils.UrlDecode(c.GetString("title"))
	img_id, _ := c.GetInt64("img_id")
	uid, _ := c.GetInt64("uid")
	enabled, err := c.GetBool("enbaled")
	if err != nil {
		enabled = true
	}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_live_channel_add_fail", "GM016_001", "必须设置标题", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("admincp_live_channel_add_fail", "GM016_002", "必须设置所属播主", ""))
		return
	}
	channel := &lives.LiveChannel{
		Name:    title,
		Img:     img_id,
		Uid:     uid,
		Enabled: enabled,
	}
	lvs := &lives.LiveOrgs{}
	_, err = lvs.CreateChannel(channel)
	if err != nil {
		c.Json(libs.NewError("admincp_live_channel_add_fail", "GM016_003", "添加失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_channel_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
}

// @Title 更新机构直播频道
// @Description 更新机构直播频道
// @Param   id   path	int true  "id"
// @Param   title   path	string true  "标题"
// @Param   img_id  path	int true  "所属游戏"
// @Param   uid   path	int true  "所属主播"
// @Param   enabled   path	int false  "是否开启(默认开启)"
// @Success 200  {object} libs.Error
// @router /org/channel/update [post]
func (c *LiveCPController) LiveChannelUpdate() {
	id, _ := c.GetInt64("id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img_id, _ := c.GetInt64("img_id")
	uid, _ := c.GetInt64("uid")
	enabled, err := c.GetBool("enabled")
	if err != nil {
		enabled = true
	}
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_channel_update_fail", "GM016_010", "必须提供id", ""))
		return
	}
	if len(title) == 0 {
		c.Json(libs.NewError("admincp_live_channel_update_fail", "GM016_011", "必须设置标题", ""))
		return
	}
	if uid <= 0 {
		c.Json(libs.NewError("admincp_live_channel_update_fail", "GM016_012", "必须设置所属播主", ""))
		return
	}
	lvs := &lives.LiveOrgs{}
	channel := lvs.GetChannel(id)
	if channel == nil {
		c.Json(libs.NewError("admincp_live_channel_update_fail", "GM016_013", "频道不存在", ""))
		return
	}
	channel.Name = title
	channel.Img = img_id
	channel.Uid = uid
	channel.Enabled = enabled
	err = lvs.UpdateChannel(channel)
	if err != nil {
		c.Json(libs.NewError("admincp_live_channel_update_fail", "GM016_014", "更新失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_channel_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 机构直播频道播放流列表
// @Description 机构直播频道播放流列表
// @Param   channel_id   path	int true  "频道id"
// @Success 200 {object} outobjs.OutChannelStreamForAdmin
// @router /org/channel_stream/list [get]
func (c *LiveCPController) LiveChannelStreamList() {
	channel_id, _ := c.GetInt64("channel_id")
	lvs := &lives.LiveOrgs{}
	streams := lvs.GetStreams(channel_id)
	outs := []*outobjs.OutChannelStreamForAdmin{}
	for _, stream := range streams {
		outs = append(outs, outobjs.GetOutLiveChannelStreamForAdmin(stream))
	}
	c.Json(outs)
}

// @Title 添加机构直播频道播放流
// @Description 添加机构直播频道播放流
// @Param   channel_id   path	int true  "所属频道id"
// @Param   img_id   path	int true  "img_id"
// @Param   rep_url   path	string true  "抓取地址"
// @Param   is_def   path	bool true  "是否默认"
// @Param   enabled   path	bool false  "是否开启(默认开启)"
// @Param   allow_rep   path	bool false  "允许抓取"
// @Success 200 {object} libs.Error
// @router /org/channel_stream/add [post]
func (c *LiveCPController) LiveChannelStreamAdd() {
	channel_id, _ := c.GetInt64("channel_id")
	rep_url, _ := utils.UrlDecode(c.GetString("rep_url"))
	img_id, _ := c.GetInt64("img_id")
	is_def, _ := c.GetBool("is_def")
	enabled, err := c.GetBool("enabled")
	allow_rep, _ := c.GetBool("allow_rep")
	if err != nil {
		enabled = true
	}
	if channel_id <= 0 {
		c.Json(libs.NewError("admincp_live_channel_stream_add_fail", "GM016_020", "必须提供channel_id", ""))
		return
	}
	if len(rep_url) == 0 {
		c.Json(libs.NewError("admincp_live_channel_stream_add_fail", "GM016_021", "抓取地址必须设置", ""))
		return
	}
	lvs := &lives.LiveOrgs{}
	stream := &lives.LiveStream{}
	stream.ReptileUrl = rep_url
	stream.ChannelId = channel_id
	stream.Img = img_id
	stream.Default = is_def
	stream.Enabled = enabled
	stream.AllowRep = allow_rep
	_, err = lvs.CreateStream(stream)
	if err != nil {
		c.Json(libs.NewError("admincp_live_channel_stream_add_fail", "GM016_022", "新增失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_channel_stream_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
}

// @Title 更新机构直播频道播放流
// @Description 更新机构直播频道播放流
// @Param   id   path	int true  "流id"
// @Param   img_id   path	int true  "img_id"
// @Param   rep_url   path	string true  "抓取地址"
// @Param   is_def   path	bool true  "是否默认"
// @Param   enabled   path	bool false  "是否开启(默认开启)"
// @Param   allow_rep   path	bool false  "允许抓取"
// @Success 200 {object} libs.Error
// @router /org/channel_stream/update [post]
func (c *LiveCPController) LiveChannelStreamUpdate() {
	id, _ := c.GetInt64("id")
	rep_url, _ := utils.UrlDecode(c.GetString("rep_url"))
	img_id, _ := c.GetInt64("img_id")
	is_def, _ := c.GetBool("is_def")
	enabled, err := c.GetBool("enabled")
	allow_rep, _ := c.GetBool("allow_rep")
	if err != nil {
		enabled = true
	}
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_channel_stream_update_fail", "GM016_030", "id非法", ""))
		return
	}
	if len(rep_url) == 0 {
		c.Json(libs.NewError("admincp_live_channel_stream_update_fail", "GM016_032", "抓取地址必须设置", ""))
		return
	}
	lvs := &lives.LiveOrgs{}
	stream := lvs.GetStream(id)
	if stream == nil {
		c.Json(libs.NewError("admincp_live_channel_stream_update_fail", "GM016_033", "流不存在", ""))
		return
	}
	stream.ReptileUrl = rep_url
	stream.Img = img_id
	stream.Default = is_def
	stream.Enabled = enabled
	stream.AllowRep = allow_rep
	err = lvs.UpdateStream(stream)
	if err != nil {
		c.Json(libs.NewError("admincp_live_channel_stream_update_fail", "GM016_034", "更新失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_channel_stream_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title  删除机构直播频道播放流
// @Description 删除机构直播频道播放流
// @Param   id   path	int true  "流id"
// @Success 200 {object} libs.Error
// @router /org/channel_stream/del [delete]
func (c *LiveCPController) LiveChannelStreamDel() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_channel_stream_del_fail", "GM016_040", "id非法", ""))
		return
	}
	lvs := &lives.LiveOrgs{}
	err := lvs.DeleteStream(id)
	if err != nil {
		c.Json(libs.NewError("admincp_live_channel_stream_del_fail", "GM016_041", "删除失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_channel_stream_del_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
}

// @Title 机构直播频道节目单列表
// @Description 机构直播频道节目单列表
// @Param   date   path	string true  "日期yyyy-MM-dd"
// @Success 200 {object} outobjs.OutProgramObj
// @router /org/program/list [get]
func (c *LiveCPController) LiveProgramList() {
	strd, _ := utils.UrlDecode(c.GetString("date"))
	date, err := utils.StrToTime(strd + " 00:00:00")
	outps := []*outobjs.OutProgramObj{}
	if err != nil {
		c.Json(outps)
		return
	}
	pms := &lives.LivePrograms{}
	programs := pms.Gets(date)
	for _, p := range programs {
		outps = append(outps, outobjs.ConvertOutProgramObj(p, 0))
	}
	c.Json(outps)
}

// @Title 机构直播频道节目单
// @Description 机构直播频道节目单
// @Param   id   path	int true  "id"
// @Success 200 {object} outobjs.OutProgramObj
// @router /org/program/:id([0-9]+) [get]
func (c *LiveCPController) LiveProgramGet() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_program_get_fail", "GM016_050", "id非法", ""))
		return
	}
	pms := &lives.LivePrograms{}
	program := pms.Get(id)
	if program == nil {
		c.Json(libs.NewError("admincp_live_program_get_fail", "GM016_051", "节目单不存在", ""))
		return
	}
	out := *outobjs.ConvertOutProgramObj(program, 0)
	c.Json(out)
}

// @Title 添加机构直播频道主节目单
// @Description 添加机构直播频道主节目单
// @Param   title   path	string true  "标题"
// @Param   sub_title   path	string false  "副标题"
// @Param   start_time   path	string true  "开始时间yyyy-MM-dd hh:mm:ss(填入子节目单后,主节目单开始和结束时间将自动调整)"
// @Param   match_id   path	int false  "关联赛事"
// @Param   game_id   path	int false  "关联游戏"
// @Param   def_channel   path	int true  "默认频道"
// @Param   channel_ids   path	string true  "关联频道"
// @Param   img_id   path	int true  "图片"
// @Success 200 {object} libs.Error
// @router /org/program/add [post]
func (c *LiveCPController) LiveProgramAdd() {
	uid := c.CurrentUid()
	title, _ := utils.UrlDecode(c.GetString("title"))
	sub_title, _ := utils.UrlDecode(c.GetString("sub_title"))
	start_time, _ := utils.UrlDecode(c.GetString("start_time"))
	match_id, _ := c.GetInt64("match_id")
	game_id, _ := c.GetInt64("game_id")
	channel_idsstr := c.GetString("channel_ids")
	def_channel, _ := c.GetInt64("def_channel")
	img_id, _ := c.GetInt64("img_id")

	channel_ids := []int64{}
	cidstr_arr := strings.Split(channel_idsstr, ",")
	for _, cidstr := range cidstr_arr {
		_cid, err := strconv.ParseInt(cidstr, 10, 64)
		if err != nil {
			continue
		}
		channel_ids = append(channel_ids, _cid)
	}

	if len(title) == 0 {
		c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_060", "标题不能为空", ""))
		return
	}
	if img_id <= 0 {
		c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_065", "必须指定图片", ""))
		return
	}
	startTime, err := utils.StrToTime(start_time)
	if err != nil {
		c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_061", "时间格式错误", ""))
		return
	}
	//if game_id <= 0 {
	//	c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_062", "必须选择对应的游戏", ""))
	//	return
	//}
	if def_channel <= 0 {
		c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_063", "未指定默认频道", ""))
		return
	}
	if len(channel_ids) <= 0 {
		c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_063", "未指定对应频道s", ""))
		return
	}
	orgs := &lives.LiveOrgs{}
	for _, cid := range channel_ids {
		if orgs.GetChannel(cid) == nil {
			c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_064", "选择对应的频道不存在", ""))
			return
		}
	}
	//end_time := fmt.Sprintf("%d-%d-%d 23:59:59", startTime.Year(), startTime.Month(), startTime.Day())
	//endTime, _ := utils.StrToTime(end_time)

	program := &lives.LiveProgram{
		Title:            title,
		SubTitle:         sub_title,
		Date:             startTime,
		StartTime:        startTime,
		EndTime:          startTime, //由子菜单自动调整时段
		MatchId:          int(match_id),
		PostUid:          uid,
		GameId:           int(game_id),
		DefaultChannelId: def_channel,
		Img:              img_id,
	}
	pms := &lives.LivePrograms{}
	_, err = pms.Create(program, channel_ids)
	if err != nil {
		c.Json(libs.NewError("admincp_live_program_add_fail", "GM016_065", "添加失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_program_add_succ", controllers.RESPONSE_SUCCESS, "添加成功", ""))
}

// @Title 更改机构直播频道主节目单
// @Description 更改机构直播频道主节目单
// @Param   id   path	int true  "节目单id"
// @Param   title   path	string true  "标题"
// @Param   sub_title   path	string false  "副标题"
// @Param   match_id   path	int false  "关联赛事"
// @Param   game_id   path	int false  "关联游戏"
// @Param   def_channel   path	int true  "默认频道"
// @Param   channel_ids   path	string true  "关联频道"
// @Param   img_id   path	int true  "图片"
// @Success 200 {object} libs.Error
// @router /org/program/update [post]
func (c *LiveCPController) LiveProgramUpdate() {
	uid := c.CurrentUid()
	id, _ := c.GetInt64("id")
	title, _ := utils.UrlDecode(c.GetString("title"))
	sub_title, _ := utils.UrlDecode(c.GetString("sub_title"))
	match_id, _ := c.GetInt64("match_id")
	game_id, _ := c.GetInt64("game_id")
	channel_idsstr := c.GetString("channel_ids")
	def_channel, _ := c.GetInt64("def_channel")
	img_id, _ := c.GetInt64("img_id")

	channel_ids := []int64{}
	cidstr_arr := strings.Split(channel_idsstr, ",")
	for _, cidstr := range cidstr_arr {
		_cid, err := strconv.ParseInt(cidstr, 10, 64)
		if err != nil {
			continue
		}
		channel_ids = append(channel_ids, _cid)
	}

	if len(title) == 0 {
		c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_070", "标题不能为空", ""))
		return
	}
	if img_id <= 0 {
		c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_076", "必须指定图片", ""))
		return
	}
	//if game_id <= 0 {
	//	c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_071", "必须选择对应的游戏", ""))
	//	return
	//}
	if def_channel <= 0 {
		c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_072", "未指定默认频道", ""))
		return
	}
	if len(channel_ids) <= 0 {
		c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_072", "未指定对应频道s", ""))
		return
	}
	orgs := &lives.LiveOrgs{}
	for _, cid := range channel_ids {
		if orgs.GetChannel(cid) == nil {
			c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_073", "选择对应的频道不存在", ""))
			return
		}
	}
	pms := &lives.LivePrograms{}
	program := pms.Get(id)
	if program == nil {
		c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_074", "要修改的节目单不存在", ""))
		return
	}
	program.Title = title
	program.SubTitle = sub_title
	program.MatchId = int(match_id)
	program.GameId = int(game_id)
	program.PostUid = uid
	program.Img = img_id
	program.DefaultChannelId = def_channel
	err := pms.Update(program, channel_ids)
	if err != nil {
		c.Json(libs.NewError("admincp_live_program_update_fail", "GM016_075", "更新失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_program_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 删除机构直播频道主节目单
// @Description 删除机构直播频道主节目单
// @Param   id   path	int true  "节目单id"
// @Success 200 {object} libs.Error
// @router /org/program/del [delete]
func (c *LiveCPController) LiveProgramDelete() {
	id, _ := c.GetInt64("id")
	spms := &lives.LiveSubPrograms{}
	sps := spms.Gets(id)
	if len(sps) > 0 {
		c.Json(libs.NewError("admincp_live_program_del_fail", "GM016_080", "必须先删除子节目单后才能删除主节目单", ""))
		return
	}
	pms := &lives.LivePrograms{}
	err := pms.Delete(id)
	if err != nil {
		c.Json(libs.NewError("admincp_live_program_del_fail", "GM016_081", "删除失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_program_del_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
}

// @Title 添加频道子节目单
// @Description 添加频道子节目单
// @Param   pid   path	int true  "主节目单id"
// @Param   game_id  path	int true  "游戏id"
// @Param   vs1   path	string true  "选手1"
// @Param   vs1_uid   path	int false  "选手1_uid"
// @Param   vs1_img   path	int false  "选手1_img"
// @Param   vs2   path	string true  "选手2"
// @Param   vs2_uid   path	int false  "选手2_uid"
// @Param   vs2_img   path	int false  "选手2_img"
// @Param   t   path	int true  "节目单类型(1对阵,2单个)"
// @Param   title  path	string false  "标题"
// @Param   img  path	int false  "图片"
// @Param   start_time  path	string true  "开始时间(yyyy-MM-dd hh:mm:ss)"
// @Param   end_time  path	string true  "结束时间"
// @Success 200  {object} libs.Error
// @router /org/subprogram/add [post]
func (c *LiveCPController) LiveSubProgramAdd() {
	uid := c.CurrentUid()
	pid, _ := c.GetInt64("pid")
	gameId, _ := c.GetInt64("game_id")
	vs1, _ := utils.UrlDecode(c.GetString("vs1"))
	vs1uid, _ := c.GetInt64("vs1_uid")
	vs1img, _ := c.GetInt64("vs1_img")
	vs2, _ := utils.UrlDecode(c.GetString("vs2"))
	vs2uid, _ := c.GetInt64("vs2_uid")
	vs2img, _ := c.GetInt64("vs2_img")
	t, _ := c.GetInt64("t")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img, _ := c.GetInt64("img")
	start_time, _ := utils.UrlDecode(c.GetString("start_time"))
	end_time, _ := utils.UrlDecode(c.GetString("end_time"))
	if pid <= 0 {
		c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_100", "主频道id大于0", ""))
		return
	}
	if int(t) != int(lives.LIVE_SUBPROGRAM_VIEW_VS) && int(t) != int(lives.LIVE_SUBPROGRAM_VIEW_SINGLE) {
		c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_102", "节目单类型只支持1对阵2单个", ""))
		return
	}
	ptype := lives.LIVE_SUBPROGRAM_VIEW_TYPE(int(t))
	if ptype == lives.LIVE_SUBPROGRAM_VIEW_VS {
		if len(vs1) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_103", "必须指定对阵vs1名称", ""))
			return
		}
		if len(vs2) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_102", "必须指定对阵vs2名称", ""))
			return
		}
		if gameId <= 0 {
			c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_101", "必须指定gameid", ""))
			return
		}
	}
	if ptype == lives.LIVE_SUBPROGRAM_VIEW_SINGLE {
		if len(title) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_103", "必须指定标题", ""))
			return
		}
	}
	if len(start_time) == 0 {
		c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_104", "必须指定开始时间", ""))
		return
	}
	if len(end_time) == 0 {
		c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_105", "必须指定结束时间", ""))
		return
	}

	startTime, s_err := utils.StrToTime(start_time)
	endTime, e_err := utils.StrToTime(end_time)
	if s_err != nil || e_err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_106", "开始或结束时间格式错误", ""))
		return
	}
	lsp := lives.LiveSubProgram{
		ProgramId: pid,
		GameId:    int(gameId),
		Vs1Name:   vs1,
		Vs1Img:    vs1img,
		Vs1Uid:    vs1uid,
		Vs2Name:   vs2,
		Vs2Img:    vs2img,
		Vs2Uid:    vs2uid,
		ViewType:  ptype,
		Title:     title,
		Img:       img,
		StartTime: startTime,
		EndTime:   endTime,
		PostTime:  time.Now(),
		PostUid:   uid,
	}
	lsps := &lives.LiveSubPrograms{}
	_, err := lsps.Create(&lsp)
	if err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_add_fail", "GM016_107", "新增失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_subprogram_add_succ", controllers.RESPONSE_SUCCESS, "新增成功", ""))
}

// @Title 更新频道子节目单
// @Description 更新频道子节目单
// @Param   id   path	int true  "id"
// @Param   game_id  path	int true  "游戏id"
// @Param   vs1   path	string true  "选手1"
// @Param   vs1_uid   path	int false  "选手1_uid"
// @Param   vs1_img   path	int false  "选手1_img"
// @Param   vs2   path	string true  "选手2"
// @Param   vs2_uid   path	int false  "选手2_uid"
// @Param   vs2_img   path	int false  "选手2_img"
// @Param   title  path	string false  "标题"
// @Param   img  path	int false  "图片"
// @Param   start_time  path	string true  "开始时间(yyyy-MM-dd hh:mm:ss)"
// @Param   end_time  path	string true  "结束时间"
// @Success 200  {object} libs.Error
// @router /org/subprogram/update [post]
func (c *LiveCPController) LiveSubProgramUpdate() {
	uid := c.CurrentUid()
	id, _ := c.GetInt64("id")
	gameId, _ := c.GetInt64("game_id")
	vs1, _ := utils.UrlDecode(c.GetString("vs1"))
	vs1uid, _ := c.GetInt64("vs1_uid")
	vs1img, _ := c.GetInt64("vs1_img")
	vs2, _ := utils.UrlDecode(c.GetString("vs2"))
	vs2uid, _ := c.GetInt64("vs2_uid")
	vs2img, _ := c.GetInt64("vs2_img")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img, _ := c.GetInt64("img")
	start_time, _ := utils.UrlDecode(c.GetString("start_time"))
	end_time, _ := utils.UrlDecode(c.GetString("end_time"))
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_110", "id非法", ""))
		return
	}

	lsps := &lives.LiveSubPrograms{}
	lps := lsps.Get(id)
	if lps == nil {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_118", "二级节目单不存在", ""))
		return
	}
	ptype := lps.ViewType
	if ptype == lives.LIVE_SUBPROGRAM_VIEW_VS {
		if len(vs1) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_112", "必须指定对阵vs1名称", ""))
			return
		}
		if len(vs2) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_112", "必须指定对阵vs2名称", ""))
			return
		}
		if gameId <= 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_111", "必须指定gameid", ""))
			return
		}
	}
	if ptype == lives.LIVE_SUBPROGRAM_VIEW_SINGLE {
		if len(title) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_113", "必须指定标题", ""))
			return
		}
	}
	if len(start_time) == 0 {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_114", "必须指定开始时间", ""))
		return
	}
	if len(end_time) == 0 {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_115", "必须指定结束时间", ""))
		return
	}
	startTime, s_err := utils.StrToTime(start_time)
	endTime, e_err := utils.StrToTime(end_time)
	if s_err != nil || e_err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_116", "开始或结束时间格式错误", ""))
		return
	}
	if startTime.After(endTime) {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_118", "开始时间不能大于结束时间", ""))
		return
	}
	if time.Now().After(startTime) {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_119", "开始时间不能小于当前时间", ""))
		return
	}

	lps.GameId = int(gameId)
	lps.Vs1Name = vs1
	lps.Vs1Img = vs1img
	lps.Vs1Uid = vs1uid
	lps.Vs2Name = vs2
	lps.Vs2Img = vs2img
	lps.Vs2Uid = vs2uid
	lps.Title = title
	lps.Img = img
	lps.StartTime = startTime
	lps.EndTime = endTime
	lps.PostTime = time.Now()
	lps.PostUid = uid
	err := lsps.Update(*lps)
	if err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_117", "更新失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_subprogram_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 更新锁定中频道子节目单
// @Description 更新锁定中频道子节目单
// @Param   id   path	int true  "id"
// @Param   game_id  path	int true  "游戏id"
// @Param   vs1   path	string true  "选手1"
// @Param   vs1_uid   path	int false  "选手1_uid"
// @Param   vs1_img   path	int false  "选手1_img"
// @Param   vs2   path	string true  "选手2"
// @Param   vs2_uid   path	int false  "选手2_uid"
// @Param   vs2_img   path	int false  "选手2_img"
// @Param   title  path	string false  "标题"
// @Param   img  path	int false  "图片"
// @Param   end_time  path	string true  "结束时间"
// @Success 200  {object} libs.Error
// @router /org/locking_subprogram/update [post]
func (c *LiveCPController) LiveLockingSubProgramUpdate() {
	uid := c.CurrentUid()
	id, _ := c.GetInt64("id")
	gameId, _ := c.GetInt64("game_id")
	vs1, _ := utils.UrlDecode(c.GetString("vs1"))
	vs1uid, _ := c.GetInt64("vs1_uid")
	vs1img, _ := c.GetInt64("vs1_img")
	vs2, _ := utils.UrlDecode(c.GetString("vs2"))
	vs2uid, _ := c.GetInt64("vs2_uid")
	vs2img, _ := c.GetInt64("vs2_img")
	title, _ := utils.UrlDecode(c.GetString("title"))
	img, _ := c.GetInt64("img")
	end_time, _ := utils.UrlDecode(c.GetString("end_time"))
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_140", "id非法", ""))
		return
	}

	lsps := &lives.LiveSubPrograms{}
	lps := lsps.Get(id)
	if lps == nil {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_148", "二级节目单不存在", ""))
		return
	}
	ptype := lps.ViewType
	if ptype == lives.LIVE_SUBPROGRAM_VIEW_VS {
		if len(vs1) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_142", "必须指定对阵vs1名称", ""))
			return
		}
		if len(vs2) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_142", "必须指定对阵vs2名称", ""))
			return
		}
		if gameId <= 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_141", "必须指定gameid", ""))
			return
		}
	}
	if ptype == lives.LIVE_SUBPROGRAM_VIEW_SINGLE {
		if len(title) == 0 {
			c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_143", "必须指定标题", ""))
			return
		}
	}
	if len(end_time) == 0 {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_145", "必须指定结束时间", ""))
		return
	}
	endTime, e_err := utils.StrToTime(end_time)
	if e_err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_146", "开始或结束时间格式错误", ""))
		return
	}
	if lps.StartTime.After(endTime) {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_148", "结束时间不能小于结束时间", ""))
		return
	}

	lps.GameId = int(gameId)
	lps.Vs1Name = vs1
	lps.Vs1Img = vs1img
	lps.Vs1Uid = vs1uid
	lps.Vs2Name = vs2
	lps.Vs2Img = vs2img
	lps.Vs2Uid = vs2uid
	lps.Title = title
	lps.Img = img
	lps.EndTime = endTime
	lps.PostTime = time.Now()
	lps.PostUid = uid
	err := lsps.UpdateLockingSubProgram(*lps)
	if err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_update_fail", "GM016_147", "更新失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_subprogram_update_succ", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 删除频道子节目单
// @Description 删除频道子节目单
// @Param   id  path	int true  "子节目单id"
// @Success 200  {object} libs.Error
// @router /org/subprogram/del [delete]
func (c *LiveCPController) LiveSubProgramDel() {
	id, _ := c.GetInt64("id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_subprogram_del_fail", "GM016_120", "id非法", ""))
		return
	}
	lsps := &lives.LiveSubPrograms{}
	err := lsps.Delete(id)
	if err != nil {
		c.Json(libs.NewError("admincp_live_subprogram_del_fail", "GM016_121", "删除失败:"+err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_live_subprogram_del_succ", controllers.RESPONSE_SUCCESS, "删除成功", ""))
}

// @Title 频道子节目单列表
// @Description 频道子节目单列表
// @Param   program_id  path	int true  "主节目单id"
// @Success 200  {object} outobjs.OutSubProgramObj
// @router /org/subprogram/list [get]
func (c *LiveCPController) LiveSubProgramList() {
	program_id, _ := c.GetInt64("program_id")
	lsps := &lives.LiveSubPrograms{}
	lsp := &lives.LivePrograms{}
	outs := []outobjs.OutSubProgramObj{}
	program := lsp.Get(program_id)
	if program == nil {
		c.Json(outs)
		return
	}
	list := lsps.Gets(program_id)
	for _, p := range list {
		out_p := outobjs.ConvertOutSubProgramObj(p, program, 0)
		if out_p != nil {
			outs = append(outs, *out_p)
		}
	}
	c.Json(outs)
}

// @Title 频道子节目单
// @Description 频道子节目单
// @Success 200  {object} outobjs.OutSubProgramObj
// @router /org/subprogram/:id([0-9]+) [get]
func (c *LiveCPController) LiveSubProgramGet() {
	id, _ := c.GetInt64(":id")
	if id <= 0 {
		c.Json(libs.NewError("admincp_live_subprogram_get_fail", "GM016_130", "id非法", ""))
		return
	}
	lsp := &lives.LivePrograms{}
	lsps := &lives.LiveSubPrograms{}
	sp := lsps.Get(id)
	if sp == nil {
		c.Json(libs.NewError("admincp_live_subprogram__get_fail", "GM016_131", "节目单不存在", ""))
		return
	}
	p := lsp.Get(sp.ProgramId)
	if p == nil {
		c.Json(libs.NewError("admincp_live_subprogram__get_fail", "GM016_132", "主节目单不存在", ""))
		return
	}
	out := *outobjs.ConvertOutSubProgramObj(sp, p, 0)
	c.Json(out)
}
