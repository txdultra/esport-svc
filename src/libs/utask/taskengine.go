package utask

import (
	"dbs"
	"fmt"
	"libs/credits/proxy"
	msgclient "libs/message/client"
	"libs/vars"
	"sync"
	"time"
	"utils/ssdb"

	"github.com/thrift"
)

type MissionRewardParameter struct {
	Uid    int64
	Points int64
	UMId   string //任务记录单号
	Info   string
}

type IMissionReward interface {
	Reward(mrp *MissionRewardParameter) (string, error)
}

type CreditMissionReward struct{}

func (c CreditMissionReward) Reward(mrp *MissionRewardParameter) (string, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTSocket(credit_service_host)
	if err != nil {
		return "", fmt.Errorf("error resolving address:%s", err.Error())
	}

	useTransport := transportFactory.GetTransport(transport)
	client := proxy.NewCreditServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return "", fmt.Errorf("Error opening socket to %s;error:%s", credit_service_host, err.Error())
	}
	defer transport.Close()

	param := &proxy.OperationCreditParameter{
		Uid:    mrp.Uid,
		Points: mrp.Points,
		Desc:   fmt.Sprintf("完成任务(%s)获得积分", mrp.Info),
		Action: proxy.OPERATION_ACTOIN_INCR,
		Ref:    "utask",
		RefId:  mrp.UMId,
	}

	result, err := client.Do(param)
	if result.State == proxy.OPERATION_STATE_SUCCESS {
		//发送系统消息
		go func() {
			msgc, msgt, err := msgclient.NewClient(utask_message_send_host)
			if err != nil {
				return
			}
			defer msgt.Close()
			msgTxt := fmt.Sprintf("您完成%s任务获得了%d积分", mrp.Info, mrp.Points)
			err = msgc.Send(vars.MESSAGE_SYS_ID, mrp.Uid, string(vars.MSG_TYPE_SYS), msgTxt, result.No)
			if err != nil {
				return
			}
		}()
		return result.No, nil
	}
	return "", fmt.Errorf(result.Err)
}

type IMissionEngine interface {
	IsMissionPassed(uid int64, taskId int64) bool
	GetMissions(uid int64) []*MissionGroup
	MissionEventHandler(uid int64, eventName string, obj interface{}) (*EventResult, error)
}

type MISSION_STATUS int

const (
	MISSION_STATUS_FAILURE MISSION_STATUS = 0
	MISSION_STATUS_UNDONE  MISSION_STATUS = 1
	MISSION_STATUS_DONE    MISSION_STATUS = 2
)

type MissionGroup struct {
	Group *TaskGroup
	Tasks []*Mission
	Dones int
}

type Mission struct {
	Task   *Task
	Limit  int //需要完成的计数
	Count  int //已完成的计数
	Status MISSION_STATUS
}

type MemberMissionStatus struct {
	Var    int
	Limit  int //需要完成的计数
	Count  int //已完成的计数
	Status MISSION_STATUS
}

type EventResult struct {
	DoneTaskIds      []int64
	RelevanceTaskIds []int64
	Rewards          []*MissionReward
}

type MissionReward struct {
	RewardType TASK_REWARD_TYPE
	Prize      int64
}

func NewMissionEngine(reward IMissionReward) *MissionEngine {
	engine := &MissionEngine{}
	engine.tasksys = NewTasker()
	engine.creditReward = reward
	return engine
}

//用户锁
var m_lockers map[int64]*sync.Mutex = make(map[int64]*sync.Mutex)
var _tlock *sync.Mutex = new(sync.Mutex)

func mLock(uid int64) {
	if l, ok := m_lockers[uid]; ok {
		l.Lock()
		return
	}
	_tlock.Lock()
	defer _tlock.Unlock()
	if l, ok := m_lockers[uid]; ok {
		l.Lock()
		return
	}
	_l := new(sync.Mutex)
	m_lockers[uid] = _l
	_l.Lock()
}

func mUnLock(uid int64) {
	if l, ok := m_lockers[uid]; ok {
		l.Unlock()
	}
}

type MissionEngine struct {
	tasksys      *Tasker
	creditReward IMissionReward
}

func (t MissionEngine) IsMissionPassed(uid int64, taskId int64) bool {
	memberMissionStatus := t.getMemberMissionStatus(uid)
	task := t.tasksys.GetTask(taskId)
	if task == nil {
		return false
	}
	mstate := t.memberMissionStatus(memberMissionStatus, task)
	if mstate == MISSION_STATUS_DONE {
		return true
	}
	return false
}

func (t *MissionEngine) member_task_state_key(uid int64) string {
	return fmt.Sprintf("utask_member_task_state_maps_uid:%d", uid)
}

func (t *MissionEngine) getMemberMissionStatus(uid int64) map[int64]MemberMissionStatus {
	memberMissionStatus := make(map[int64]MemberMissionStatus)
	//获取用户任务状态列表
	_key := t.member_task_state_key(uid)
	ssdb.New(use_ssdb_utask_db).Get(_key, &memberMissionStatus)
	return memberMissionStatus
}

func (t *MissionEngine) setMemberMissioinStatus(uid int64, state map[int64]MemberMissionStatus) {
	_key := t.member_task_state_key(uid)
	ssdb.New(use_ssdb_utask_db).Set(_key, state)
}

func (t *MissionEngine) checkTaskTime(task *Task, ts int64) bool {
	if task.StartTime > ts || task.EndTime < ts {
		return false
	}
	return true
}

func (t MissionEngine) GetMissions(uid int64) []*MissionGroup {
	mgs := []*MissionGroup{}
	allgroups := t.tasksys.GetGroups()
	ts := time.Now().Unix()
	memberMissionStatus := t.getMemberMissionStatus(uid)
	for _, group := range allgroups {
		if !group.Enabled || group.StartTime > ts || group.EndTime < ts {
			continue
		}

		groupTasks := t.tasksys.GetGroupTasks(group.GroupId)
		_gtask := []*Mission{}
		dones := 0

		//group_type => web类型只有一个子任务
		if group.GroupType == TASK_GROUP_TYPE_WEB && len(groupTasks) == 1 {
			task := groupTasks[0]
			if !t.checkTaskTime(task, ts) {
				continue
			}
			_c := 0
			_status := MISSION_STATUS_UNDONE
			_limit := 0
			if state, ok := memberMissionStatus[task.TaskId]; ok {
				_c = state.Count
				_status = state.Status
				_limit = state.Limit
			} else {
				_limit = task.TaskLimits
			}
			if _status == MISSION_STATUS_DONE {
				dones++
			}
			_gtask = append(_gtask, &Mission{
				Task:   task,
				Limit:  _limit,
				Count:  _c,
				Status: _status,
			})
			mgs = append(mgs, &MissionGroup{
				Group: group,
				Tasks: _gtask,
				Dones: dones,
			})
			continue
		}

		//group_type => list类型装载子任务
		if group.GroupType == TASK_GROUP_TYPE_LIST {
			//一次性任务
			if group.TaskType == TASK_TYPE_ONCE {
				for _, task := range groupTasks {
					if !t.checkTaskTime(task, ts) {
						continue
					}
					_c := 0
					_status := MISSION_STATUS_UNDONE
					_limit := 0
					if state, ok := memberMissionStatus[task.TaskId]; ok {
						_c = state.Count
						_status = state.Status
						_limit = state.Limit
					} else {
						_limit = task.TaskLimits
					}
					if _status == MISSION_STATUS_DONE {
						dones++
					}
					_gtask = append(_gtask, &Mission{
						Task:   task,
						Limit:  _limit,
						Count:  _c,
						Status: _status,
					})
				}
				if len(_gtask) > 0 {
					mgs = append(mgs, &MissionGroup{
						Group: group,
						Tasks: _gtask,
						Dones: dones,
					})
				}
			}
			//日常任务
			if group.TaskType == TASK_TYPE_DUPLICATE {
				for _, task := range groupTasks {
					if !t.checkTaskTime(task, ts) {
						continue
					}
					_c := 0
					_status := MISSION_STATUS_UNDONE
					_limit := 0
					if state, ok := memberMissionStatus[task.TaskId]; ok {
						if state.Var == task.ResetVar {
							_c = state.Count
							_status = state.Status
							_limit = state.Limit
						} else {
							_limit = task.TaskLimits
						}
					} else {
						_limit = task.TaskLimits
					}
					if _status == MISSION_STATUS_DONE {
						dones++
					}
					_gtask = append(_gtask, &Mission{
						Task:   task,
						Limit:  _limit,
						Count:  _c,
						Status: _status,
					})
				}
				if len(_gtask) > 0 {
					mgs = append(mgs, &MissionGroup{
						Group: group,
						Tasks: _gtask,
						Dones: dones,
					})
				}
			}
		}
	}
	return mgs
}

func (t *MissionEngine) memberMissionStatus(memberState map[int64]MemberMissionStatus, task *Task) MISSION_STATUS {
	ts := time.Now().Unix()
	if task.StartTime > ts || task.EndTime < ts {
		return MISSION_STATUS_FAILURE
	}
	if state, ok := memberState[task.TaskId]; ok {
		if task.TaskType == TASK_TYPE_DUPLICATE {
			if state.Var == task.ResetVar && state.Status == MISSION_STATUS_DONE {
				return MISSION_STATUS_DONE
			}
			return MISSION_STATUS_UNDONE
		}
		if task.TaskType == TASK_TYPE_ONCE {
			if state.Status == MISSION_STATUS_DONE {
				return MISSION_STATUS_DONE
			}
			return MISSION_STATUS_UNDONE
		}
	}
	return MISSION_STATUS_UNDONE
}

func (t *MissionEngine) calculateMissionStatus(task *Task, mms *MemberMissionStatus, event *TaskEvent, n int) (MISSION_STATUS, *MissionReward) {
	if event.EventType == EVENT_TYPE_COUNT { //计数型
		mms.Count += n
		if mms.Limit <= mms.Count {
			mms.Status = MISSION_STATUS_DONE
			return MISSION_STATUS_DONE, &MissionReward{
				RewardType: task.Reward,
				Prize:      task.Prize,
			}
		}
	}
	if event.EventType == EVENT_TYPE_NUM { //总数型
		if task.TaskLimits <= n {
			mms.Count = task.TaskLimits
			mms.Status = MISSION_STATUS_DONE
			return MISSION_STATUS_DONE, &MissionReward{
				RewardType: task.Reward,
				Prize:      task.Prize,
			}
		}
	}
	return MISSION_STATUS_UNDONE, nil
}

func (t *MissionEngine) missionDoneToDb(uid int64, task Task, tt time.Time) string {
	tbl := hash_task_records_tbl(uid)
	o := dbs.NewOrm(db_aliasname)
	sql := fmt.Sprintf("insert into %s(no,uid,taskid,ts) values(?,?,?,?)", tbl)
	no := BuildTaskRecordNo(uid, &task, tt)
	o.Raw(sql, no, uid, task.TaskId, tt.Unix()).Exec()
	return no
}

func (t *MissionEngine) missionDoneReward(uid int64, reward *MissionReward, recordNo string, info string) error {
	if reward == nil {
		return fmt.Errorf("任务奖励对象不能为nil")
	}
	//接入积分系统api
	mrp := &MissionRewardParameter{
		Uid:    uid,
		Points: reward.Prize,
		UMId:   recordNo,
		Info:   info,
	}
	rewardNo, err := t.creditReward.Reward(mrp)
	if err != nil {
		return err
	}

	//发送系统消息
	//	go func() {
	//		message.
	//	}()

	tbl := hash_task_records_tbl(uid)
	o := dbs.NewOrm(db_aliasname)
	//现在只支持credit币种=> rewardno1
	sql := fmt.Sprintf("update %s set rewardno1=? where no=?", tbl)
	o.Raw(sql, rewardNo, recordNo).Exec()
	return nil
}

func (t MissionEngine) MissionEventHandler(uid int64, eventName string, obj interface{}) (*EventResult, error) {
	if len(eventName) == 0 {
		return nil, fmt.Errorf("eventName 参数不能为空")
	}
	event := t.tasksys.GetEvent(eventName)
	if event == nil || !event.Enabled {
		return nil, fmt.Errorf("不存在处理事件")
	}
	tasks := t.tasksys.GetEventTasks(eventName)
	if len(tasks) == 0 {
		return nil, fmt.Errorf("没有事件对应的任务")
	}

	//ssdb.New(use_ssdb_utask_db).Del(t.member_task_state_key(uid)) //测试时重置

	mLock(uid)
	defer mUnLock(uid)

	doneTaskIds := []int64{}
	relevanceTaskIds := []int64{}
	_n, ok := obj.(int32)
	if !ok {
		return nil, fmt.Errorf("obj参数错误")
	}
	n := int(_n)

	mms := t.getMemberMissionStatus(uid)
	rewards := []*MissionReward{}
	now := time.Now()
	for _, task := range tasks {
		mstate := t.memberMissionStatus(mms, task)
		if mstate == MISSION_STATUS_DONE || mstate == MISSION_STATUS_FAILURE {
			continue
		}
		//日常任务
		status := MISSION_STATUS_UNDONE
		var _rwd *MissionReward = nil

		if task.TaskType == TASK_TYPE_DUPLICATE {
			if state, ok := mms[task.TaskId]; ok {
				if state.Var == task.ResetVar {
					status, _rwd = t.calculateMissionStatus(task, &state, event, n)
				} else { //重置版本不一致
					state = MemberMissionStatus{
						Var:    task.ResetVar,
						Limit:  task.TaskLimits,
						Count:  0,
						Status: MISSION_STATUS_UNDONE,
					}
					status, _rwd = t.calculateMissionStatus(task, &state, event, n)
				}
				mms[task.TaskId] = state
			} else {
				state = MemberMissionStatus{
					Var:    task.ResetVar,
					Limit:  task.TaskLimits,
					Count:  0,
					Status: MISSION_STATUS_UNDONE,
				}
				status, _rwd = t.calculateMissionStatus(task, &state, event, n)
				mms[task.TaskId] = state
			}
			t.setMemberMissioinStatus(uid, mms)
		}
		//一次性任务
		if task.TaskType == TASK_TYPE_ONCE {
			if state, ok := mms[task.TaskId]; ok {
				status, _rwd = t.calculateMissionStatus(task, &state, event, n)
				mms[task.TaskId] = state
			} else {
				state = MemberMissionStatus{
					Var:    task.ResetVar,
					Limit:  task.TaskLimits,
					Count:  0,
					Status: MISSION_STATUS_UNDONE,
				}
				status, _rwd = t.calculateMissionStatus(task, &state, event, n)
				mms[task.TaskId] = state
			}
			t.setMemberMissioinStatus(uid, mms)
		}

		if _rwd != nil {
			rewards = append(rewards, _rwd)
		}
		if status == MISSION_STATUS_DONE {
			doneTaskIds = append(doneTaskIds, task.TaskId)
			recordNo := t.missionDoneToDb(uid, *task, now)
			_task := task
			go func() {
				t.missionDoneReward(uid, _rwd, recordNo, _task.Name)
			}()
		}
		relevanceTaskIds = append(relevanceTaskIds, task.TaskId)
	}
	return &EventResult{
		DoneTaskIds:      doneTaskIds,
		RelevanceTaskIds: relevanceTaskIds,
		Rewards:          rewards,
	}, nil
}
