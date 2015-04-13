package admincp

import (
	"controllers"
	"libs"
	ut_client "libs/utask/client"
	"libs/utask/proxy"
	"utils"
)

// 用户任务 API
type UTaskCPController struct {
	AdminController
}

func (c *UTaskCPController) Prepare() {
	c.AdminController.Prepare()
}

// @Title 添加任务组
// @Description 添加任务组
// @Param   gt   path	int true  "组类型"
// @Param   tt   path	int true  "任务类型"
// @Param   name   path	string true  "名称"
// @Param   description   path	string false  "描述"
// @Param   icon   path	int true  "图标"
// @Param   bgimg   path int true  "背景图"
// @Param   enabled path bool true  "开启"
// @Param   starttime path int true "开始时间戳"
// @Param   endtime path int true "结束时间戳"
// @Param   displayorder   path	int false  "排序"
// @Param   ex1   path	string false  "扩展字段1"
// @Param   ex2   path	string false  "扩展字段2"
// @Param   ex3   path	string false  "扩展字段3"
// @Param   ex4   path	string false  "扩展字段4"
// @Param   ex5   path	string false  "扩展字段5"
// @Success 200 {object} libs.Error
// @router /add_group [post]
func (c *UTaskCPController) AddGroup() {
	gt, _ := c.GetInt("gt")
	tt, _ := c.GetInt("tt")
	name, _ := utils.UrlDecode(c.GetString("name"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	icon, _ := c.GetInt64("icon")
	bgimg, _ := c.GetInt64("bgimg")
	enabled, _ := c.GetBool("enabled")
	starttime, _ := c.GetInt64("starttime")
	endtime, _ := c.GetInt64("endtime")
	displayorder, _ := c.GetInt("displayorder")
	ex1, _ := utils.UrlDecode(c.GetString("ex1"))
	ex2, _ := utils.UrlDecode(c.GetString("ex2"))
	ex3, _ := utils.UrlDecode(c.GetString("ex3"))
	ex4, _ := utils.UrlDecode(c.GetString("ex4"))
	ex5, _ := utils.UrlDecode(c.GetString("ex5"))

	if gt <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_001", "参数错误gt", ""))
		return
	}
	if tt <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_002", "参数错误tt", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_003", "参数错误name", ""))
		return
	}
	if bgimg <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_004", "参数错误bgimg", ""))
		return
	}
	if starttime <= 0 || endtime <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_005", "参数错误starttime或endtime", ""))
		return
	}

	group := &proxy.TaskGroup{}
	group.GroupType = proxy.TASK_GROUP_TYPE(gt)
	group.TaskType = proxy.TASK_TYPE(tt)
	group.Name = name
	group.Description = description
	group.Icon = icon
	group.BgImg = bgimg
	group.Enabled = enabled
	group.StartTime = starttime
	group.EndTime = endtime
	group.DisplayOrder = int32(displayorder)
	group.Ex1 = ex1
	group.Ex2 = ex2
	group.Ex3 = ex3
	group.Ex4 = ex4
	group.Ex5 = ex5

	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_006", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	result, err := client.CreateGroup(group)
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_007", "连接任务系统错误02:"+err.Error(), ""))
		return
	}
	if result.Success {
		c.Json(libs.NewError("admincp_user_taskgroup_add_succ", controllers.RESPONSE_SUCCESS, "新任务组添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_user_taskgroup_add_fail", "GM020_008", "添加失败:"+result.Exception, ""))
}

// @Title 更新任务组
// @Description 更新任务组
// @Param   groupid   path	int true  "组id"
// @Param   name   path	string true  "名称"
// @Param   description   path	string false  "描述"
// @Param   icon   path	int true  "图标"
// @Param   bgimg   path int true  "背景图"
// @Param   enabled path bool true  "开启"
// @Param   starttime path int true "开始时间戳"
// @Param   endtime path int true "结束时间戳"
// @Param   displayorder   path	int false  "排序"
// @Param   ex1   path	string false  "扩展字段1"
// @Param   ex2   path	string false  "扩展字段2"
// @Param   ex3   path	string false  "扩展字段3"
// @Param   ex4   path	string false  "扩展字段4"
// @Param   ex5   path	string false  "扩展字段5"
// @Success 200 {object} libs.Error
// @router /update_group [post]
func (c *UTaskCPController) UpdateGroup() {
	groupid, _ := c.GetInt64("groupid")
	name, _ := utils.UrlDecode(c.GetString("name"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	icon, _ := c.GetInt64("icon")
	bgimg, _ := c.GetInt64("bgimg")
	enabled, _ := c.GetBool("enabled")
	starttime, _ := c.GetInt64("starttime")
	endtime, _ := c.GetInt64("endtime")
	displayorder, _ := c.GetInt("displayorder")
	ex1, _ := utils.UrlDecode(c.GetString("ex1"))
	ex2, _ := utils.UrlDecode(c.GetString("ex2"))
	ex3, _ := utils.UrlDecode(c.GetString("ex3"))
	ex4, _ := utils.UrlDecode(c.GetString("ex4"))
	ex5, _ := utils.UrlDecode(c.GetString("ex5"))

	if groupid <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_011", "参数错误groupid", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_012", "参数错误name", ""))
		return
	}
	if bgimg <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_013", "参数错误bgimg", ""))
		return
	}
	if starttime <= 0 || endtime <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_014", "参数错误starttime或endtime", ""))
		return
	}

	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_015", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	group, err := client.GetGroup(groupid)
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_016", "连接任务系统错误02:"+err.Error(), ""))
		return
	}

	group.Name = name
	group.Description = description
	group.Icon = icon
	group.BgImg = bgimg
	group.Enabled = enabled
	group.StartTime = starttime
	group.EndTime = endtime
	group.DisplayOrder = int32(displayorder)
	group.Ex1 = ex1
	group.Ex2 = ex2
	group.Ex3 = ex3
	group.Ex4 = ex4
	group.Ex5 = ex5
	result, err := client.UpdateGroup(group)
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_017", "连接任务系统错误03:"+err.Error(), ""))
		return
	}
	if result.Success {
		c.Json(libs.NewError("admincp_user_taskgroup_update_succ", controllers.RESPONSE_SUCCESS, "更新任务组成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_user_taskgroup_update_fail", "GM020_018", "更新失败:"+result.Exception, ""))
}

// @Title 删除任务组
// @Description 更新任务组
// @Param   groupid   path	int true  "任务组id"
// @Success 200 {object} libs.Error
// @router /del_group [delete]
func (c *UTaskCPController) DelGroup() {
	groupid, _ := c.GetInt64("groupid")
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_del_fail", "GM020_021", "参数错误groupid", ""))
		return
	}

	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	result, err := client.DeleteGroup(groupid)
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_del_fail", "GM020_022", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	if result.Success {
		c.Json(libs.NewError("admincp_user_taskgroup_del_succ", controllers.RESPONSE_SUCCESS, "删除任务组成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_user_taskgroup_del_fail", "GM020_023", "删除失败:"+result.Exception, ""))
}

// @Title 获取所有任务组
// @Description 获取所有任务组(数组)
// @Success 200
// @router /get_group_all [get]
func (c *UTaskCPController) GetAllGroups() {
	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	groups, err := client.GetGroups()
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_getall_fail", "GM020_100", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	c.Json(groups)
}

// @Title 获取指定任务组
// @Description 获取指定任务组
// @Param   groupid   path	int true  "任务组id"
// @Success 200
// @router /get_group [get]
func (c *UTaskCPController) GetGroup() {
	groupid, _ := c.GetInt64("groupid")
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_user_taskgroup_get_fail", "GM020_110", "参数错误groupid", ""))
		return
	}
	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	group, err := client.GetGroup(groupid)
	if err != nil {
		c.Json(libs.NewError("admincp_user_taskgroup_get_fail", "GM020_111", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	c.Json(group)
}

// @Title 添加新任务
// @Description 添加新任务
// @Param   groupid   path	int true  "任务组id"
// @Param   tt   path	int true  "任务类型"
// @Param   name   path	string true  "名称"
// @Param   description   path	string false  "描述"
// @Param   icon   path	int true  "图标"
// @Param   tasklimits   path int true  "任务操作计数"
// @Param   applyperm   path	string false  "获取任务条件"
// @Param   eventname   path    string true  "事件名称"
// @Param   starttime path int true "开始时间戳"
// @Param   endtime path int true "结束时间戳"
// @Param   resettime path int true "重置时间点"
// @Param   period path int true "重置时间计数"
// @Param   periodtype path int true "重置时间类型"
// @Param   reward path int true "奖励类型"
// @Param   prize path int true "奖励积分数"
// @Param   version   path string false  "版本"
// @Param   displayorder   path	int false  "排序"
// @Param   ex1   path	string false  "扩展字段1"
// @Param   ex2   path	string false  "扩展字段2"
// @Param   ex3   path	string false  "扩展字段3"
// @Param   ex4   path	string false  "扩展字段4"
// @Param   ex5   path	string false  "扩展字段5"
// @Success 200 {object} libs.Error
// @router /add_task [post]
func (c *UTaskCPController) AddTask() {
	uid := c.CurrentUid()
	groupid, _ := c.GetInt64("groupid")
	tt, _ := c.GetInt("tt")
	name, _ := utils.UrlDecode(c.GetString("name"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	icon, _ := c.GetInt64("icon")
	tasklimits, _ := c.GetInt("tasklimits")
	applyperm, _ := utils.UrlDecode(c.GetString("applyperm"))
	eventname, _ := utils.UrlDecode(c.GetString("eventname"))
	starttime, _ := c.GetInt64("starttime")
	endtime, _ := c.GetInt64("endtime")
	resettime, _ := c.GetInt("resettime")
	period, _ := c.GetInt("period")
	periodtype, _ := c.GetInt("periodtype")
	reward, _ := c.GetInt("reward")
	prize, _ := c.GetInt64("prize")
	version, _ := utils.UrlDecode(c.GetString("version"))
	displayorder, _ := c.GetInt("displayorder")
	ex1, _ := utils.UrlDecode(c.GetString("ex1"))
	ex2, _ := utils.UrlDecode(c.GetString("ex2"))
	ex3, _ := utils.UrlDecode(c.GetString("ex3"))
	ex4, _ := utils.UrlDecode(c.GetString("ex4"))
	ex5, _ := utils.UrlDecode(c.GetString("ex5"))

	if groupid <= 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_031", "参数错误groupid", ""))
		return
	}
	if tt <= 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_032", "参数错误tt", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_033", "参数错误name", ""))
		return
	}
	if starttime <= 0 || endtime <= 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_034", "参数错误starttime或endtime", ""))
		return
	}
	if tasklimits <= 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_035", "参数错误tasklimits", ""))
		return
	}
	if len(eventname) == 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_036", "参数错误eventname", ""))
		return
	}
	if periodtype <= 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_037", "参数错误periodtype", ""))
		return
	}
	if reward <= 0 {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_038", "参数错误reward", ""))
		return
	}

	task := &proxy.Task{}
	task.GroupId = groupid
	task.TaskType = proxy.TASK_TYPE(tt)
	task.Name = name
	task.Description = description
	task.Icon = icon
	task.TaskLimits = int32(tasklimits)
	task.ApplyPerm = applyperm
	task.EventName = eventname
	task.StartTime = starttime
	task.EndTime = endtime
	task.ResetTime = int32(resettime)
	task.Period = int32(period)
	task.PeriodType = proxy.TASK_PERIOD_TYPE(periodtype)
	task.Reward = proxy.TASK_REWARD_TYPE(reward)
	task.Prize = prize
	task.Applicants = uid
	task.Version = version
	task.DisplayOrder = int32(displayorder)
	task.Ex1 = ex1
	task.Ex2 = ex2
	task.Ex3 = ex3
	task.Ex4 = ex4
	task.Ex5 = ex5

	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_039", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	result, err := client.CreateTask(task)
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_040", "连接任务系统错误02:"+err.Error(), ""))
		return
	}
	if result.Success {
		c.Json(libs.NewError("admincp_user_task_add_succ", controllers.RESPONSE_SUCCESS, "新任务添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_user_task_add_fail", "GM020_041", "添加失败:"+result.Exception, ""))
}

// @Title 更新任务
// @Description 更新任务
// @Param   taskid   path	int true  "任务id"
// @Param   groupid   path	int true  "任务组id"
// @Param   name   path	string true  "名称"
// @Param   description   path	string false  "描述"
// @Param   icon   path	int true  "图标"
// @Param   tasklimits   path int true  "任务操作计数"
// @Param   applyperm   path	string false  "获取任务条件"
// @Param   starttime path int true "开始时间戳"
// @Param   endtime path int true "结束时间戳"
// @Param   resettime path int true "重置时间点"
// @Param   period path int true "重置时间计数"
// @Param   periodtype path int true "重置时间类型"
// @Param   reward path int true "奖励类型"
// @Param   prize path int true "奖励积分数"
// @Param   version   path string false  "版本"
// @Param   displayorder   path	int false  "排序"
// @Param   ex1   path	string false  "扩展字段1"
// @Param   ex2   path	string false  "扩展字段2"
// @Param   ex3   path	string false  "扩展字段3"
// @Param   ex4   path	string false  "扩展字段4"
// @Param   ex5   path	string false  "扩展字段5"
// @Success 200 {object} libs.Error
// @router /update_task [post]
func (c *UTaskCPController) UpdateTask() {
	uid := c.CurrentUid()
	taskid, _ := c.GetInt64("taskid")
	groupid, _ := c.GetInt64("groupid")
	name, _ := utils.UrlDecode(c.GetString("name"))
	description, _ := utils.UrlDecode(c.GetString("description"))
	icon, _ := c.GetInt64("icon")
	tasklimits, _ := c.GetInt("tasklimits")
	applyperm, _ := utils.UrlDecode(c.GetString("applyperm"))
	starttime, _ := c.GetInt64("starttime")
	endtime, _ := c.GetInt64("endtime")
	resettime, _ := c.GetInt("resettime")
	period, _ := c.GetInt("period")
	periodtype, _ := c.GetInt("periodtype")
	reward, _ := c.GetInt("reward")
	prize, _ := c.GetInt64("prize")
	version, _ := utils.UrlDecode(c.GetString("version"))
	displayorder, _ := c.GetInt("displayorder")
	ex1, _ := utils.UrlDecode(c.GetString("ex1"))
	ex2, _ := utils.UrlDecode(c.GetString("ex2"))
	ex3, _ := utils.UrlDecode(c.GetString("ex3"))
	ex4, _ := utils.UrlDecode(c.GetString("ex4"))
	ex5, _ := utils.UrlDecode(c.GetString("ex5"))

	if taskid <= 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_050", "参数错误taskid", ""))
		return
	}
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_051", "参数错误groupid", ""))
		return
	}
	if len(name) == 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_052", "参数错误name", ""))
		return
	}
	if starttime <= 0 || endtime <= 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_053", "参数错误starttime或endtime", ""))
		return
	}
	if tasklimits <= 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_054", "参数错误tasklimits", ""))
		return
	}
	if periodtype <= 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_055", "参数错误periodtype", ""))
		return
	}
	if reward <= 0 {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_056", "参数错误reward", ""))
		return
	}

	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_057", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	task, err := client.GetTask(taskid)
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_058", "连接任务系统错误02:"+err.Error(), ""))
		return
	}
	task.GroupId = groupid
	task.Name = name
	task.Description = description
	task.Icon = icon
	task.TaskLimits = int32(tasklimits)
	task.ApplyPerm = applyperm
	task.StartTime = starttime
	task.EndTime = endtime
	task.ResetTime = int32(resettime)
	task.Period = int32(period)
	task.PeriodType = proxy.TASK_PERIOD_TYPE(periodtype)
	task.Reward = proxy.TASK_REWARD_TYPE(reward)
	task.Prize = prize
	task.Applicants = uid
	task.Version = version
	task.DisplayOrder = int32(displayorder)
	task.Ex1 = ex1
	task.Ex2 = ex2
	task.Ex3 = ex3
	task.Ex4 = ex4
	task.Ex5 = ex5
	result, err := client.UpdateTask(task)
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_059", "连接任务系统错误03:"+err.Error(), ""))
		return
	}
	if result.Success {
		c.Json(libs.NewError("admincp_user_task_update_succ", controllers.RESPONSE_SUCCESS, "更新任务添加成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_user_task_update_fail", "GM020_060", "更新失败:"+result.Exception, ""))
}

// @Title 删除任务
// @Description 删除任务
// @Param   taskid   path	int true  "任务id"
// @Success 200 {object} libs.Error
// @router /del_task [delete]
func (c *UTaskCPController) DelTask() {
	taskid, _ := c.GetInt64("taskid")
	if taskid <= 0 {
		c.Json(libs.NewError("admincp_user_task_del_fail", "GM020_070", "参数错误taskid", ""))
		return
	}

	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	task, err := client.GetTask(taskid)
	if err != nil || task == nil {
		c.Json(libs.NewError("admincp_user_task_del_fail", "GM020_071", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	result, err := client.DeleteTask(taskid, task.GroupId)
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_del_fail", "GM020_072", "连接任务系统错误02:"+err.Error(), ""))
		return
	}
	if result.Success {
		c.Json(libs.NewError("admincp_user_task_del_succ", controllers.RESPONSE_SUCCESS, "删除任务成功", ""))
		return
	}
	c.Json(libs.NewError("admincp_user_task_del_fail", "GM020_073", "删除失败:"+result.Exception, ""))
}

// @Title 获取组任务
// @Description 获取组任务(数组)
// @Param   groupid   path	int true  "任务组id"
// @Success 200
// @router /get_tasks [get]
func (c *UTaskCPController) GetTasks() {
	groupid, _ := c.GetInt64("groupid")
	if groupid <= 0 {
		c.Json(libs.NewError("admincp_user_task_gets_fail", "GM020_150", "参数错误groupid", ""))
		return
	}
	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	tasks, err := client.GetTasks(groupid)
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_gets_fail", "GM020_151", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	c.Json(tasks)
}

// @Title 获取任务
// @Description 获取任务
// @Param   taskid   path	int true  "任务id"
// @Success 200
// @router /get_task [get]
func (c *UTaskCPController) GetTask() {
	taskid, _ := c.GetInt64("taskid")
	if taskid <= 0 {
		c.Json(libs.NewError("admincp_user_task_get_fail", "GM020_160", "参数错误taskid", ""))
		return
	}
	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	task, err := client.GetTask(taskid)
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_get_fail", "GM020_161", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	c.Json(task)
}

// @Title 获取事件列表
// @Description 获取事件列表(数组)
// @Success 200
// @router /get_eventnames [get]
func (c *UTaskCPController) GetEventNames() {
	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	events, err := client.GetEventNames()
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_getevents_fail", "GM020_170", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	c.Json(events)
}

// @Title 获取任务重置定时器
// @Description 获取任务重置定时器(数组)
// @Success 200
// @router /get_tasktimers [get]
func (c *UTaskCPController) GetTaskTimers() {
	client, transport, err := ut_client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	timers, err := client.GetTaskTimers()
	if err != nil {
		c.Json(libs.NewError("admincp_user_task_gettimers_fail", "GM020_180", "连接任务系统错误01:"+err.Error(), ""))
		return
	}
	c.Json(timers)
}
