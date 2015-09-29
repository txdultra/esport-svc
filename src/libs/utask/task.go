package utask

import (
	"dbs"
	"fmt"
	"sync"
	"time"
)

var taskRecordTbls = make(map[string]bool)
var n_locker *sync.Mutex = new(sync.Mutex)

type TASK_GROUP_TYPE string

const (
	TASK_GROUP_TYPE_LIST TASK_GROUP_TYPE = "list"
	TASK_GROUP_TYPE_WEB  TASK_GROUP_TYPE = "web"
)

type TASK_TYPE string

const (
	TASK_TYPE_ONCE      TASK_TYPE = "once"
	TASK_TYPE_DUPLICATE TASK_TYPE = "duplicate"
)

type TASK_PERIOD_TYPE string

const (
	TASK_PERIOD_TYPE_NULL TASK_PERIOD_TYPE = "null"
	TASK_PERIOD_TYPE_DAY  TASK_PERIOD_TYPE = "day"
)

type TASK_REWARD_TYPE string

const (
	TASK_REWARD_TYPE_CREDIT TASK_REWARD_TYPE = "credit"
	TASK_REWARD_TYPE_JING   TASK_REWARD_TYPE = "jing"
)

type TaskGroup struct {
	GroupId      int64           `orm:"column(groupid);pk"`
	GroupType    TASK_GROUP_TYPE `orm:"column(grouptype)"`
	TaskType     TASK_TYPE       `orm:"column(tasktype)"`
	Name         string          `orm:"column(name)"`
	Description  string          `orm:"column(description)"`
	Icon         int64           `orm:"column(icon)"`
	BgImg        int64           `orm:"column(bgimg)"`
	Tasks        int             `orm:"column(tasks)"`
	Enabled      bool            `orm:"column(enabled)"`
	StartTime    int64           `orm:"column(starttime)"`
	EndTime      int64           `orm:"column(endtime)"`
	DisplayOrder int             `orm:"column(displayorder)"`
	Ex1          string          `orm:"column(ex1)"`
	Ex2          string          `orm:"column(ex2)"`
	Ex3          string          `orm:"column(ex3)"`
	Ex4          string          `orm:"column(ex4)"`
	Ex5          string          `orm:"column(ex5)"`
}

func (self *TaskGroup) TableName() string {
	return "task_groups"
}

func (self *TaskGroup) TableEngine() string {
	return "INNODB"
}

type Task struct {
	TaskId        int64            `orm:"column(taskid);pk"`
	TaskType      TASK_TYPE        `orm:"column(tasktype)"`
	GroupId       int64            `orm:"column(groupid)"`
	Name          string           `orm:"column(name)"`
	Description   string           `orm:"column(description)"`
	Icon          int64            `orm:"column(icon)"`
	TaskLimits    int              `orm:"column(tasklimits)"`
	ApplyPerm     string           `orm:"column(applyperm)"`
	EventName     string           `orm:"column(eventname)"`
	StartTime     int64            `orm:"column(starttime)"`
	EndTime       int64            `orm:"column(endtime)"`
	ResetTime     int              `orm:"column(resettime)"`
	Period        int              `orm:"column(period)"`
	PeriodType    TASK_PERIOD_TYPE `orm:"column(periodtype)"`
	Reward        TASK_REWARD_TYPE `orm:"column(reward)"`
	Prize         int64            `orm:"column(prize)"`
	Applicants    int64            `orm:"column(applicants)"`
	Achievers     int              `orm:"column(achievers)"`
	Version       string           `orm:"column(version)"`
	DisplayOrder  int              `orm:"column(displayorder)"`
	ResetVar      int              `orm:"column(resetvar)"`
	LastResetTime int64            `orm:"column(lastresettime)"`
	Ex1           string           `orm:"column(ex1)"`
	Ex2           string           `orm:"column(ex2)"`
	Ex3           string           `orm:"column(ex3)"`
	Ex4           string           `orm:"column(ex4)"`
	Ex5           string           `orm:"column(ex5)"`
}

func (self *Task) TableName() string {
	return "tasks"
}

func (self *Task) TableEngine() string {
	return "INNODB"
}

type EVENT_TYPE int

const (
	EVENT_TYPE_COUNT EVENT_TYPE = 1 //计数模式
	EVENT_TYPE_NUM   EVENT_TYPE = 2 //总数模式
)

type TaskEvent struct {
	Id          int64      `orm:"column(id);pk"`
	Name        string     `orm:"column(name)"`
	Description string     `orm:"column(description)"`
	Enabled     bool       `orm:"column(enabled)"`
	EventType   EVENT_TYPE `orm:"column(eventtype)"`
}

func (self *TaskEvent) TableName() string {
	return "task_events"
}

func (self *TaskEvent) TableEngine() string {
	return "INNODB"
}

type TaskRecord struct {
	No        string `orm:"column(no);pk"`
	Uid       int64  `orm:"column(uid)"`
	TaskId    int64  `orm:"column(taskid)"`
	Ts        int64  `orm:"column(ts)"`
	RewardNo1 string `orm:"column(rewardno1)"`
	RewardNo2 string `orm:"column(rewardno2)"`
	RewardNo3 string `orm:"column(rewardno3)"`
	RewardNo4 string `orm:"column(rewardno4)"`
	RewardNo5 string `orm:"column(rewardno5)"`
}

func task_records_tbl_tag(uid int64) string {
	mod := uid % 999
	mtag := ""
	if mod > 99 {
		mtag = fmt.Sprintf("%d", mod)
	} else if mod < 10 {
		mtag = fmt.Sprintf("00%d", mod)
	} else {
		mtag = fmt.Sprintf("0%d", mod)
	}
	return mtag
}

func hash_task_records_tbl(uid int64) string {
	tag := task_records_tbl_tag(uid)
	tbl := fmt.Sprintf("%s_%s", db_tbl_pfx, tag)
	if _, ok := taskRecordTbls[tbl]; ok {
		return tbl
	}
	return task_create_tbl(tbl)
}

func task_create_tbl(tbl string) string {
	n_locker.Lock()
	defer n_locker.Unlock()
	if _, ok := taskRecordTbls[tbl]; ok {
		return tbl
	}
	o := dbs.NewOrm(db_aliasname)
	create_tbl_sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(no char(32) NOT NULL,
	  uid int(11) NOT NULL,
	  taskid int(11) NOT NULL,
	  ts int(11) NOT NULL,
	  rewardno1 varchar(30) NOT NULL DEFAULT '',
	  rewardno2 varchar(30) NOT NULL DEFAULT '',
	  rewardno3 varchar(30) NOT NULL DEFAULT '',
	  rewardno4 varchar(30) NOT NULL DEFAULT '',
	  rewardno5 varchar(30) NOT NULL DEFAULT '',
	  PRIMARY KEY (no),
	  KEY idx_uid_taskid_ts (uid,taskid,ts) USING BTREE) ENGINE=InnoDB DEFAULT CHARSET=utf8`, tbl)
	_, err := o.Raw(create_tbl_sql).Exec()
	if err == nil {
		taskRecordTbls[tbl] = true
	}
	return tbl
}

func BuildTaskRecordNo(uid int64, task *Task, t time.Time) string {
	fstr := ""
	switch task.TaskType {
	case TASK_TYPE_DUPLICATE:
		fstr = "U%dS%dT%sV%d"
		break
	default:
		fstr = "U%dS%dV%dONCE"
		return fmt.Sprintf(fstr, uid, task.TaskId, task.ResetVar)
	}
	ts := fmt.Sprintf("%d", t.Year())
	if task.PeriodType == TASK_PERIOD_TYPE_DAY {
		ts = fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	}
	return fmt.Sprintf(fstr, uid, task.TaskId, ts, task.ResetVar)
}
