package controllers

// 组 API
type GroupController struct {
	BaseController
}

func (f *GroupController) Prepare() {
}

func (f *GroupController) URLMapping() {
	f.Mapping("GetSetting", f.GetSetting)
	f.Mapping("GetGroup", f.GetGroup)
	f.Mapping("GetGroups", f.GetGroups)
	f.Mapping("GetRecruitingGroups", f.GetRecruitingGroups)
	f.Mapping("GetMyGroups", f.GetMyGroups)
	f.Mapping("GetMyJoinGroups", f.GetMyJoinGroups)
	f.Mapping("CreateGroup", f.CreateGroup)
	f.Mapping("UpdateGroup", f.UpdateGroup)
	f.Mapping("JoinGroup", f.JoinGroup)
	f.Mapping("ExitGroup", f.ExitGroup)
	f.Mapping("GetThreads", f.GetThreads)
	f.Mapping("GetThreads", f.GetThreads)
	f.Mapping("GetPosts", f.GetPosts)
	f.Mapping("CreatePost", f.CreatePost)
	f.Mapping("ReportOptions", f.ReportOptions)
	f.Mapping("Report", f.Report)
}

// @Title 申请组的设定值
// @Description 申请组的设定值
// @Success 200 {object} outobjs.OutGroupSetting
// @router /group/setting [get]
func (f *GroupController) GetSetting() {

}

// @Title 获取组信息
// @Description 获取组信息
// @Param   access_token  path  string  false  "access_token"
// @Success 200 {object} outobjs.OutGroup
// @router /group/get [get]
func (f *GroupController) GetGroup() {

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
func (f *GroupController) GetGroups() {

}

// @Title 获取招募中的组列表
// @Description 获取招募中的组列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   game_ids   path  string  false  "游戏ids(逗号,分隔)"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutGroupPagedList
// @router /group/recruiting [get]
func (f *GroupController) GetRecruitingGroups() {

}

// @Title 用户创建的组列表
// @Description 用户创建的组列表(返回数组)
// @Param   access_token  path  string  true  "access_token"
// @Success 200 {object} outobjs.OutMyGroups
// @router /group/my [get]
func (f *GroupController) GetMyGroups() {

}

// @Title 获取我加入的组列表
// @Description 获取我加入的组列表(返回数组)
// @Param   access_token  path  string  true  "access_token"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutGroupPagedList
// @router /group/myjoins [get]
func (f *GroupController) GetMyJoinGroups() {

}

// @Title 申请建组
// @Description 申请建组
// @Param   access_token  path  string  true  "access_token"
// @Param   name   path  string  true  "名称"
// @Param   description   path  string  true  "描述"
// @Param   country  path  string  false  "国家"
// @Param   city  path  string  false  "城市"
// @Param   game_ids  path  string  false  "选择游戏(逗号,分隔)"
// @Param   img   path  int  false  "图片id"
// @Param   bgimg   path  int  true  "背景图片id"
// @Param   longitude   path  float  false  "经度"
// @Param   latitude   path  float  false  "维度"
// @Success 200 {object} libs.Error
// @router /group/create [post]
func (f *GroupController) CreateGroup() {

}

// @Title 建组前校验用户是否满足条件
// @Description 建组前校验用户是否满足条件
// @Param   access_token  path  string  true  "access_token"
// @Success 200 {object} libs.Error
// @router /group/create_check [post]
func (f *GroupController) CreateGroupCheck() {

}

// @Title 更新组属性
// @Description 更新组属性
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Param   description   path  string  true  "描述"
// @Param   img   path  int  true  "图片id"
// @Param   bgimg   path  int  true  "背景图片id"
// @Success 200 {object} libs.Error
// @router /group/update [post]
func (f *GroupController) UpdateGroup() {

}

// @Title 加入组
// @Description 加入组
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Success 200 {object} libs.Error
// @router /group/join [post]
func (f *GroupController) JoinGroup() {

}

// @Title 离开组
// @Description 离开组
// @Param   access_token  path  string  true  "access_token"
// @Param   groupid   path  int  true  "组id"
// @Success 200 {object} libs.Error
// @router /group/exit [post]
func (f *GroupController) ExitGroup() {

}

// @Title 获取组帖子
// @Description 获取组帖子列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   group_id   path  int  true  "组id"
// @Param   page   path  int  false  "页"
// @Success 200 {object} outobjs.OutThreadPagedList
// @router /thread/list [get]
func (f *GroupController) GetThreads() {

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
func (f *GroupController) CreateThread() {

}

// @Title 帖子评论列表
// @Description 帖子评论列表(返回数组)
// @Param   access_token  path  string  false  "access_token"
// @Param   thread_id   path  int  true  "帖子id"
// @Param   orderby   path  string  false  "排序规则(pos默认,rev)"
// @Success 200 {object} outobjs.OutPostPagedList
// @router /post/list [get]
func (f *GroupController) GetPosts() {

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
func (f *GroupController) CreatePost() {

}

// @Title 举报选项
// @Description 举报选项(字符串数组)
// @Success 200
// @router /report/options [get]
func (f *GroupController) ReportOptions() {

}

// @Title 举报
// @Description 举报
// @Param   access_token  path  string  true  "access_token"
// @Param   refid   path  string  true  "关联id"
// @Param   c   path  int  true  "关联id的类型"
// @Param   thread_id   path  int  true  "帖子id"
// @Param   msg  path  string  false  "举报内容"
// @Success 200 {object} libs.Error
// @router /report [post]
func (f *GroupController) Report() {

}
