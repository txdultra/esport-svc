package proxy

import (
	"fmt"
	"libs/utask"
)

func NewUserTaskServiceProxy(engine utask.IMissionEngine) *UserTaskServiceProxy {
	utsp := &UserTaskServiceProxy{}
	utsp.tasker = utask.NewTasker()
	utsp.enigne = engine
	return utsp
}

type UserTaskServiceProxy struct {
	tasker *utask.Tasker
	enigne utask.IMissionEngine
}

func ConvertGROUP_TYPE(gt TASK_GROUP_TYPE) utask.TASK_GROUP_TYPE {
	switch gt {
	case TASK_GROUP_TYPE_WEB:
		return utask.TASK_GROUP_TYPE_WEB
	default:
		return utask.TASK_GROUP_TYPE_LIST
	}
}

func UnConvertGROUP_TYPE(gt utask.TASK_GROUP_TYPE) TASK_GROUP_TYPE {
	switch gt {
	case utask.TASK_GROUP_TYPE_WEB:
		return TASK_GROUP_TYPE_WEB
	default:
		return TASK_GROUP_TYPE_LIST
	}
}

func ConvertTASK_TYPE(tt TASK_TYPE) utask.TASK_TYPE {
	switch tt {
	case TASK_TYPE_DUPLICATE:
		return utask.TASK_TYPE_DUPLICATE
	default:
		return utask.TASK_TYPE_ONCE
	}
}

func UnConvertTASK_TYPE(tt utask.TASK_TYPE) TASK_TYPE {
	switch tt {
	case utask.TASK_TYPE_DUPLICATE:
		return TASK_TYPE_DUPLICATE
	default:
		return TASK_TYPE_ONCE
	}
}

func ConvertTASK_PERIOD_TYPE(tpt TASK_PERIOD_TYPE) utask.TASK_PERIOD_TYPE {
	switch tpt {
	case TASK_PERIOD_TYPE_DAY:
		return utask.TASK_PERIOD_TYPE_DAY
	default:
		return utask.TASK_PERIOD_TYPE_NULL
	}
}

func UnConvertTASK_PERIOD_TYPE(tpt utask.TASK_PERIOD_TYPE) TASK_PERIOD_TYPE {
	switch tpt {
	case utask.TASK_PERIOD_TYPE_DAY:
		return TASK_PERIOD_TYPE_DAY
	default:
		return TASK_PERIOD_TYPE_NULL
	}
}

func ConvertTASK_REWARD_TYPE(trt TASK_REWARD_TYPE) utask.TASK_REWARD_TYPE {
	switch trt {
	case TASK_REWARD_TYPE_CREDIT:
		return utask.TASK_REWARD_TYPE_CREDIT
	case TASK_REWARD_TYPE_JING:
		return utask.TASK_REWARD_TYPE_JING
	}
	return utask.TASK_REWARD_TYPE_CREDIT
}

func UnConvertTASK_REWARD_TYPE(trt utask.TASK_REWARD_TYPE) TASK_REWARD_TYPE {
	switch trt {
	case utask.TASK_REWARD_TYPE_CREDIT:
		return TASK_REWARD_TYPE_CREDIT
	case utask.TASK_REWARD_TYPE_JING:
		return TASK_REWARD_TYPE_JING
	}
	return TASK_REWARD_TYPE_CREDIT
}

func UnConvertMISSION_STATUS(status utask.MISSION_STATUS) MISSION_STATUS {
	switch status {
	case utask.MISSION_STATUS_DONE:
		return MISSION_STATUS_DONE
	case utask.MISSION_STATUS_FAILURE:
		return MISSION_STATUS_FAILURE
	default:
		return MISSION_STATUS_UNDONE
	}
}

func UnConvertMissionReward(mr *utask.MissionReward) *MissionReward {
	return &MissionReward{
		RewardType: UnConvertTASK_REWARD_TYPE(mr.RewardType),
		Prize:      mr.Prize,
	}
}

func UnConvertEventResult(result *utask.EventResult, err error) *EventResult_ {
	doneTasks := []*Task{}
	relTasks := []*Task{}
	rewards := []*MissionReward{}
	tasker := utask.NewTasker()
	for _, taskId := range result.DoneTaskIds {
		_task := tasker.GetTask(taskId)
		if _task != nil {
			doneTasks = append(doneTasks, UnConvertTask(_task))
		}
	}

	for _, taskId := range result.RelevanceTaskIds {
		_task := tasker.GetTask(taskId)
		if _task != nil {
			relTasks = append(relTasks, UnConvertTask(_task))
		}
	}

	for _, mr := range result.Rewards {
		rewards = append(rewards, UnConvertMissionReward(mr))
	}
	ex := ""
	if err != nil {
		ex = err.Error()
	}
	return &EventResult_{
		DoneTasks:      doneTasks,
		RelevanceTasks: relTasks,
		Rewards:        rewards,
		Exception:      ex,
	}
}

func ConvertTaskGroup(group *TaskGroup) *utask.TaskGroup {
	return &utask.TaskGroup{
		GroupId:      group.GroupId,
		GroupType:    ConvertGROUP_TYPE(group.GroupType),
		TaskType:     ConvertTASK_TYPE(group.TaskType),
		Name:         group.Name,
		Description:  group.Description,
		Icon:         group.Icon,
		BgImg:        group.BgImg,
		Tasks:        int(group.Tasks),
		Enabled:      group.Enabled,
		StartTime:    group.StartTime,
		EndTime:      group.EndTime,
		DisplayOrder: int(group.DisplayOrder),
		Ex1:          group.Ex1,
		Ex2:          group.Ex2,
		Ex3:          group.Ex3,
		Ex4:          group.Ex4,
		Ex5:          group.Ex5,
	}
}

func UnConvertTaskGroup(group *utask.TaskGroup) *TaskGroup {
	return &TaskGroup{
		GroupId:      group.GroupId,
		GroupType:    UnConvertGROUP_TYPE(group.GroupType),
		TaskType:     UnConvertTASK_TYPE(group.TaskType),
		Name:         group.Name,
		Description:  group.Description,
		Icon:         group.Icon,
		BgImg:        group.BgImg,
		Tasks:        int32(group.Tasks),
		Enabled:      group.Enabled,
		StartTime:    group.StartTime,
		EndTime:      group.EndTime,
		DisplayOrder: int32(group.DisplayOrder),
		Ex1:          group.Ex1,
		Ex2:          group.Ex2,
		Ex3:          group.Ex3,
		Ex4:          group.Ex4,
		Ex5:          group.Ex5,
	}
}

func ConvertTask(task *Task) *utask.Task {
	return &utask.Task{
		TaskId:        task.TaskId,
		TaskType:      ConvertTASK_TYPE(task.TaskType),
		GroupId:       task.GroupId,
		Name:          task.Name,
		Description:   task.Description,
		Icon:          task.Icon,
		TaskLimits:    int(task.TaskLimits),
		ApplyPerm:     task.ApplyPerm,
		EventName:     task.EventName,
		StartTime:     task.StartTime,
		EndTime:       task.EndTime,
		ResetTime:     int(task.ResetTime),
		Period:        int(task.Period),
		PeriodType:    ConvertTASK_PERIOD_TYPE(task.PeriodType),
		Reward:        ConvertTASK_REWARD_TYPE(task.Reward),
		Prize:         task.Prize,
		Applicants:    task.Applicants,
		Achievers:     int(task.Achievers),
		Version:       task.Version,
		DisplayOrder:  int(task.DisplayOrder),
		ResetVar:      int(task.ResetVar),
		LastResetTime: task.LastResetTime,
		Ex1:           task.Ex1,
		Ex2:           task.Ex2,
		Ex3:           task.Ex3,
		Ex4:           task.Ex4,
		Ex5:           task.Ex5,
	}
}

func UnConvertTask(task *utask.Task) *Task {
	return &Task{
		TaskId:        task.TaskId,
		TaskType:      UnConvertTASK_TYPE(task.TaskType),
		GroupId:       task.GroupId,
		Name:          task.Name,
		Description:   task.Description,
		Icon:          task.Icon,
		TaskLimits:    int32(task.TaskLimits),
		ApplyPerm:     task.ApplyPerm,
		EventName:     task.EventName,
		StartTime:     task.StartTime,
		EndTime:       task.EndTime,
		ResetTime:     int32(task.ResetTime),
		Period:        int32(task.Period),
		PeriodType:    UnConvertTASK_PERIOD_TYPE(task.PeriodType),
		Reward:        UnConvertTASK_REWARD_TYPE(task.Reward),
		Prize:         task.Prize,
		Applicants:    task.Applicants,
		Achievers:     int32(task.Achievers),
		Version:       task.Version,
		DisplayOrder:  int32(task.DisplayOrder),
		ResetVar:      int32(task.ResetVar),
		LastResetTime: task.LastResetTime,
		Ex1:           task.Ex1,
		Ex2:           task.Ex2,
		Ex3:           task.Ex3,
		Ex4:           task.Ex4,
		Ex5:           task.Ex5,
	}
}

func UnConvertTaskTimer(tt *utask.TaskTimer) *TaskTimer {
	return &TaskTimer{
		TaskId:       tt.TaskId,
		RunTime:      tt.RunTime,
		PreResetTime: tt.PreResetTime,
	}
}

func UnConvertEventName(en *utask.TaskEvent) *TaskEventName {
	return &TaskEventName{
		Id:          en.Id,
		Name:        en.Name,
		Description: en.Description,
		Enabled:     en.Enabled,
	}
}

func ConvertMissionGroup(mg *utask.MissionGroup) *MissionGroup {
	_tasks := []*Mission{}
	for _, _m := range mg.Tasks {
		_tasks = append(_tasks, &Mission{
			Task:   UnConvertTask(_m.Task),
			Limit:  int32(_m.Limit),
			Count:  int32(_m.Count),
			Status: UnConvertMISSION_STATUS(_m.Status),
		})
	}
	return &MissionGroup{
		Group: UnConvertTaskGroup(mg.Group),
		Tasks: _tasks,
		Dones: int32(mg.Dones),
	}
}

func (t UserTaskServiceProxy) CreateGroup(group *TaskGroup) (r *ActionResult_, err error) {
	_group := ConvertTaskGroup(group)
	_err := t.tasker.CreateGroup(_group)
	if _err == nil {
		return &ActionResult_{
			Id:        _group.GroupId,
			Success:   true,
			Exception: "",
		}, nil
	}
	return &ActionResult_{
		Id:        0,
		Success:   false,
		Exception: _err.Error(),
	}, _err
}

func (t UserTaskServiceProxy) UpdateGroup(group *TaskGroup) (r *ActionResult_, err error) {
	_group := ConvertTaskGroup(group)
	_err := t.tasker.UpdateGroup(_group)
	if _err == nil {
		return &ActionResult_{
			Id:        _group.GroupId,
			Success:   true,
			Exception: "",
		}, nil
	}
	return &ActionResult_{
		Id:        0,
		Success:   false,
		Exception: _err.Error(),
	}, _err
}

func (t UserTaskServiceProxy) DeleteGroup(groupId int64) (r *ActionResult_, err error) {
	_err := t.tasker.DeleteGroup(groupId)
	if _err == nil {
		return &ActionResult_{
			Id:        groupId,
			Success:   true,
			Exception: "",
		}, nil
	}
	return &ActionResult_{
		Id:        0,
		Success:   false,
		Exception: _err.Error(),
	}, _err
}

func (t UserTaskServiceProxy) GetGroup(groupId int64) (r *TaskGroup, err error) {
	_group := t.tasker.GetGroup(groupId)
	if _group == nil {
		return nil, fmt.Errorf("分组不存在")
	}
	return UnConvertTaskGroup(_group), nil
}

func (t UserTaskServiceProxy) GetGroups() (r []*TaskGroup, err error) {
	_groups := t.tasker.GetGroups()
	outgs := []*TaskGroup{}
	for _, group := range _groups {
		outgs = append(outgs, UnConvertTaskGroup(group))
	}
	return outgs, nil
}

func (t UserTaskServiceProxy) CreateTask(task *Task) (r *ActionResult_, err error) {
	_task := ConvertTask(task)
	_err := t.tasker.CreateTask(_task)
	if _err == nil {
		return &ActionResult_{
			Id:        _task.TaskId,
			Success:   true,
			Exception: "",
		}, nil
	}
	return &ActionResult_{
		Id:        0,
		Success:   false,
		Exception: _err.Error(),
	}, _err
}

func (t UserTaskServiceProxy) UpdateTask(task *Task) (r *ActionResult_, err error) {
	_task := ConvertTask(task)
	_err := t.tasker.UpdateTask(_task)
	if _err == nil {
		return &ActionResult_{
			Id:        _task.TaskId,
			Success:   true,
			Exception: "",
		}, nil
	}
	return &ActionResult_{
		Id:        0,
		Success:   false,
		Exception: _err.Error(),
	}, _err
}

func (t UserTaskServiceProxy) DeleteTask(taskId int64, groupId int64) (r *ActionResult_, err error) {
	_task := t.tasker.GetTask(taskId)
	var _err error
	if _task != nil {
		_err = t.tasker.DeleteTask(taskId, _task.GroupId)
		if _err == nil {
			return &ActionResult_{
				Id:        taskId,
				Success:   true,
				Exception: "",
			}, nil
		}
	} else {
		_err = fmt.Errorf("任务不存在")
	}
	return &ActionResult_{
		Id:        0,
		Success:   false,
		Exception: _err.Error(),
	}, _err
}

func (t UserTaskServiceProxy) GetTask(taskId int64) (r *Task, err error) {
	_task := t.tasker.GetTask(taskId)
	if _task == nil {
		return nil, fmt.Errorf("任务不存在")
	}
	return UnConvertTask(_task), nil
}

func (t UserTaskServiceProxy) GetTasks(groupId int64) (r []*Task, err error) {
	_tasks := t.tasker.GetGroupTasks(groupId)
	outks := []*Task{}
	for _, _task := range _tasks {
		outks = append(outks, UnConvertTask(_task))
	}
	return outks, nil
}

func (t UserTaskServiceProxy) GetEventNames() (r []*TaskEventName, err error) {
	_ens := t.tasker.GetEvents()
	outs := []*TaskEventName{}
	for _, _en := range _ens {
		if _en.Enabled {
			outs = append(outs, UnConvertEventName(_en))
		}
	}
	return outs, nil
}

func (t UserTaskServiceProxy) GetTaskTimers() (r []*TaskTimer, err error) {
	outtts := []*TaskTimer{}
	tts := t.tasker.GetTaskTimers()
	for _, tt := range tts {
		outtts = append(outtts, UnConvertTaskTimer(tt))
	}
	return outtts, nil
}

func (t UserTaskServiceProxy) IsMissionPassed(uid int64, taskId int64) (r bool, err error) {
	passed := t.enigne.IsMissionPassed(uid, taskId)
	return passed, nil
}

func (t UserTaskServiceProxy) GetMissions(uid int64) (r []*MissionGroup, err error) {
	mgroups := t.enigne.GetMissions(uid)
	outs := []*MissionGroup{}
	for _, mg := range mgroups {
		outs = append(outs, ConvertMissionGroup(mg))
	}
	return outs, nil
}

func (t UserTaskServiceProxy) EventHandler(uid int64, event string, n int32) (r *EventResult_, err error) {
	result, err := t.enigne.MissionEventHandler(uid, event, n)
	return UnConvertEventResult(result, err), err
}
