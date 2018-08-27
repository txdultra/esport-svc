package matchrace

import (
	"time"
)

type MODE_TYPE int16

const (
	MODE_TYPE_GROUP  MODE_TYPE = 1 //小组模型
	MODE_TYPE_RECENT MODE_TYPE = 2 //胜负模型
	MODE_TYPE_ELIMIN MODE_TYPE = 3 //淘汰模型
)

type RaceMode struct {
	Id           int64     `orm:"column(id);pk"`
	MatchId      int64     `orm:"column(matchid)"`
	ModeType     MODE_TYPE `orm:"column(modetype)"`
	DisplayOrder int       `orm:"column(displayorder)"`
	IsView       bool      `orm:"column(isview)"`
	Title        string    `orm:"column(title)"`
}

func (self *RaceMode) TableName() string {
	return "match_modes"
}

func (self *RaceMode) TableEngine() string {
	return "INNODB"
}

type MatchPlayer struct {
	Id       int64     `orm:"column(id);pk"`
	Name     string    `orm:"column(pname)"`
	Img      int64     `orm:"column(img)"`
	PostTime time.Time `orm:"column(posttime)"`
}

func (self *MatchPlayer) TableName() string {
	return "match_players"
}

func (self *MatchPlayer) TableEngine() string {
	return "INNODB"
}

type WLP int16

const (
	WLP_UNDEFINED WLP = 0 //未知
	WLP_W         WLP = 1 //赢
	WLP_L         WLP = 2 //输
	WLP_P         WLP = 3 //平
)

type MatchRecent struct {
	Id           int64 `orm:"column(id);pk"`
	ModeId       int64 `orm:"column(modeid)"`
	Player       int64 `orm:"column(player)"`
	M1           WLP   `orm:"column(m1)"`
	M2           WLP   `orm:"column(m2)"`
	M3           WLP   `orm:"column(m3)"`
	M4           WLP   `orm:"column(m4)"`
	M5           WLP   `orm:"column(m5)"`
	DisplayOrder int   `orm:"column(displayorder)"`
	Disabled     bool  `orm:"column(disabled)"`
}

func (self *MatchRecent) TableName() string {
	return "match_recents"
}

func (self *MatchRecent) TableEngine() string {
	return "INNODB"
}

type MatchGroup struct {
	Id           int64     `orm:"column(id);pk"`
	ModeId       int64     `orm:"column(modeid)"`
	Title        string    `orm:"column(title)"`
	DisplayOrder int       `orm:"column(displayorder)"`
	PostTime     time.Time `orm:"column(posttime)"`
}

func (self *MatchGroup) TableName() string {
	return "match_groups"
}

func (self *MatchGroup) TableEngine() string {
	return "INNODB"
}

type MatchGroupPlayer struct {
	Id           int64 `orm:"column(id);pk"`
	GroupId      int64 `orm:"column(groupid)"`
	Player       int64 `orm:"column(player)"`
	Wins         int16 `orm:"column(wins)"`
	Pings        int16 `orm:"column(pings)"`
	Loses        int16 `orm:"column(loses)"`
	Points       int   `orm:"column(points)"`
	DisplayOrder int   `orm:"column(displayorder)"`
	Outlet       bool  `orm:"column(outlet)"`
}

func (self *MatchGroupPlayer) TableName() string {
	return "match_group_players"
}

func (self *MatchGroupPlayer) TableEngine() string {
	return "INNODB"
}

type MatchVs struct {
	Id      int64  `orm:"column(id);pk"`
	A       int64  `orm:"column(a)"`
	AName   string `orm:"column(a_name)"`
	AImg    int64  `orm:"column(a_img)"`
	AScore  int16  `orm:"column(a_score)"`
	B       int64  `orm:"column(b)"`
	BName   string `orm:"column(b_name)"`
	BImg    int64  `orm:"column(b_img)"`
	BScore  int16  `orm:"column(b_score)"`
	MatchId int64  `orm:"column(matchid)"`
	ModeId  int64  `orm:"column(modeid)"`
	RefId   int64  `orm:"column(refid)"`
}

func (self *MatchVs) TableName() string {
	return "match_vs"
}

func (self *MatchVs) TableEngine() string {
	return "INNODB"
}

type ELIMIN_MSTYPE int16

const (
	ELIMIN_MSTYPE_STANDARD ELIMIN_MSTYPE = 0 //标准
	ELIMIN_MSTYPE_CHAMPION ELIMIN_MSTYPE = 1 //冠军
	ELIMIN_MSTYPE_THIRD    ELIMIN_MSTYPE = 2 //季军
)

type MatchEliminMs struct {
	Id           int64         `orm:"column(id);pk"`
	ModeId       int64         `orm:"column(modeid)"`
	Title        string        `orm:"column(title)"`
	PostTime     int64         `orm:"column(posttime)"`
	Icon         int64         `orm:"column(icon)"`
	DisplayOrder int           `orm:"column(displayorder)"`
	T            ELIMIN_MSTYPE `orm:"column(t)"`
}

func (self *MatchEliminMs) TableName() string {
	return "match_elimin_ms"
}

func (self *MatchEliminMs) TableEngine() string {
	return "INNODB"
}

type MatchEliminVs struct {
	Id       int64 `orm:"column(id);pk"`
	MsId     int64 `orm:"column(msid)"`
	VsId     int64 `orm:"column(vsid)"`
	OutletId int64 `orm:"column(outletid)"`
	PostTime int64 `orm:"column(posttime)"`
}

func (self *MatchEliminVs) TableName() string {
	return "match_elimin_vs"
}

func (self *MatchEliminVs) TableEngine() string {
	return "INNODB"
}
