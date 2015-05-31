package controllers

import (
	"fmt"
	"libs"
	"libs/groups"
	"libs/search"
	"libs/share"
	"outobjs"
	"strconv"
	"strings"
	"time"
	"utils"
)

// 组 API
type GroupController struct {
	BaseController
}

func (c *GroupController) Prepare() {
}

func (c *GroupController) URLMapping() {
	c.Mapping("GetSetting", c.GetSetting)
	c.Mapping("GetGroup", c.GetGroup)
	c.Mapping("GetGroups", c.GetGroups)
	c.Mapping("GetRecruitingGroups", c.GetRecruitingGroups)
	c.Mapping("GetMyGroups", c.GetMyGroups)
	c.Mapping("GetMyJoinGroups", c.GetMyJoinGroups)
	c.Mapping("CreateGroup", c.CreateGroup)
	c.Mapping("UpdateGroup", c.UpdateGroup)
	c.Mapping("JoinGroup", c.JoinGroup)
	c.Mapping("ExitGroup", c.ExitGroup)
	c.Mapping("InvitedFriendList", c.InvitedFriendList)
	c.Mapping("Invite", c.Invite)
	c.Mapping("GetThreads", c.GetThreads)
	c.Mapping("GetThreads", c.GetThreads)
	c.Mapping("GetPost", c.GetPost)
	c.Mapping("GetPosts", c.GetPosts)
	c.Mapping("CreatePost", c.CreatePost)
	c.Mapping("ActionPost", c.ActionPost)
	c.Mapping("ReportOptions", c.ReportOptions)
	c.Mapping("Report", c.Report)
	c.Mapping("Share", c.Share)
	c.Mapping("MsgCount", c.MsgCount)
}

// @Title 申请组的设定值
// @Description 申请组的设定值
// @Success 200 {object} outobjs.OutGroupSetting
// @router /group/setting [get]
func (c *GroupController) GetSetting() {
	setting := groups.GetDefaultCfg()
	out_setting := &outobjs.OutGroupSetting{
		GroupNameLen:    setting.GroupNameLen,
		GroupDescMaxLen: setting.GroupDescMaxLen,
		GroupDescMinLen: setting.GroupDescMinLen,
		DeductPoint:     setting.CreateGroupBasePoint,
		MinUsers:        setting.CreateGroupMinUsers,
		LimitDay:        setting.CreateGroupRecruitDay,
		GroupClause:     setting.CreateGroupClause,
	}
	c.Json(out_setting)
}

// @Title 获取组信息
// @Description 获取组信息
// @Param   access_token  path  string  false  "access_token"
// @Param   group_id  path  int  true  "组id"
// @Success 200 {object} outobjs.OutGroup
// @router /group/get [get]
func (c *GroupController) GetGroup() {
	current_uid := c.CurrentUid()
	groupid, _ := c.GetInt64("group_id")
	if groupid <= 0 {
		c.Json(libs.NewError("group_get_fail", "GP1001", "group_id参数错误", ""))
		return
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := gs.Get(groupid)
	if group == nil {
		c.Json(libs.NewError("group_get_notexist", "GP1002", "组不存在", ""))
		return
	}
	c.Json(outobjs.GetOutGroup(group, current_uid))
}

// @Title 获取组列表
// @Description 获取组列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   page   path  int  false  "页"
// @Param   words   path  int  false  "搜索关键字"
// @Param   game_ids   path  string  false  "游戏ids(逗号,分隔)"
// @Param   orderby   path  string  false  "排序规则(recommend默认,hot,fans,official)"
// @Success 200 {object} outobjs.OutGroupPagedList
// @router /group/list [get]
func (c *GroupController) GetGroups() {
	current_uid := c.CurrentUid()
	words, _ := utils.UrlDecode(c.GetString("words"))
	gameids := c.GetString("game_ids")
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
	fgameids := []uint64{}
	_strids := strings.Split(gameids, ",")
	for _, _id := range _strids {
		gid, _ := strconv.ParseUint(_id, 10, 64)
		fgameids = append(fgameids, gid)
	}
	var filters []search.SearchFilter
	if len(fgameids) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "gameids",
			Values:  fgameids,
			Exclude: false,
		})
	}
	filters = append(filters, search.SearchFilter{
		Attr:    "status",
		Values:  []uint64{uint64(groups.GROUP_STATUS_OPENING), uint64(groups.GROUP_STATUS_LOWMEMBER)},
		Exclude: false,
	})

	match_mode := "any"
	if len(words) == 0 {
		match_mode = "all"
	}
	if size <= 0 {
		size = 20
	}
	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	total, groups := gs.Search(words, page, size, match_mode, sortBy, filters, nil)
	out_groups := outobjs.GetOutGroups(groups, current_uid)
	out_p := &outobjs.OutGroupPagedList{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  utils.TotalPages(total, size),
		Groups:      out_groups,
		Total:       total,
	}
	c.Json(out_p)
}

// @Title 获取招募中的组列表
// @Description 获取招募中的组列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   game_ids   path  string  false  "游戏ids(逗号,分隔)"
// @Param   page   path  int  false  "页"
// @Param   size   path  int  false  "页数"
// @Success 200 {object} outobjs.OutGroupPagedList
// @router /group/recruiting [get]
func (c *GroupController) GetRecruitingGroups() {
	current_uid := c.CurrentUid()
	gameids := c.GetString("game_ids")
	page, _ := c.GetInt("page")
	size, _ := c.GetInt("size")
	fgameids := []uint64{}
	_strids := strings.Split(gameids, ",")
	for _, _id := range _strids {
		gid, _ := strconv.ParseUint(_id, 10, 64)
		fgameids = append(fgameids, gid)
	}
	var filters []search.SearchFilter
	if len(fgameids) > 0 {
		filters = append(filters, search.SearchFilter{
			Attr:    "gameids",
			Values:  fgameids,
			Exclude: false,
		})
	}
	filters = append(filters, search.SearchFilter{
		Attr:    "status",
		Values:  []uint64{uint64(groups.GROUP_STATUS_RECRUITING)},
		Exclude: false,
	})
	if size <= 0 {
		size = 20
	}

	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	total, groups := gs.Search("", page, size, "all", groups.GP_SEARCH_SORT_ENDTIME, filters, nil)
	out_groups := outobjs.GetOutGroups(groups, current_uid)

	out_p := &outobjs.OutGroupPagedList{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  utils.TotalPages(total, size),
		Groups:      out_groups,
	}
	c.Json(out_p)
}

// @Title 用户创建的组列表
// @Description 用户创建的组列表(返回数组)
// @Param   access_token  path  string  true  "access_token"
// @Success 200 {object} outobjs.OutMyGroups
// @router /group/my [get]
func (c *GroupController) GetMyGroups() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_my_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询", ""))
		return
	}
	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	mygroups := gs.MyGroups(current_uid)
	out_mygroups := &outobjs.OutMyGroups{
		MaxAllowGroupCount: gs.AllowMaxCreateLimits(current_uid),
	}
	out_groups := []*outobjs.OutGroup{}
	for _, gp := range mygroups {
		out_g := outobjs.GetOutGroup(gp, 0)
		out_g.IsJoined = true //自己的组肯定是加入状态
		out_groups = append(out_groups, out_g)
	}
	out_mygroups.Groups = out_groups
	c.Json(out_mygroups)
}

// @Title 获取我加入的组列表
// @Description 获取我加入的组列表(返回数组)
// @Param   access_token  path  string  true  "access_token"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutGroupPagedList
// @router /group/myjoins [get]
func (c *GroupController) GetMyJoinGroups() {
	current_uid := c.CurrentUid()
	page, _ := c.GetInt("page", 1)
	if current_uid <= 0 {
		c.Json(libs.NewError("group_myjoins_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询", ""))
		return
	}
	if page <= 0 {
		page = 1
	}
	size := 20
	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	total, joinGroups := gs.MyJoins(current_uid, page, size)
	out_groups := []*outobjs.OutGroup{}
	for _, gp := range joinGroups {
		out_g := outobjs.GetOutGroup(gp, 0)
		out_g.IsJoined = true //肯定是true
		out_groups = append(out_groups, out_g)
	}
	out_p := &outobjs.OutGroupPagedList{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  utils.TotalPages(total, size),
		Groups:      out_groups,
	}
	c.Json(out_p)
}

// @Title 申请建组
// @Description 申请建组
// @Param   access_token  path  string  true  "access_token"
// @Param   name   path  string  true  "名称"
// @Param   description   path  string  true  "描述"
// @Param   country  path  string  false  "国家"
// @Param   city  path  string  false  "城市"
// @Param   game_ids  path  string  false  "选择游戏(逗号,分隔)"
// @Param   bgimg   path  int  true  "背景图片id"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Param   invite_uids   path  string  false  "用户uids(用,分隔)"
// @Success 200 {object} libs.Error
// @router /group/create [post]
func (c *GroupController) CreateGroup() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_create_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能创建", ""))
		return
	}
	name, _ := utils.UrlDecode(c.GetString("name"))
	desc, _ := utils.UrlDecode(c.GetString("description"))
	country, _ := utils.UrlDecode(c.GetString("country"))
	city, _ := utils.UrlDecode(c.GetString("city"))
	gameids := c.GetString("game_ids")
	bgimg, _ := c.GetInt64("bgimg")
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	inviteuids := c.GetString("invite_uids")
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
		Name:        name,
		Description: desc,
		Uid:         current_uid,
		Country:     country,
		City:        city,
		GameIds:     gameids,
		BgImg:       bgimg,
		Belong:      groups.GROUP_BELONG_MEMBER,
		Type:        groups.GROUP_TYPE_NORMAL,
		LongiTude:   float32(longitude),
		LatiTude:    float32(latitude),
		InviteUids:  inv_uids,
	}
	err := gs.Create(group)
	if err != nil {
		c.Json(libs.NewError("group_create_fail", "GP1050", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_create_success", RESPONSE_SUCCESS, fmt.Sprintf("%d", group.Id), ""))
}

// @Title 建组前校验用户是否满足条件
// @Description 建组前校验用户是否满足条件
// @Param   access_token  path  string  true  "access_token"
// @Success 200 {object} libs.Error
// @router /group/create_check [post]
func (c *GroupController) CreateGroupCheck() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_create_check_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能检查", ""))
		return
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	err := gs.CheckMemberNewGroupPass(current_uid, groups.GROUP_BELONG_MEMBER)
	if err != nil {
		c.Json(libs.NewError("group_check_fail", "GP1100", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_check_success", RESPONSE_SUCCESS, "检查成功", ""))
}

// @Title 更新组属性
// @Description 更新组属性
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Param   description   path  string  true  "描述"
// @Param   bgimg   path  int  true  "背景图片id"
// @Success 200 {object} libs.Error
// @router /group/update [post]
func (c *GroupController) UpdateGroup() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_update_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能更新", ""))
		return
	}
	groupid, _ := c.GetInt64("groupid")
	desc, _ := utils.UrlDecode(c.GetString("description"))
	bgimg, _ := c.GetInt64("bgimg")
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	group := gs.Get(groupid)
	if group == nil || group.Status == groups.GROUP_STATUS_CLOSED {
		c.Json(libs.NewError("group_update_fail", "GP1060", "组不存在", ""))
		return
	}
	if group.Uid != current_uid {
		c.Json(libs.NewError("group_update_fail", "GP1060", "非本人不能更新组属性", ""))
		return
	}
	group.Description = desc
	if bgimg > 0 {
		group.BgImg = bgimg
	}
	err := gs.Update(group)
	if err != nil {
		c.Json(libs.NewError("group_update_fail", "GP1060", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_update_success", RESPONSE_SUCCESS, "更新成功", ""))
}

// @Title 加入组
// @Description 加入组
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Success 200 {object} libs.Error
// @router /group/join [post]
func (c *GroupController) JoinGroup() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_join_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能加入", ""))
		return
	}
	groupid, _ := c.GetInt64("groupid")
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	err := gs.Join(current_uid, groupid)
	if err != nil {
		c.Json(libs.NewError("group_join_fail", "GP1070", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_join_success", RESPONSE_SUCCESS, "成功加入小组", ""))
}

// @Title 离开组
// @Description 离开组
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Success 200 {object} libs.Error
// @router /group/exit [post]
func (c *GroupController) ExitGroup() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_exit_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能加入", ""))
		return
	}
	groupid, _ := c.GetInt64("groupid")
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	err := gs.Exit(current_uid, groupid)
	if err != nil {
		c.Json(libs.NewError("group_exit_fail", "GP1080", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_exit_success", RESPONSE_SUCCESS, "成功离开小组", ""))
}

// @Title 邀请好友列表
// @Description 邀请好友列表
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Success 200
// @router /group/invite_friends [get]
func (c *GroupController) InvitedFriendList() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_invitefriend_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能查询", ""))
		return
	}
	groupid, _ := c.GetInt64("groupid")
	if groupid <= 0 {
		c.Json(libs.NewError("group_invitefriend_fail", "GP1090", "组id错误", ""))
		return
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	maps := gs.InviteFriends(current_uid, groupid)
	type OutPy struct {
		Key     string                     `json:"w"`
		Members []*outobjs.OutInviteMember `json:"ims"`
	}
	outs := []*OutPy{}
	for _, v := range py_chars {
		vs, ok := maps[string(v)]
		if ok {
			out_ims := []*outobjs.OutInviteMember{}
			for _, im := range vs {
				out_member := outobjs.GetOutSimpleMember(im.Uid)
				if out_member != nil {
					out_ims = append(out_ims, &outobjs.OutInviteMember{
						Uid:     im.Uid,
						Member:  out_member,
						Invited: im.Invited,
						Joined:  im.Joined,
					})
				}
			}
			outs = append(outs, &OutPy{
				Key:     string(v),
				Members: out_ims,
			})
		}
	}
	c.Json(outs)
}

// @Title 邀请好友
// @Description 邀请好友
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Param   uids   path  string  true  "用户uids(用,分隔)"
// @Success 200  {object} libs.Error
// @router /group/invite [post]
func (c *GroupController) Invite() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_invite_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能邀请", ""))
		return
	}
	groupid, _ := c.GetInt64("groupid")
	if groupid <= 0 {
		c.Json(libs.NewError("group_invite_fail", "GP1190", "组id错误", ""))
		return
	}
	uidstrs := c.GetString("uids")
	uids := []int64{}
	uidss := strings.Split(uidstrs, ",")
	for _, str := range uidss {
		_id, _ := strconv.ParseInt(str, 10, 64)
		if _id > 0 {
			uids = append(uids, _id)
		}
	}
	gs := groups.NewGroupService(groups.GetDefaultCfg())
	gs.Invite(current_uid, groupid, uids)
	c.Json(libs.NewError("group_invite_success", RESPONSE_SUCCESS, "成功邀请用户", ""))
}

// @Title 获取组帖子
// @Description 获取组帖子列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   group_id   path  int  true  "组id"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutThreadPagedList
// @router /thread/list [get]
func (c *GroupController) GetThreads() {
	current_uid := c.CurrentUid()
	groupid, _ := c.GetInt64("group_id")
	page, _ := c.GetInt("page")

	if groupid <= 0 {
		c.Json(libs.NewError("group_getthreads_fail", "GP1200", "组id错误", ""))
	}
	if page <= 0 {
		page = 1
	}
	size := 20

	//5秒缓冲
	query_cache_key := fmt.Sprintf("front_fast_cache.group.threads:g_%d_p_%d_s_%d", groupid, page, size)
	c_obj := utils.GetLocalFastExpriesTimePartCache(query_cache_key)
	if c_obj != nil {
		c.Json(c_obj)
		return
	}

	ths := groups.NewThreadService(groups.GetDefaultCfg())
	_, threads := ths.Gets(groupid, page, size, current_uid)
	out_threads := []*outobjs.OutThread{}
	for _, t := range threads {
		out_threads = append(out_threads, outobjs.GetOutThread(t))
	}
	out_p := &outobjs.OutThreadPagedList{
		CurrentPage: page,
		PageSize:    size,
		Threads:     out_threads,
	}
	utils.SetLocalFastExpriesTimePartCache(5*time.Second, query_cache_key, out_p)
	c.Json(out_p)
}

// @Title 新建帖子
// @Description 新建帖子
// @Param   access_token  path  string  false  "access_token"
// @Param   group_id   path  int  true  "组id"
// @Param   subject   path  string  true  "标题"
// @Param   message   path  string  true  "内容"
// @Param   img_ids   path  string  true  "图片集(最大9张 逗号,分隔)"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Param   fromdev   path  string  false  "设备标识(android,ios,ipad,wphone,web)"
// @Success 200 {object} libs.Error
// @router /thread/submit [post]
func (c *GroupController) CreateThread() {
	current_uid := c.CurrentUid()
	groupid, _ := c.GetInt64("group_id")
	subject, _ := utils.UrlDecode(c.GetString("subject"))
	message, _ := utils.UrlDecode(c.GetString("message"))
	imgids := c.GetString("img_ids")
	fromdev := c.GetString("fromdev")
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	ths := groups.NewThreadService(groups.GetDefaultCfg())
	thread := &groups.Thread{
		GroupId:  groupid,
		Subject:  subject,
		AuthorId: current_uid,
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
		c.Json(libs.NewError("group_newthread_fail", "GP1300", "图片数量不能大于9张", ""))
		return
	}
	post := &groups.Post{
		AuthorId:   current_uid,
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
		c.Json(libs.NewError("group_newthread_fail", "GP1301", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_newthread_success", RESPONSE_SUCCESS, "新建帖子成功", ""))
}

// @Title 获取单个评论
// @Description 获取单个评论
// @Param   access_token  path  string  false  "access_token"
// @Param   post_id   path  string  true  "评论id"
// @Success 200 {object} outobjs.OutSinglePost
// @router /post/get [get]
func (c *GroupController) GetPost() {
	current_uid := c.CurrentUid()
	postid := c.GetString("post_id")
	cfg := groups.GetDefaultCfg()
	ps := groups.NewPostService(cfg)
	post := ps.Get(postid)
	if post == nil {
		c.Json(libs.NewError("group_getpost_fail", "GP1400", "评论不存在", ""))
		return
	}
	res := ps.GetSrcToRes(post.Resources)
	outp := outobjs.GetOutPost(post, res)
	outsp := outobjs.OutSinglePost{
		JoinedGroup: false,
		Post:        outp,
	}

	//顶记录
	dingCs := ps.GetPostActionCounts(post.ThreadId, []string{post.Id}, groups.POST_ACTIONTAG_DING)
	if cs, ok := dingCs[outp.Id]; ok {
		outp.Ding = cs
	}

	if current_uid > 0 {
		//用户顶踩记录
		actionTags := ps.MemberThreadPostActionTags(post.ThreadId, current_uid)
		if dc, ok := actionTags[outp.Id]; ok {
			if dc == groups.POST_ACTIONTAG_DING {
				outp.Dinged = true
			}
			if dc == groups.POST_ACTIONTAG_CAI {
				outp.Caied = true
			}
		}
		ts := groups.NewThreadService(cfg)
		thread := ts.Get(post.ThreadId)
		if thread != nil {
			gs := groups.NewGroupService(cfg)
			outsp.Thread = outobjs.GetOutThread(thread)
			outsp.JoinedGroup = gs.IsJoined(current_uid, thread.GroupId)
		}
	}
	c.Json(outsp)
}

// @Title 帖子评论列表
// @Description 帖子评论列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   thread_id   path  int  true  "帖子id"
// @Param   page   path  int  false  "页"
// @Param   orderby   path  string  false  "排序规则(pos默认,rev)"
// @Param   onlylz   path  bool  false  "只看楼主"
// @Success 200 {object} outobjs.OutPostPagedList
// @router /post/list [get]
func (c *GroupController) GetPosts() {
	current_uid := c.CurrentUid()
	threadid, _ := c.GetInt64("thread_id")
	page, _ := c.GetInt("page")
	orderby := c.GetString("orderby")
	onlylz, _ := c.GetBool("onlylz")

	oby := groups.POST_ORDERBY_ASC
	if orderby == "rev" {
		oby = groups.POST_ORDERBY_DESC
	}
	if page <= 0 {
		page = 1
	}
	size := 20
	cfg := groups.GetDefaultCfg()
	ps := groups.NewPostService(cfg)
	ths := groups.NewThreadService(cfg)
	gs := groups.NewGroupService(cfg)
	thread := ths.Get(threadid)
	if thread == nil || thread.Closed {
		c.Json(libs.NewError("group_getpost_fail", "GP1410", "帖子不存在", ""))
		return
	}
	//用户顶踩记录
	actionTags := ps.MemberThreadPostActionTags(threadid, current_uid)
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
	total, posts := ps.Gets(threadid, page, size, oby, onlylz)

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
		PageSize:    size,
		Posts:       out_posts,
		MaxDingPost: maxding,
		MaxCaiPost:  nil,
		Thread:      outobjs.GetOutThread(thread),
		JoinedGroup: gs.IsJoined(current_uid, thread.GroupId),
		LordPost:    lordOutPost,
	}
	c.Json(out_p)
}

// @Title 创建评论
// @Description 创建评论
// @Param   access_token  path  string  true  "access_token"
// @Param   thread_id   path  int  true  "帖子id"
// @Param   subject   path  string  false  "标题"
// @Param   content   path  string  true  "内容"
// @Param   img_ids   path  string  false  "图片集(最大9张 逗号,分隔)"
// @Param   replyid   path  string  false  "回复id"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Param   fromdev   path  string  false  "设备标识(android,ios,ipad,wphone,web)"
// @Success 200 {object} libs.Error
// @router /post/submit [post]
func (c *GroupController) CreatePost() {
	current_uid := c.CurrentUid()
	threadid, _ := c.GetInt64("thread_id")
	subject, _ := utils.UrlDecode(c.GetString("subject"))
	content, _ := utils.UrlDecode(c.GetString("content"))
	imgids := c.GetString("img_ids")
	replyid := c.GetString("replyid")
	longitude, _ := c.GetFloat("longitude")
	latitude, _ := c.GetFloat("latitude")
	fromdev := c.GetString("fromdev")

	cfg := groups.GetDefaultCfg()
	ps := groups.NewPostService(cfg)

	img_ids := []int64{}
	arrImg := strings.Split(imgids, ",")
	for _, _ai := range arrImg {
		_id, _ := strconv.ParseInt(_ai, 10, 64)
		if _id > 0 {
			img_ids = append(img_ids, _id)
		}
	}
	if len(img_ids) > 9 {
		c.Json(libs.NewError("group_newpost_fail", "GP1430", "图片数量不能大于9张", ""))
		return
	}
	post := &groups.Post{
		ThreadId:   threadid,
		AuthorId:   current_uid,
		Subject:    subject,
		Message:    content,
		Ip:         c.Ctx.Input.IP(),
		FromDevice: groups.GetFromDevice(fromdev),
		ImgIds:     img_ids,
		ReplyId:    replyid,
		LongiTude:  float32(longitude),
		LatiTude:   float32(latitude),
	}
	err := ps.Create(post)
	if err != nil {
		c.Json(libs.NewError("group_newpost_fail", "GP1431", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_newpost_success", RESPONSE_SUCCESS, "评论成功", ""))
}

// @Title 顶踩评论
// @Description 顶踩评论
// @Param   access_token  path  string  true  "access_token"
// @Param   post_id   path  string  true  "评论id"
// @Param   action   path  string  true  "动作(ding,cancel_ding)"
// @Success 200 {object} libs.Error
// @router /post/action [post]
func (c *GroupController) ActionPost() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_actionpost_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能操作", ""))
		return
	}
	postid := c.GetString("post_id")
	action := c.GetString("action")
	var act groups.POST_ACTION
	if action == "ding" {
		act = groups.POST_ACTION_DING
	}
	if action == "cancel_ding" {
		act = groups.POST_ACTION_CANCEL_DING
	}
	if int(act) == 0 {
		c.Json(libs.NewError("group_actionpost_fail", "GP1500", "action操作符不被支持", ""))
		return
	}
	cfg := groups.GetDefaultCfg()
	ps := groups.NewPostService(cfg)
	err := ps.Action(postid, current_uid, act)
	if err != nil {
		c.Json(libs.NewError("group_actionpost_fail", "GP1501", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_actionpost_success", RESPONSE_SUCCESS, "操作成功", ""))
}

// @Title 举报选项
// @Description 举报选项(字符串数组)
// @Success 200
// @router /report/options [get]
func (c *GroupController) ReportOptions() {
	c.Json([]string{
		"广告或垃圾信息",
		"色情、淫秽或低俗内容",
		"激进时政或意识形态话题",
		"人身攻击或文字侮辱",
	})
}

// @Title 举报
// @Description 举报
// @Param   access_token  path  string  true  "access_token"
// @Param   refid   path  string  true  "关联id"
// @Param   c   path  int  true  "关联id的类型"
// @Param   msg  path  string  false  "举报内容"
// @Success 200 {object} libs.Error
// @router /report [post]
func (c *GroupController) Report() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_report_premission_denied", UNAUTHORIZED_CODE, "必须登录后才能操作", ""))
		return
	}
	refid := c.GetString("refid")
	cc, _ := c.GetInt("c")
	msg, _ := utils.UrlDecode(c.GetString("msg"))
	if cc != int(groups.REPORT_CATEGORY_POST) &&
		cc != int(groups.REPORT_CATEGORY_THREAD) &&
		cc != int(groups.REPORT_CATEGORY_GROUP) {
		c.Json(libs.NewError("group_report_fail", "GP1551", "举报类型不存在", ""))
		return
	}
	report := &groups.Report{
		RefId:   refid,
		C:       groups.REPORT_CATEGORY(cc),
		Ts:      time.Now().Unix(),
		RefTxt:  "",
		PostUid: current_uid,
		Msg:     msg,
	}
	rs := groups.NewReportService()
	err := rs.Create(report)
	if err != nil {
		c.Json(libs.NewError("group_report_fail", "GP1552", err.Error(), ""))
		return
	}
	c.Json(libs.NewError("group_report_success", RESPONSE_SUCCESS, "提交成功", ""))
}

// @Title 分享
// @Description 分享
// @Param   access_token  path   string  true  "access_token"
// @Param   thread_id     path   int  false  "帖子id"
// @Param   text     path   string  false  "文本内容"
// @Param   ref_uids     path   string  false  "提到了哪些好友uid(示例:1001,1002,1023))"
// @Success 200 成功返回error_code:REP000
// @router /thread/share [post]
func (c *GroupController) Share() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_share_premission_denied", UNAUTHORIZED_CODE, "请登陆后重新尝试", ""))
		return
	}
	thread_id, _ := c.GetInt64("thread_id")
	text, _ := utils.UrlDecode(c.GetString("text"))
	ref_uids := c.GetString("ref_uids")

	if thread_id <= 0 {
		c.Json(libs.NewError("group_share_fail", "GP1600", "必须提供帖子id", ""))
		return
	}
	ts := groups.NewThreadService(groups.GetDefaultCfg())
	thread := ts.Get(thread_id)
	if thread == nil {
		c.Json(libs.NewError("group_share_fail", "GP1601", "帖子不存在", ""))
		return
	}
	nts := &share.Shares{}
	_share_res := []string{}
	_share_res = append(_share_res, nts.TranformResource(groups.SHARE_KIND_GROUP, fmt.Sprintf("%d", thread_id)))
	_s := share.Share{
		Uid:       current_uid,
		Source:    "",
		Text:      text,
		ShareType: int(groups.SHARE_KIND_GROUP),
		Type:      share.SHARE_TYPE_ORIGINAL, //默认是原创
		Resources: nts.CombinResources(_share_res),
		RefUids:   ref_uids,
	}
	err, _id := nts.Create(&_s, false)
	if err != nil {
		c.Json(libs.NewError("group_share_create", "GP1602", "发布分享失败", ""))
		return
	}
	c.Json(libs.NewError("group_share_succ", RESPONSE_SUCCESS, "分享成功", strconv.FormatInt(_id, 10)))
}

// @Title 组消息计数
// @Description 组消息计数
// @Param   access_token  path   string  true  "access_token"
// @Success 200 成功返回error_code:REP000,error_description=数量
// @router /msg/c [get]
func (c *GroupController) MsgCount() {
	current_uid := c.CurrentUid()
	if current_uid <= 0 {
		c.Json(libs.NewError("group_msgcount_premission_denied", UNAUTHORIZED_CODE, "请登陆后重新尝试", ""))
		return
	}
	gcm := gms.NewEventCount(current_uid)
	if gcm > 99 {
		gcm = 99
	}
	c.Json(libs.NewError("group_msgcount_succ", RESPONSE_SUCCESS, fmt.Sprintf("%d", gcm), ""))
}
