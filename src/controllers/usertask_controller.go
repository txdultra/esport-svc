package controllers

import (
	"libs"
	"libs/utask/client"
	"libs/utask/proxy"
	"outobjs"
)

// 用户任务模块 API
type UserTaskController struct {
	AuthorizeController
}

func (c *UserTaskController) Prepare() {
	c.AuthorizeController.Prepare()
}

func (c *UserTaskController) URLMapping() {
	c.Mapping("All", c.All)
}

// @Title 获取任务列表
// @Description 获取任务列表(数组)
// @Param   access_token   path   string  true  "access_token"
// @Success 200 {object} outobjs.OutUserTaskGroup
// @router /all [get]
func (c *UserTaskController) All() {
	uid := c.CurrentUid()
	ut_client, transport, err := client.NewClient()
	defer func() {
		if transport != nil {
			transport.Close()
		}
	}()
	if err != nil {
		c.Json(libs.NewError("member_task_getall_fail", "UT1001", "系统错误", ""))
		return
	}
	mgroups, err := ut_client.GetMissions(uid)
	if err != nil {
		c.Json(libs.NewError("member_task_getall_fail", "UT1002", "系统错误", ""))
		return
	}

	taskGroups := []*outobjs.OutUserTaskGroup{}
	for _, mg := range mgroups {
		tasks := make([]*outobjs.OutUserTask, len(mg.Tasks), len(mg.Tasks))
		for i, mt := range mg.Tasks {
			if mt.Task == nil {
				continue
			}
			tasks[i] = &outobjs.OutUserTask{
				TaskId:       mt.Task.TaskId,
				TaskType:     proxy.ConvertTASK_TYPE(mt.Task.TaskType),
				GroupId:      mt.Task.GroupId,
				Name:         mt.Task.Name,
				Description:  mt.Task.Description,
				Icon:         mt.Task.Icon,
				IconUrl:      file_storage.GetFileUrl(mt.Task.Icon),
				Limit:        int(mt.Task.TaskLimits),
				Dones:        int(mt.Count),
				Period:       int(mt.Task.Period),
				PeriodType:   proxy.ConvertTASK_PERIOD_TYPE(mt.Task.PeriodType),
				Reward:       proxy.ConvertTASK_REWARD_TYPE(mt.Task.Reward),
				Prize:        mt.Task.Prize,
				DisplayOrder: int(mt.Task.DisplayOrder),
				ResetVar:     int(mt.Task.ResetVar),
				ResetTime:    int(mt.Task.ResetTime),
				Ex1:          mt.Task.Ex1,
				Ex2:          mt.Task.Ex2,
				Ex3:          mt.Task.Ex3,
				Ex4:          mt.Task.Ex4,
				Ex5:          mt.Task.Ex5,
			}
		}
		if mg.Group == nil {
			continue
		}
		taskGroups = append(taskGroups, &outobjs.OutUserTaskGroup{
			GroupId:      mg.Group.GroupId,
			GroupType:    proxy.ConvertGROUP_TYPE(mg.Group.GroupType),
			TaskType:     proxy.ConvertTASK_TYPE(mg.Group.TaskType),
			Name:         mg.Group.Name,
			Description:  mg.Group.Description,
			Icon:         mg.Group.Icon,
			IconUrl:      file_storage.GetFileUrl(mg.Group.Icon),
			BgImg:        mg.Group.BgImg,
			BgImgUrl:     file_storage.GetFileUrl(mg.Group.BgImg),
			TaskCount:    len(mg.Tasks),
			DoneCount:    int(mg.Dones),
			DisplayOrder: int(mg.Group.DisplayOrder),
			Tasks:        tasks,
			Ex1:          mg.Group.Ex1,
			Ex2:          mg.Group.Ex2,
			Ex3:          mg.Group.Ex3,
			Ex4:          mg.Group.Ex4,
			Ex5:          mg.Group.Ex5,
		})
	}
	c.Json(taskGroups)
}
