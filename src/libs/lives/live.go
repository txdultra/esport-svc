package lives

import (
	"libs/reptile"
	"time"

	"labix.org/v2/mgo/bson"
)

type LivePerson struct {
	Id         int64
	Name       string `orm:"column(title)"`
	Img        int64
	Uid        int64
	Des        string
	ReptileUrl string
	Rep        reptile.REP_SUPPORT
	StreamUrl  string
	Succ       bool                `orm:"column(rep_succ)"`
	ReptileDes string              `orm:"column(rep_des)"`
	LiveStatus reptile.LIVE_STATUS `orm:"column(live_status)"`
	PostTime   time.Time
	LastTime   time.Time
	Enabled    bool `orm:"column(enabled)"`
	//RepMethod  reptile.REP_METHOD
	ShowOnlineMin int `orm:"column(show_online_min)"`
	ShowOnlineMax int `orm:"column(show_online_max)"`
}

func (self *LivePerson) TableName() string {
	return "live_personal"
}

func (self *LivePerson) TableEngine() string {
	return "INNODB"
}

//主频道
type LiveChannel struct {
	Id      int64
	Name    string `orm:"column(title)"`
	Img     int64
	Childs  int
	Uid     int64
	Enabled bool
}

func (self *LiveChannel) TableName() string {
	return "live_channel"
}

func (self *LiveChannel) TableEngine() string {
	return "INNODB"
}

//直播流
//type LIVE_REP_SUPPORT string

//const (
//	LIVE_REP_17173  LIVE_REP_SUPPORT = "17173"
//	LIVE_REP_PPLIVE LIVE_REP_SUPPORT = "pplive"
//	LIVE_REP_FYZB   LIVE_REP_SUPPORT = "fyzb"
//)

type LiveStream struct {
	Id         int64
	Rep        reptile.REP_SUPPORT
	Img        int64
	ReptileUrl string
	StreamUrl  string
	LastTime   time.Time
	ChannelId  int64 `orm:"column(pid)"`
	Default    bool  `orm:"column(def)"`
	Enabled    bool
	AllowRep   bool
}

func (self *LiveStream) TableName() string {
	return "live_streams"
}

func (self *LiveStream) TableEngine() string {
	return "INNODB"
}

type LiveProgram struct {
	Id               int64
	Title            string
	SubTitle         string
	Date             time.Time `orm:"column(d);type(date)"`
	StartTime        time.Time
	EndTime          time.Time
	MatchId          int
	PostTime         time.Time
	PostUid          int64
	GameId           int   `orm:"column(gid)"`
	Img              int64 `orm:"column(img)"`
	DefaultChannelId int64 `orm:"column(def_cid)"`
}

func (self *LiveProgram) TableName() string {
	return "live_program"
}

func (self *LiveProgram) TableEngine() string {
	return "INNODB"
}

type LIVE_SUBPROGRAM_VIEW_TYPE int

const (
	LIVE_SUBPROGRAM_VIEW_VS     LIVE_SUBPROGRAM_VIEW_TYPE = 1 //对阵
	LIVE_SUBPROGRAM_VIEW_SINGLE LIVE_SUBPROGRAM_VIEW_TYPE = 2 //单个
)

type LiveSubProgram struct {
	Id        int64
	ProgramId int64                     `orm:"column(pid)"`
	GameId    int                       `orm:"column(gid)"`
	Vs1Name   string                    `orm:"column(vs1)"`
	Vs1Img    int64                     `orm:"column(vs1_img)"`
	Vs1Uid    int64                     `orm:"column(vs1_uid)"`
	Vs2Name   string                    `orm:"column(vs2)"`
	Vs2Img    int64                     `orm:"column(vs2_img)"`
	Vs2Uid    int64                     `orm:"column(vs2_uid)"`
	ViewType  LIVE_SUBPROGRAM_VIEW_TYPE `orm:"column(t)"`
	Title     string
	Img       int64
	StartTime time.Time
	EndTime   time.Time
	PostTime  time.Time
	PostUid   int64
}

func (self *LiveSubProgram) TableName() string {
	return "live_subprogram"
}

func (self *LiveSubProgram) TableEngine() string {
	return "INNODB"
}

//用户提交的原始提醒列表
type CommitOriginalNotices struct {
	ID       bson.ObjectId `bson:"_id"`
	EventId  int64         `bson:"event_id"`
	FromId   int64         `bson:"from_id"`
	RefIds   []int64       `bson:"refids"`
	LastTime time.Time     `bson:"last_time"`
}

type LIVE_TYPE int //直播类型

const (
	LIVE_TYPE_ORGANIZATION LIVE_TYPE = 1
	LIVE_TYPE_PERSONAL     LIVE_TYPE = 2
)

//个人直播 用于sphinx搜索数据类型
type LiveSearchData struct {
	LiveId              int64               `orm:"column(id);pk"`
	PersonalChannelName string              `orm:"column(pc_name)"`
	AnchorName          string              `orm:"column(anchor)"`
	Status              reptile.LIVE_STATUS `orm:"column(status)"`
	GameNames           string              `orm:"column(games)"`
	GameIds             string              `orm:"column(game_ids)"`
	CustomWeight        int                 `orm:"column(cw)"`
	Onlines             int                 `orm:"column(onlines)"`
	Enabled             bool                `orm:"column(enabled)"`
}

func (self *LiveSearchData) TableName() string {
	return "live_serach_data"
}

func (self *LiveSearchData) TableEngine() string {
	return "INNODB"
}

//live节目单 用于sphinx搜索数据类型
type LiveProgramSearchData struct {
	ProgramId    int64     `orm:"column(id);pk"`
	Title        string    `orm:"column(title)"`
	SubTitle     string    `orm:"column(sub_title)"`
	ChannelName  string    `orm:"column(channel_name)"`
	GameName     string    `orm:"column(game_name)"`
	StartTime    time.Time `orm:"column(stime)"`
	EndTime      time.Time `orm:"column(etime)"`
	CustomWeight int       `orm:"column(cw)"`
	Onlines      int       `orm:"column(onlines)"`
}

func (self *LiveProgramSearchData) TableName() string {
	return "live_program_serach_data"
}

func (self *LiveProgramSearchData) TableEngine() string {
	return "INNODB"
}
