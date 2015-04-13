package utask

import (
	"dbs"
	"fmt"
	"sync"
	"utils"

	"github.com/astaxie/beego/orm"
	"github.com/pmylund/sortutil"
)

var allTaskGroups = make(map[int64]*TaskGroup)
var allTaskGroupsOrderByDesc []*TaskGroup
var allTasks = make(map[int64]*Task)
var allTaskGroupTasks = make(map[int64][]*Task)
var allEventTasks = make(map[string][]*Task)
var allEvents = make(map[string]*TaskEvent)
var lock *sync.RWMutex = new(sync.RWMutex)
var once sync.Once

func taskInit() {
	once.Do(func() {
		fmt.Println("utask init...")

		initTaskGroups()
		initTasks()
		initEvents()
		initTaskTimers()

		fmt.Println("utask init completed")
	})
}

//载入所有分组
func initTaskGroups() {
	allTaskGroupsOrderByDesc = []*TaskGroup{}
	o := dbs.NewOrm(db_aliasname)
	groups := []*TaskGroup{}
	o.QueryTable(&TaskGroup{}).OrderBy("-displayorder").All(&groups)
	for _, group := range groups {
		allTaskGroupsOrderByDesc = append(allTaskGroupsOrderByDesc, group)
		allTaskGroups[group.GroupId] = group
	}
}

//载入所有任务
func initTasks() {
	o := dbs.NewOrm(db_aliasname)
	tasks := []*Task{}
	o.QueryTable(&Task{}).OrderBy("-displayorder").All(&tasks)
	for _, task := range tasks {
		allTasks[task.TaskId] = task
		//gourp -> tasks mapping
		if gts, ok := allTaskGroupTasks[task.GroupId]; ok {
			gts = append(gts, task)
			allTaskGroupTasks[task.GroupId] = gts
		} else {
			gts := []*Task{}
			gts = append(gts, task)
			allTaskGroupTasks[task.GroupId] = gts
		}
		//scriptname -> tasks mapping
		if gts, ok := allEventTasks[task.EventName]; ok {
			gts = append(gts, task)
			allEventTasks[task.EventName] = gts
		} else {
			gts := []*Task{}
			gts = append(gts, task)
			allEventTasks[task.EventName] = gts
		}
	}
}

//载入事件列表
func initEvents() {
	o := dbs.NewOrm(db_aliasname)
	events := []*TaskEvent{}
	o.QueryTable(&TaskEvent{}).All(&events)
	for _, en := range events {
		if _, ok := allEvents[en.Name]; ok {
			continue
		}
		allEvents[en.Name] = en
	}
}

func NewTasker() *Tasker {
	taskInit()
	return &Tasker{}
}

type Tasker struct{}

func (t *Tasker) setTaskGroupMapping(group *TaskGroup) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := allTaskGroups[group.GroupId]; ok {
		allTaskGroups[group.GroupId] = group
		for i := range allTaskGroupsOrderByDesc {
			if allTaskGroupsOrderByDesc[i].GroupId == group.GroupId {
				allTaskGroupsOrderByDesc[i] = group
			}
		}
	} else {
		allTaskGroups[group.GroupId] = group
		allTaskGroupsOrderByDesc = append(allTaskGroupsOrderByDesc, group)
	}
	sortutil.DescByField(allTaskGroupsOrderByDesc, "DisplayOrder")
}

func (t *Tasker) getTaskGroupMapping(groupId int64) *TaskGroup {
	lock.RLock()
	defer lock.RUnlock()
	if g, ok := allTaskGroups[groupId]; ok {
		return g
	}
	return nil
}

func (t *Tasker) delTaskGroupMapping(groupId int64) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := allTaskGroups[groupId]; ok {
		delete(allTaskGroups, groupId)
	}
}

func (t *Tasker) checkGroup(group *TaskGroup) error {
	if group == nil {
		return fmt.Errorf("参数错误")
	}
	if group.GroupId > 0 {
		_org := t.GetGroup(group.GroupId)
		if _org != nil && _org.TaskType != group.TaskType {
			return fmt.Errorf("不能更新分组任务类型")
		}
	}
	if len(string(group.TaskType)) == 0 {
		return fmt.Errorf("TaskType属性未设置")
	}
	if len(group.Name) == 0 {
		return fmt.Errorf("名称未设置")
	}
	return nil
}

func (t *Tasker) CreateGroup(group *TaskGroup) error {
	err := t.checkGroup(group)
	if err != nil {
		return err
	}

	group.Tasks = 0
	o := dbs.NewOrm(db_aliasname)
	id, err := o.Insert(group)
	if err != nil {
		return err
	}
	group.GroupId = id
	//	cache := utils.GetCache()
	//	cache.Add(t.groupCacheKey(id), *group, 0)
	//	t.clearGroupsCache()
	t.setTaskGroupMapping(group)
	return nil
}

func (t *Tasker) UpdateGroup(group *TaskGroup) error {
	err := t.checkGroup(group)
	if err != nil {
		return err
	}
	if group.GroupId <= 0 {
		return fmt.Errorf("groupid未设置")
	}

	o := dbs.NewOrm(db_aliasname)
	_, err = o.Update(group, "grouptype", "tasktype", "name", "description", "icon", "bgimg", "enabled", "starttime", "endtime", "displayorder",
		"ex1", "ex2", "ex3", "ex4", "ex5")
	if err != nil {
		return err
	}
	//	cache := utils.GetCache()
	//	cache.Set(t.groupCacheKey(group.GroupId), group, 0)
	//	t.clearGroupsCache()
	t.setTaskGroupMapping(group)
	return nil
}

func (t *Tasker) DeleteGroup(groupId int64) error {
	group := t.GetGroup(groupId)
	if group == nil {
		return fmt.Errorf("不存在分组")
	}
	if group.Tasks > 0 {
		return fmt.Errorf("必须删除子任务后才能删除分组")
	}
	o := dbs.NewOrm(db_aliasname)
	num, err := o.QueryTable(&TaskGroup{}).Filter("groupid", groupId).Delete()
	if err != nil {
		return err
	}
	if num <= 0 {
		return fmt.Errorf("不存在分组")
	}
	t.delTaskGroupMapping(groupId)
	return nil
}

func (t *Tasker) GetGroup(groupId int64) *TaskGroup {
	group := t.getTaskGroupMapping(groupId)
	if group != nil {
		return group
	}
	return nil
}

func (t *Tasker) GetGroups() []*TaskGroup {
	return allTaskGroupsOrderByDesc
}

////////////////////////////////////////////////////////////////////////////////
//task
////////////////////////////////////////////////////////////////////////////////

func (t *Tasker) setTaskMapping(task *Task) {
	lock.Lock()
	defer lock.Unlock()
	var originalGroupId int64 = 0
	var originalEventName string
	if t, ok := allTasks[task.TaskId]; ok {
		originalGroupId = t.GroupId
		originalEventName = t.EventName
	}
	allTasks[task.TaskId] = task
	//删除原来的任务分组
	if originalGroupId > 0 {
		if gts, ok := allTaskGroupTasks[originalGroupId]; ok {
			newTasks := []*Task{}
			for _, gt := range gts {
				if gt.TaskId != task.TaskId {
					newTasks = append(newTasks, gt)
				}
			}
			allTaskGroupTasks[originalGroupId] = newTasks
		}
	}
	if len(originalEventName) > 0 {
		if gts, ok := allEventTasks[originalEventName]; ok {
			newTasks := []*Task{}
			for _, gt := range gts {
				if gt.TaskId != task.TaskId {
					newTasks = append(newTasks, gt)
				}
			}
			allEventTasks[originalEventName] = newTasks
		}
	}
	//重新加入到任务分组中
	if gts, ok := allTaskGroupTasks[task.GroupId]; ok {
		gts = append(gts, task)
		sortutil.DescByField(gts, "DisplayOrder")
		allTaskGroupTasks[task.GroupId] = gts
	} else {
		gts = []*Task{}
		gts = append(gts, task)
		allTaskGroupTasks[task.GroupId] = gts
	}
	//重新加入到时间任务分组中
	if gts, ok := allEventTasks[task.EventName]; ok {
		gts = append(gts, task)
		sortutil.DescByField(gts, "DisplayOrder")
		allEventTasks[task.EventName] = gts
	} else {
		gts = []*Task{}
		gts = append(gts, task)
		allEventTasks[task.EventName] = gts
	}
}

func (t *Tasker) getTaskMapping(taskId int64) *Task {
	lock.RLock()
	defer lock.RUnlock()
	if t, ok := allTasks[taskId]; ok {
		return t
	}
	return nil
}

func (t *Tasker) delTaskMapping(taskId int64) {
	lock.Lock()
	defer lock.Unlock()
	if t, ok := allTasks[taskId]; ok {
		if gts, ok2 := allTaskGroupTasks[t.GroupId]; ok2 {
			newTasks := []*Task{}
			for _, gt := range gts {
				if gt.TaskId != taskId {
					newTasks = append(newTasks, gt)
				}
			}
			allTaskGroupTasks[t.GroupId] = newTasks
		}
		if gts, ok3 := allEventTasks[t.EventName]; ok3 {
			newTasks := []*Task{}
			for _, gt := range gts {
				if gt.TaskId != taskId {
					newTasks = append(newTasks, gt)
				}
			}
			allEventTasks[t.EventName] = newTasks
		}
		delete(allTasks, taskId)
	}
}

func (t *Tasker) checkTask(task *Task) error {
	if task == nil {
		return fmt.Errorf("参数错误")
	}
	if len(string(task.TaskType)) == 0 {
		return fmt.Errorf("TaskType属性未设置")
	}
	if len(task.Name) == 0 {
		return fmt.Errorf("名称未设置")
	}
	if task.GroupId <= 0 {
		return fmt.Errorf("groupid未设置")
	}
	group := t.GetGroup(task.GroupId)
	if group == nil {
		return fmt.Errorf("分组不存在")
	}
	if group.GroupType == TASK_GROUP_TYPE_WEB && task.TaskId == 0 && group.Tasks > 0 {
		return fmt.Errorf("分组类型是页面类型只能有一个子任务")
	}
	if group.GroupType == TASK_GROUP_TYPE_WEB && !utils.IsUrl(task.Ex1) {
		return fmt.Errorf("分组类型为页面类型必须填写跳转URL")
	}
	if group.TaskType != task.TaskType {
		return fmt.Errorf("任务类型和分组不一致")
	}
	if len(task.EventName) == 0 {
		return fmt.Errorf("eventname未设置")
	}
	if t.GetEvent(task.EventName) == nil {
		return fmt.Errorf("事件名称不存在")
	}
	if len(string(task.PeriodType)) == 0 {
		return fmt.Errorf("周期类型未设置")
	}
	if len(string(task.Reward)) == 0 {
		return fmt.Errorf("奖励类型未设置")
	}
	return nil
}

func (t *Tasker) CreateTask(task *Task) error {
	err := t.checkTask(task)
	if err != nil {
		return err
	}

	task.ResetVar = 0
	task.LastResetTime = 0
	if task.TaskType == TASK_TYPE_ONCE {
		task.PeriodType = TASK_PERIOD_TYPE_NULL
		task.Period = 0
	}

	o := dbs.NewOrm(db_aliasname)
	o.Begin()
	id, err := o.Insert(task)
	if err != nil {
		o.Rollback()
		return err
	}
	task.TaskId = id
	_, err = o.QueryTable(&TaskGroup{}).Filter("groupid", task.GroupId).Update(orm.Params{
		"tasks": orm.ColValue(orm.Col_Add, 1),
	})
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	//更新缓冲
	group := t.getTaskGroupMapping(task.GroupId)
	if group != nil {
		group.Tasks++ //指针直接累加
	}
	t.setTaskMapping(task)
	addTaskTimer(task)
	return nil
}

func (t *Tasker) UpdateTask(task *Task) error {
	err := t.checkTask(task)
	if err != nil {
		return err
	}
	if task.TaskId <= 0 {
		return fmt.Errorf("taskid未设置")
	}

	o := dbs.NewOrm(db_aliasname)
	_, err = o.Update(task, "tasktype", "groupid", "name", "description", "icon", "tasklimits", "applyperm", "eventname", "starttime",
		"endtime", "resettime", "period", "periodtype", "reward", "prize", "applicants", "achievers", "version", "displayorder",
		"ex1", "ex2", "ex3", "ex4", "ex5")
	if err != nil {
		return err
	}
	t.setTaskMapping(task)
	resetTaskTimer(task)
	return nil
}

func (t *Tasker) ResetTaskVar(taskId int64, incrN int, ts int64) error {
	if taskId <= 0 {
		return fmt.Errorf("taskid错误")
	}
	task := t.GetTask(taskId)
	if task == nil {
		return fmt.Errorf("任务不存在")
	}
	o := dbs.NewOrm(db_aliasname)
	_, err := o.QueryTable(&Task{}).Filter("taskid", taskId).Update(orm.Params{
		"resetvar":      orm.ColValue(orm.Col_Add, incrN),
		"lastresettime": ts,
	})
	if err == nil {
		task.ResetVar += incrN //指针直接累加
		task.LastResetTime = ts
	}
	return err
}

func (t *Tasker) DeleteTask(taskId int64, groupId int64) error {
	o := dbs.NewOrm(db_aliasname)
	o.Begin()
	_, err := o.QueryTable(&Task{}).Filter("taskid", taskId).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryTable(&TaskGroup{}).Filter("groupid", groupId).Update(orm.Params{
		"tasks": orm.ColValue(orm.Col_Minus, 1),
	})
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	//更新缓冲
	group := t.getTaskGroupMapping(groupId)
	if group != nil {
		group.Tasks-- //指针直接累加
	}
	t.delTaskMapping(taskId)
	deleteTaskTimer(taskId)
	return nil
}

func (t *Tasker) GetTask(taskId int64) *Task {
	if task, ok := allTasks[taskId]; ok {
		return task
	}
	return nil
}

func (t *Tasker) GetGroupTasks(groupId int64) []*Task {
	if ts, ok := allTaskGroupTasks[groupId]; ok {
		return ts
	}
	return nil
}

func (t *Tasker) GetEventTasks(eventName string) []*Task {
	if ts, ok := allEventTasks[eventName]; ok {
		return ts
	}
	return nil
}

func (t *Tasker) GetEvents() []*TaskEvent {
	o := dbs.NewOrm(db_aliasname)
	var tsns []*TaskEvent
	o.QueryTable(&TaskEvent{}).All(&tsns)
	return tsns
}

func (t *Tasker) GetEvent(name string) *TaskEvent {
	if e, ok := allEvents[name]; ok {
		return e
	}
	return nil
}

func (t *Tasker) GetTaskTimers() []*TaskTimer {
	return getTaskTimers()
}
