package libs

import (
	"time"
)

//type BaseModel struct {
//	Ext map[string]interface{} `orm:"-"`
//}

type Game struct {
	Id           int
	Name         string `orm:"column(title)"`
	En           string
	Img          int64
	Enabled      bool
	PostTime     time.Time
	DisplayOrder int
}

func (self *Game) TableName() string {
	return "common_games"
}

func (self *Game) TableEngine() string {
	return "INNODB"
}

type Match struct {
	Id       int
	Name     string `orm:"column(title)"`
	SubTitle string
	En       string
	Img      int64
	Des1     string
	Des2     string
	Des3     string
}

func (self *Match) TableName() string {
	return "common_matchs"
}

func (self *Match) TableEngine() string {
	return "INNODB"
}

type File struct {
	Id           int64
	FileName     string
	OriginalName string
	Volume       string
	ExtName      string `orm:"column(ext)"`
	Size         int64
	PostTime     time.Time
	MimeType     string `orm:"column(mime)"`
	IsDeleted    bool   `orm:"column(is_del)"`
	Width        int    `orm:"column(w)"`
	Height       int    `orm:"column(h)"`
	Source       int64
	C            string
}

func (self *File) TableName() string {
	return "common_files"
}

func (self *File) TableEngine() string {
	return "INNODB"
}

////////////////////////////////////////////////////////////////////////////////
//MSQ
////////////////////////////////////////////////////////////////////////////////

type MSQ_STATE string

const (
	MSQ_STATE_RUNNING   MSQ_STATE = "running"
	MSQ_STATE_READY     MSQ_STATE = "ready"
	MSQ_STATE_COMPLETED MSQ_STATE = "completed"
	MSQ_STATE_DISCARD   MSQ_STATE = "discard"
)

//type MSQ_QUEUE_STATE string

//const (
//	MSQ_QUEUE_STATE_CREATED  MSQ_QUEUE_STATE = "created"
//	MSQ_QUEUE_STATE_REDELETE MSQ_QUEUE_STATE = "ready_del"
//	MSQ_QUEUE_STATE_DELETED  MSQ_QUEUE_STATE = "deleted"
//)

type Msqtor struct {
	Id           int64
	MsqId        string `orm:"column(msqid)"`
	ScheduleTime time.Time
	CfgJson      string `orm:"column(cfg)"`
	ConsumerType string
	Consumers    int16
	Status       MSQ_STATE
	Des          string
	CreateTime   time.Time
	//QueueState   MSQ_QUEUE_STATE
	CompletedDel bool `orm:"column(cpd_del)"`
}

func (self *Msqtor) TableName() string {
	return "common_msq"
}

func (self *Msqtor) TableEngine() string {
	return "INNODB"
}

////////////////////////////////////////////////////////////////////////////////
//PUSH STATE
////////////////////////////////////////////////////////////////////////////////
type PUSH_STATE string

const (
	PUSH_STATE_NOSET   PUSH_STATE = "noset"
	PUSH_STATE_READY   PUSH_STATE = "ready"
	PUSH_STATE_SENDING PUSH_STATE = "sending"
	PUSH_STATE_SENDED  PUSH_STATE = "sended"
)

type PushState struct {
	Id         int64
	EventId    int
	EventClass string `orm:"column(event_ct)"`
	State      PUSH_STATE
	LastTime   time.Time `orm:"auto_now"`
}

func (self *PushState) TableName() string {
	return "common_push_state"
}

func (self *PushState) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (self *PushState) TableUnique() [][]string {
	return [][]string{
		[]string{"ref_id", "ref_ct"},
	}
}

////////////////////////////////////////////////////////////////////////////////
//recommend
////////////////////////////////////////////////////////////////////////////////

type Recommend struct {
	Id           int64
	RefId        int64
	RefType      string
	Title        string
	Img          int64
	Category     string
	DisplayOrder int
	Enabled      bool
	PostTime     time.Time
	PostUid      int64
}

func (self *Recommend) TableName() string {
	return "common_recommends"
}

func (self *Recommend) TableEngine() string {
	return "INNODB"
}

////////////////////////////////////////////////////////////////////////////////
//smiley
////////////////////////////////////////////////////////////////////////////////
type Smiley struct {
	Id           int64
	Code         string `orm:"unique"`
	Img          int64
	ImgPath      string
	Category     string
	DisplayOrder int
	Points       int
}

func (self *Smiley) TableName() string {
	return "common_smiley"
}

func (self *Smiley) TableEngine() string {
	return "INNODB"
}
