package controllers

import (
	"fmt"
	"libs"
	"libs/groups"
	"outobjs"
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
	c.Mapping("GetThreads", c.GetThreads)
	c.Mapping("GetThreads", c.GetThreads)
	c.Mapping("GetPosts", c.GetPosts)
	c.Mapping("CreatePost", c.CreatePost)
	c.Mapping("ReportOptions", c.ReportOptions)
	c.Mapping("Report", c.Report)
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

}

// @Title 获取招募中的组列表
// @Description 获取招募中的组列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   game_ids   path  string  false  "游戏ids(逗号,分隔)"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutGroupPagedList
// @router /group/recruiting [get]
func (c *GroupController) GetRecruitingGroups() {

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
		MaxAllowGroupCount: cfg.CreateGroupMaxCount,
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
	cfg := groups.GetDefaultCfg()
	gs := groups.NewGroupService(cfg)
	total, joinGroups := gs.MyJoins(current_uid, page, 20)
	out_groups := []*outobjs.OutGroup{}
	for _, gp := range joinGroups {
		out_g := outobjs.GetOutGroup(gp, 0)
		out_g.IsJoined = true //肯定是true
		out_groups = append(out_groups, out_g)
	}
	out_p := &outobjs.OutGroupPagedList{
		CurrentPage: page,
		PageSize:    20,
		TotalPages:  utils.TotalPages(total, 20),
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
	group.BgImg = bgimg
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
	c.Json(out_p)
}

// @Title 新建帖子
// @Description 新建帖子
// @Param   access_token  path  string  false  "access_token"
// @Param   group_id   path  int  true  "组id"
// @Param   subject   path  string  true  "标题"
// @Param   message   path  string  true  "内容"
// @Param   img_ids   path  string  true  "图片集(最大9张 逗号,分隔)"
// @Success 200 {object} libs.Error
// @router /thread/submit [post]
func (c *GroupController) CreateThread() {

}

// @Title 帖子评论列表
// @Description 帖子评论列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   thread_id   path  int  true  "帖子id"
// @Param   orderby   path  string  false  "排序规则(pos默认,rev)"
// @Param   onlylz   path  bool  false  "只看楼主"
// @Success 200 {object} outobjs.OutPostPagedList
// @router /post/list [get]
func (c *GroupController) GetPosts() {

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
// @Success 200 {object} libs.Error
// @router /post/submit [post]
func (c *GroupController) CreatePost() {

}

// @Title 举报选项
// @Description 举报选项(字符串数组)
// @Success 200
// @router /report/options [get]
func (c *GroupController) ReportOptions() {

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

}
