package admincp

import (
	"controllers"
	"fmt"
	"libs"
	"libs/groups"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 组管理 API
type GroupCPController struct {
	AdminController
}

func (c *GroupCPController) Prepare() {
	c.AdminController.Prepare()
}

// @Title 获取配置参数
// @Description 获取配置参数
// @Param   id   path  int  false  "配置参数,默认1"
// @Success 200  {object} outobjs.OutConfig
// @router /config/get [get]
func (c *GroupCPController) GetConfig() {
	id, _ := c.GetInt64("id")
	var config *groups.GroupCfg
	if id <= 0 {
		config = groups.GetDefaultCfg()
	} else {
		config = groups.GetGroupCfg(id)
	}
	if config != nil {
		c.Json(&outobjs.OutConfig{
			Id:                           config.Id,
			GroupNameLen:                 config.GroupNameLen,
			GroupDescMaxLen:              config.GroupDescMaxLen,
			GroupDescMinLen:              config.GroupDescMinLen,
			CreateGroupBasePoint:         config.CreateGroupBasePoint,
			CreateGroupRate:              config.CreateGroupRate,
			CreateGroupMinUsers:          config.CreateGroupMinUsers,
			CreateGroupRecruitDay:        config.CreateGroupRecruitDay,
			CreateGroupMaxCount:          config.CreateGroupMaxCount,
			CreateGroupCertifiedMaxCount: config.CreateGroupCertifiedMaxCount,
			CreateGroupClause:            config.CreateGroupClause,
			ReportOptions:                config.ReportOptions,
			NewThreadDefaultStatus:       config.NewThreadDefaultStatus,
		})
	}
	c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_050", "指定配置不存在", ""))
}

// @Title 更新配置参数
// @Description 更新配置参数
// @Param   id   path  int  false  "配置参数,默认1"
// @Param   groupname_len   path  int  true  "组名长度"
// @Param   groupdesc_maxlen  path  int  true  "组描述最大长度"
// @Param   groupdesc_minlen  path  int  true  "组描述最小长度"
// @Param   basepoint  path  int  true  "扣除积分数"
// @Param   minuser   path  int  true  "最小成员数"
// @Param   recruit_day   path  int  true  "招募天数"
// @Param   maxcount   path  int  true  "最大建组数"
// @Param   cretified_maxcount   path  int  true  "认证会员最大建组数"
// @Param   clause   path  string  false  "条款"
// @Param   report_options  path  string  false  "举报选择项(,分隔)"
// @Success 200 {object} libs.Error
// @router /config/update [post]
func (c *GroupCPController) UpdateConfig() {
	id, _ := c.GetInt64("id")
	groupname_len, _ := c.GetInt("groupname_len")
	groupdesc_maxlen, _ := c.GetInt("groupdesc_maxlen")
	groupdesc_minlen, _ := c.GetInt("groupdesc_minlen")
	basepoint, _ := c.GetInt64("basepoint")
	minuser, _ := c.GetInt("minuser")
	recruitday, _ := c.GetInt("recruit_day")
	maxcount, _ := c.GetInt("maxcount")
	cretified_maxcount, _ := c.GetInt("cretified_maxcount")
	clause := c.GetString("clause")
	roptions := c.GetString("report_options")

	if id <= 0 {
		id = 1
	}
	if groupname_len <= 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_060", "参数错误", ""))
		return
	}
	if groupdesc_minlen <= 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_061", "参数错误", ""))
		return
	}
	if groupdesc_maxlen <= 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_062", "参数错误", ""))
		return
	}
	if groupdesc_minlen > groupdesc_maxlen {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_063", "参数错误", ""))
		return
	}
	if basepoint < 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_064", "参数错误", ""))
		return
	}
	if minuser <= 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_065", "参数错误", ""))
		return
	}
	if recruitday <= 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_066", "参数错误", ""))
		return
	}
	if maxcount < 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_067", "参数错误", ""))
		return
	}
	if cretified_maxcount < 0 {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_068", "参数错误", ""))
		return
	}
	cfg := groups.GetGroupCfg(id)
	if cfg == nil {
		c.Json(libs.NewError("admincp_group_cfg_fail", "GM040_067", "配置不存在", ""))
		return
	}
	cfg.GroupNameLen = groupname_len
	cfg.GroupDescMinLen = groupdesc_minlen
	cfg.GroupDescMaxLen = groupdesc_maxlen
	cfg.CreateGroupBasePoint = basepoint
	cfg.CreateGroupMinUsers = minuser
	cfg.CreateGroupRecruitDay = recruitday
	cfg.CreateGroupMaxCount = maxcount
	cfg.CreateGroupCertifiedMaxCount = cretified_maxcount
	cfg.CreateGroupClause = clause
	cfg.ReportOptions = roptions
	err := groups.UpdateGroupCfg(cfg)
	if err != nil {
		c.Json(libs.NewError("admincp_group_cfg_update_fail", "GM040_011", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_group_cfg_update_success", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 获取小组列表
// @Description 获取小组列表
// @Param   page   path  int  false  "页"
// @Param   size   path  int  false  "页数量"
// @Param   words   path  int  false  "搜索关键字"
// @Param   orderby   path  string  false  "排序规则(recommend默认,hot,fans,official)"
// @Success 200  {object} outobjs.OutGroupPagedList
// @router /group/search [get]
func (c *GroupCPController) GroupSearch() {
	words, _ := utils.UrlDecode(c.GetString("words"))
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	orderby := c.GetString("orderby")
	var sortBy groups.GP_SEARCH_SORT
	switch orderby {
	case "hot":
		sortBy = groups.GP_SEARCH_SORT_VITALITY
		break
	case "fans":
		sortBy = groups.GP_SEARCH_SORT_USERS
		break
	case "official":
		sortBy = groups.GP_SEARCH_SORT_OFFICIAL
	default:
		sortBy = groups.GP_SEARCH_SORT_DEFAULT
	}

	match_mode := "any"
	if len(words) == 0 {
		match_mode = "all"
	}
	if size <= 0 {
		size = 20
	}
	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	total, groups := gs.Search(words, page, size, match_mode, sortBy, nil, nil)
	out_groups := outobjs.GetOutGroups(groups, 0)
	out_p := &outobjs.OutGroupPagedList{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  utils.TotalPages(total, size),
		Groups:      out_groups,
		Total:       total,
	}
	c.Json(out_p)
}

// @Title 获取组信息
// @Description 获取组信息
// @Param   group_id  path  int  true  "组id"
// @Success 200 {object} outobjs.OutGroup
// @router /group/get [get]
func (c *GroupCPController) GetGroup() {
	groupid, _ := c.GetInt64("group_id")
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_group_get_fail", "GM040_001", "group_id参数错误", ""))
		return
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := gs.Get(groupid)
	if group == nil {
		c.Json(libs.NewError("admincp_group_get_notexist", "GM040_002", "组不存在", ""))
		return
	}
	c.Json(outobjs.GetOutGroup(group, 0))
}

// @Title 建组
// @Description 建组
// @Param   name   path  string  true  "名称"
// @Param   description   path  string  true  "描述"
// @Param   country  path  string  false  "国家"
// @Param   city  path  string  false  "城市"
// @Param   game_ids  path  string  false  "选择游戏(逗号,分隔)"
// @Param   bgimg   path  int  true  "背景图片id"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Param   invite_uids   path  string  false  "用户uids(用,分隔)"
// @Param   belong   path  string  true  "所属类型"
// @Param   uid  path  int  true  "所属uid"
// @Param   recommend   path  bool  false  "推荐"
// @Param   displayorder   path  int  false  "排序"
// @Success 200 {object} libs.Error
// @router /group/create [post]
func (c *GroupCPController) CreateGroup() {
	name, _ := utils.UrlDecode(c.GetString("name"))
	desc, _ := utils.UrlDecode(c.GetString("description"))
	country, _ := utils.UrlDecode(c.GetString("country"))
	city, _ := utils.UrlDecode(c.GetString("city"))
	gameids := c.GetString("game_ids")
	bgimg, _ := c.GetInt64("bgimg")
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	inviteuids := c.GetString("invite_uids")
	uid, _ := c.GetInt64("uid")
	bl, _ := c.GetInt("belong")
	remd, _ := c.GetBool("recommend")
	displayorder, _ := c.GetInt("displayorder")

	var belong groups.GROUP_BELONG
	switch bl {
	case int(groups.GROUP_BELONG_MEMBER):
		belong = groups.GROUP_BELONG_MEMBER
		break
	case int(groups.GROUP_BELONG_OFFICIAL):
		belong = groups.GROUP_BELONG_OFFICIAL
		break
	default:
		c.Json(libs.NewError("admincp_group_create_belong_fail", "GM040_010", "所属类型错误", ""))
		return
	}

	inv_uids := []int64{}
	inv_uidss := strings.Split(inviteuids, ",")
	for _, str := range inv_uidss {
		_id, _ := strconv.ParseInt(str, 10, 64)
		if _id > 0 {
			inv_uids = append(inv_uids, _id)
		}
	}

	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := &groups.Group{
		Name:         name,
		Description:  desc,
		Uid:          uid,
		Country:      country,
		City:         city,
		GameIds:      gameids,
		BgImg:        bgimg,
		Belong:       belong,
		Type:         groups.GROUP_TYPE_NORMAL,
		LongiTude:    float32(longitude),
		LatiTude:     float32(latitude),
		InviteUids:   inv_uids,
		Recommend:    remd,
		DisplarOrder: displayorder,
	}
	err := gs.Create(group)
	if err != nil {
		c.Json(libs.NewError("admincp_group_create_fail", "GM040_011", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_group_create_success", controllers.RESPONSE_SUCCESS, fmt.Sprintf("%d", group.Id), ""))
}

// @Title 关闭组
// @Description 关闭组
// @Param   group_id  path  int  true  "组id"
// @Success 200 {object} libs.Error
// @router /group/close [post]
func (c *GroupCPController) CloseGroup() {
	groupid, _ := c.GetInt64("group_id")
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_group_close_fail", "GM040_020", "id错误", ""))
		return
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	err := gs.Close(groupid)
	if err != nil {
		c.Json(libs.NewError("admincp_group_close_fail", "GM040_021", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_group_close_success", controllers.RESPONSE_SUCCESS, "关闭成功", ""))
}

// @Title 更新小组信息
// @Description 更新小组信息
// @Param   group_id  path  int  true  "组id"
// @Param   description   path  string  true  "描述"
// @Param   country  path  string  false  "国家"
// @Param   city  path  string  false  "城市"
// @Param   game_ids  path  string  false  "选择游戏(逗号,分隔)"
// @Param   bgimg   path  int  true  "背景图片id"
// @Param   recommend   path  bool  false  "推荐"
// @Param   displayorder   path  int  false  "排序"
// @Success 200 {object} libs.Error
// @router /group/update [post]
func (c *GroupCPController) UpdateGroup() {
	groupid, _ := c.GetInt64("group_id")
	desc, _ := utils.UrlDecode(c.GetString("description"))
	country, _ := utils.UrlDecode(c.GetString("country"))
	city, _ := utils.UrlDecode(c.GetString("city"))
	gameids := c.GetString("game_ids")
	bgimg, _ := c.GetInt64("bgimg")
	remd, rmderr := c.GetBool("recommend")
	displayorder, _ := c.GetInt("displayorder")
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := gs.Get(groupid)
	if group == nil {
		c.Json(libs.NewError("admincp_group_update_fail", "GM040_031", "小组不存在", ""))
		return
	}
	group.Description = desc
	group.Country = country
	group.City = city
	group.GameIds = gameids
	group.BgImg = bgimg
	if rmderr == nil {
		group.Recommend = remd
	}
	group.DisplarOrder = displayorder
	err := gs.Update(group)
	if err != nil {
		c.Json(libs.NewError("admincp_group_update_fail", "GM040_032", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_group_update_success", controllers.RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 关闭帖子
// @Description 关闭帖子
// @Param   thread_id  path  int  true  "帖子id"
// @Success 200 {object} libs.Error
// @router /thread/close [post]
func (c *GroupCPController) CloseThread() {
	threadid, _ := c.GetInt64("thread_id")
	if threadid <= 0 {
		c.Json(libs.NewError("admincp_thread_close_fail", "GM040_030", "id错误", ""))
		return
	}
	ts := groups.NewThreadService(groups.GetDefaultCfg())
	err := ts.CloseThread(threadid)
	if err != nil {
		c.Json(libs.NewError("admincp_thread_close_fail", "GM040_031", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_thread_close_success", controllers.RESPONSE_SUCCESS, "关闭成功", ""))
}

// @Title 设置帖子排序
// @Description 设置帖子排序
// @Param   thread_id  path  int  true  "帖子id"
// @Param   displayorder  path  int  true  "排序数(越大排越前)"
// @Success 200 {object} libs.Error
// @router /thread/set_displayorder [post]
func (c *GroupCPController) SetThreadOrder() {
	threadid, _ := c.GetInt64("thread_id")
	if threadid <= 0 {
		c.Json(libs.NewError("admincp_thread_setorder_fail", "GM040_110", "id错误", ""))
		return
	}
	displayorder, _ := c.GetInt("displayorder")
	ts := groups.NewThreadService(groups.GetDefaultCfg())
	err := ts.SetDisplayOrder(threadid, displayorder)
	if err != nil {
		c.Json(libs.NewError("admincp_thread_setorder_fail", "GM040_111", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_thread_setorder_success", controllers.RESPONSE_SUCCESS, "设置成功", ""))
}

// @Title 获取帖子信息
// @Description 获取帖子信息
// @Param   thread_id  path  int  true  "帖子id"
// @Success 200 {object} outobjs.OutThread
// @router /thread/get [get]
func (c *GroupCPController) GetThread() {
	thread_id, _ := c.GetInt64("thread_id")
	if thread_id <= 0 {
		c.Json(libs.NewError("admincp_thread_get_fail", "GM040_080", "id错误", ""))
		return
	}
	ts := groups.NewThreadService(groups.GetDefaultCfg())
	thread := ts.Get(thread_id)
	if thread == nil {
		c.Json(libs.NewError("admincp_thread_get_fail", "GM040_081", "帖子不存在", ""))
		return
	}
	c.Json(outobjs.GetOutThread(thread))
}

// @Title 获取帖子列表
// @Description 获取帖子列表
// @Param   group_id  path  int  true  "组id"
// @Param   page  path  int  false  "页"
// @Param   size  path  int  false  "页数量"
// @Success 200 {object} outobjs.OutThreadPagedListForAdmin
// @router /thread/list [get]
func (c *GroupCPController) GetThreads() {
	group_id, _ := c.GetInt64("group_id")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if group_id <= 0 {
		c.Json(libs.NewError("admincp_thread_list_fail", "GM040_090", "group_id错误", ""))
		return
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	ts := groups.NewThreadService(groups.GetDefaultCfg())
	total, threads := ts.Gets(group_id, page, size, 0)
	out_ts := []*outobjs.OutThread{}
	for _, thread := range threads {
		out_ts = append(out_ts, outobjs.GetOutThread(thread))
	}
	out_p := &outobjs.OutThreadPagedListForAdmin{
		CurrentPage: page,
		PageSize:    size,
		Threads:     out_ts,
		Total:       total,
		Pages:       utils.TotalPages(total, size),
	}
	c.Json(out_p)
}

// @Title 新建帖子
// @Description 新建帖子
// @Param   uid  path  int  true  "uid"
// @Param   group_id   path  int  true  "组id"
// @Param   subject   path  string  true  "标题"
// @Param   message   path  string  true  "内容"
// @Param   img_ids   path  string  true  "图片集(最大9张 逗号,分隔)"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Param   fromdev   path  string  false  "设备标识(android,ios,ipad,wphone,web)"
// @Success 200 {object} libs.Error
// @router /thread/submit [post]
func (c *GroupCPController) CreateThread() {
	uid, _ := c.GetInt64("uid")
	groupid, _ := c.GetInt64("group_id")
	subject, _ := utils.UrlDecode(c.GetString("subject"))
	message, _ := utils.UrlDecode(c.GetString("message"))
	imgids := c.GetString("img_ids")
	fromdev := c.GetString("fromdev")
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	if uid <= 0 {
		c.Json(libs.NewError("admincp_createthread_fail", "GM040_120", "uid错误", ""))
		return
	}

	ths := groups.NewThreadService(groups.GetDefaultCfg())
	thread := &groups.Thread{
		GroupId:  groupid,
		Subject:  subject,
		AuthorId: uid,
	}
	img_ids := []int64{}
	arrImg := strings.Split(imgids, ",")
	for _, _ai := range arrImg {
		_id, _ := strconv.ParseInt(_ai, 10, 64)
		if _id > 0 {
			img_ids = append(img_ids, _id)
		}
	}
	if len(img_ids) > 9 {
		c.Json(libs.NewError("group_createthread_fail", "GM040_121", "图片数量不能大于9张", ""))
		return
	}
	post := &groups.Post{
		AuthorId:   uid,
		Subject:    subject,
		Message:    message,
		Ip:         c.Ctx.Input.IP(),
		FromDevice: groups.GetFromDevice(fromdev),
		ImgIds:     img_ids,
		LongiTude:  float32(longitude),
		LatiTude:   float32(latitude),
	}
	err := ths.Create(thread, post)
	if err != nil {
		c.Json(libs.NewError("group_createthread_fail", "GM040_122", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_createthread_success", controllers.RESPONSE_SUCCESS, "新建帖子成功", ""))
}

// @Title 隐藏评论
// @Description 隐藏评论
// @Param   post_id  path  string  true  "评论id"
// @Success 200 {object} libs.Error
// @router /post/invisible [post]
func (c *GroupCPController) ClosePost() {
	postid := c.GetString("post_id")
	if len(postid) == 0 {
		c.Json(libs.NewError("admincp_post_invisible_fail", "GM040_040", "id错误", ""))
		return
	}
	ps := groups.NewPostService(groups.GetDefaultCfg())
	err := ps.Invisible(postid, false)
	if err != nil {
		c.Json(libs.NewError("admincp_post_invisible_fail", "GM040_041", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("admincp_post_invisible_success", controllers.RESPONSE_SUCCESS, "关闭成功", ""))
}

// @Title 获取评论列表
// @Description 获取评论列表
// @Param   thread_id  path  int  true  "帖子id"
// @Param   page  path  int  false  "页"
// @Param   size  path  int  false  "页数量"
// @Success 200 {object} outobjs.OutPostPagedList
// @router /post/list [get]
func (c *GroupCPController) GetPosts() {
	threadid, _ := c.GetInt64("thread_id")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	oby := groups.POST_ORDERBY_ASC
	cfg := groups.GetDefaultCfg()
	ps := groups.NewPostService(cfg)
	ths := groups.NewThreadService(cfg)
	thread := ths.Get(threadid)
	if thread == nil {
		c.Json(libs.NewError("group_getpost_fail", "GM040_050", "帖子不存在", ""))
		return
	}
	//用户顶踩记录
	actionTags := ps.MemberThreadPostActionTags(threadid, 0)
	dcTag := func(outp *outobjs.OutPost) {
		if dc, ok := actionTags[outp.Id]; ok {
			if dc == groups.POST_ACTIONTAG_DING {
				outp.Dinged = true
			}
			if dc == groups.POST_ACTIONTAG_CAI {
				outp.Caied = true
			}
		}
	}
	//顶计数查询ids,进行统一获取
	postids := []string{}
	postids = append(postids, thread.Lordpid) //楼主评论

	//最大顶post
	maxdings := ps.GetTops(threadid, 2, groups.POST_ACTIONTAG_DING)
	var maxding *outobjs.OutPost
	for _, md := range maxdings {
		if md.Position != 1 {
			maxding = outobjs.GetOutPost(md, ps.GetSrcToRes(md.Resources))
			dcTag(maxding)
			postids = append(postids, maxding.Id)
			break
		}
	}
	//楼主post
	lordPost := ps.Get(thread.Lordpid)
	lordOutPost := outobjs.GetOutPost(lordPost, ps.GetSrcToRes(lordPost.Resources))
	dcTag(lordOutPost)

	//获取post列表
	total, posts := ps.Gets(threadid, page, size, oby, false)

	//顶记录
	for _, p := range posts {
		postids = append(postids, p.Id)
	}
	dingCs := ps.GetPostActionCounts(threadid, postids, groups.POST_ACTIONTAG_DING)
	dingCsF := func(outp *outobjs.OutPost) {
		if outp == nil {
			return
		}
		if cs, ok := dingCs[outp.Id]; ok {
			outp.Ding = cs
		}
	}

	//转换输出模型
	out_posts := []*outobjs.OutPost{}
	for _, p := range posts {
		res := ps.GetSrcToRes(p.Resources)
		_outp := outobjs.GetOutPost(p, res)
		dcTag(_outp)
		dingCsF(_outp)
		out_posts = append(out_posts, _outp)
	}
	dingCsF(maxding)
	dingCsF(lordOutPost)

	out_p := &outobjs.OutPostPagedList{
		CurrentPage: page,
		TotalPages:  utils.TotalPages(total, size),
		Total:       total,
		PageSize:    size,
		Posts:       out_posts,
		MaxDingPost: maxding,
		MaxCaiPost:  nil,
		Thread:      outobjs.GetOutThread(thread),
		JoinedGroup: false,
		LordPost:    lordOutPost,
	}
	c.Json(out_p)
}

// @Title 获取举报列表
// @Description 获取举报列表
// @Param   category   path  int  false  "分类"
// @Param   page   path  int  false  "页"
// @Param   size   path  int  false  "页数量"
// @Success 200  {object} outobjs.OutReportPagedList
// @router /report/list [get]
func (c *GroupCPController) Reports() {
	category, _ := c.GetInt("category")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var cg groups.REPORT_CATEGORY
	switch category {
	case int(groups.REPORT_CATEGORY_POST):
		cg = groups.REPORT_CATEGORY_POST
	case int(groups.REPORT_CATEGORY_THREAD):
		cg = groups.REPORT_CATEGORY_THREAD
	case int(groups.REPORT_CATEGORY_GROUP):
		cg = groups.REPORT_CATEGORY_GROUP
	}
	rs := groups.NewReportService()
	total, reports := rs.Gets(page, size, cg)
	outs := []*outobjs.OutReport{}
	for _, r := range reports {
		outs = append(outs, &outobjs.OutReport{
			Id:      r.Id,
			RefId:   r.RefId,
			C:       r.C,
			Ts:      time.Unix(r.Ts, 0),
			RefTxt:  r.RefTxt,
			PostUid: r.PostUid,
			Member:  outobjs.GetOutSimpleMember(r.PostUid),
			Msg:     r.Msg,
		})
	}
	outp := &outobjs.OutReportPagedList{
		CurrentPage: page,
		TotalPages:  utils.TotalPages(total, size),
		PageSize:    size,
		Reports:     outs,
	}
	c.Json(outp)
}
