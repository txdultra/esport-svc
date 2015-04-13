package utask

import (
	"fmt"
	"time"
	"utils"
)

type TaskTimer struct {
	TaskId       int64
	Timer        *time.Timer
	RunTime      int64
	PreResetTime int64
}

var taskTimers = make(map[int64]*TaskTimer)

//启动任务重置定时器
func initTaskTimers() {
	for _, task := range allTasks {
		if task.TaskType == TASK_TYPE_DUPLICATE {
			_taskId := task.TaskId
			nextResetTask(_taskId)
		}
	}
}

func getTaskTimers() []*TaskTimer {
	tts := []*TaskTimer{}
	for _, tt := range taskTimers {
		tts = append(tts, tt)
	}
	return tts
}

func nextResetTask(taskId int64) {
	tasker := &Tasker{}
	task := tasker.GetTask(taskId)
	d := nextResetTaskDuration(task)
	_taskId := taskId
	timer := time.AfterFunc(d, func() {
		_tasker := &Tasker{}
		_tasker.ResetTaskVar(_taskId, 1, time.Now().Unix())
		nextResetTask(_taskId)
	})
	taskTimers[_taskId] = &TaskTimer{
		TaskId:       _taskId,
		Timer:        timer,
		RunTime:      time.Now().Add(d).Unix(),
		PreResetTime: task.LastResetTime,
	}
}

func getDate(t time.Time) (time.Time, error) {
	y, m, d := t.Date()
	str := fmt.Sprintf("%d-%.2d-%.2d 00:00:00", y, m, d)
	return utils.StrToTime(str)
}

func nextResetTaskDuration(task *Task) time.Duration {
	var periodDuration time.Duration
	switch task.PeriodType {
	case TASK_PERIOD_TYPE_DAY:
		periodDuration = time.Duration(24) * time.Duration(task.Period) * time.Hour
	default:
		panic("日常任务的周期类型错误")
	}
	periodDuration += time.Duration(task.ResetTime) * time.Second
	var runTime time.Time
	if task.LastResetTime <= 0 {
		_t, _ := getDate(time.Unix(task.StartTime, 0))
		runTime = _t.Add(periodDuration)
	} else {
		_t, _ := getDate(time.Unix(task.LastResetTime, 0))
		runTime = _t.Add(periodDuration)
	}
	if runTime.Before(time.Now()) {
		runTime = time.Now().Add(2 * time.Second)
	}
	//测试
	//runTime = time.Now().Add(300 * time.Second)

	return runTime.Sub(time.Now())
}

func resetTaskTimer(task *Task) {
	if task.TaskType != TASK_TYPE_DUPLICATE {
		return
	}
	dur := nextResetTaskDuration(task)
	if tt, ok := taskTimers[task.TaskId]; ok {
		tt.Timer.Reset(dur)
		tt.RunTime = time.Now().Add(dur).Unix()
	}
}

func deleteTaskTimer(taskId int64) {
	if tt, ok := taskTimers[taskId]; ok {
		tt.Timer.Stop()
		delete(taskTimers, taskId)
	}
}

func addTaskTimer(task *Task) {
	if task.TaskType == TASK_TYPE_DUPLICATE {
		nextResetTask(task.TaskId)
	}
}
