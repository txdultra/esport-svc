package groups

import (
	"libs/message"
	"strconv"
	"strings"
)

const (
	MSG_TYPE_INVITED message.MSG_TYPE = "group:invite"
	MSG_TYPE_MESSAGE message.MSG_TYPE = "group:msg"
)

type GroupCfg struct {
	Id                           int64         `orm:"column(id);pk"`
	GroupNameLen                 int           `orm:"column(groupname_len)"`
	GroupDescMaxLen              int           `orm:"column(groupdesc_maxlen)"`
	GroupDescMinLen              int           `orm:"column(groupdesc_minlen)"`
	CreateGroupBasePoint         int64         `orm:"column(creategroup_basepoint)"`
	CreateGroupRate              float32       `orm:"column(creategroup_rate)"`
	CreateGroupMinUsers          int           `orm:"column(creategroup_minusers)"`
	CreateGroupRecruitDay        int           `orm:"column(creategroup_recruitday)"`
	CreateGroupMaxCount          int           `orm:"column(creategroup_maxcount)"`
	CreateGroupCertifiedMaxCount int           `orm:"column(creategroup_certifiedmaxcount)"`
	CreateGroupClause            string        `orm:"column(creategroup_clause)"`
	ReportOptions                string        `orm:"column(report_options)"`
	NewThreadDefaultStatus       THREAD_STATUS `orm:"column(new_threadstatus)"`
}

func (self *GroupCfg) TableName() string {
	return "config"
}

func (self *GroupCfg) TableEngine() string {
	return "INNODB"
}

type GroupMemberTable struct {
	Id      int64  `orm:"column(id);pk"`
	TblName string `orm:"column(tablename)"`
	Ts      int64  `orm:"column(ts)"`
}

func (self *GroupMemberTable) TableName() string {
	return "tables_groupm"
}

func (self *GroupMemberTable) TableEngine() string {
	return "INNODB"
}

type MemberGroupTable struct {
	Id      int64  `orm:"column(id);pk"`
	TblName string `orm:"column(tablename)"`
	Ts      int64  `orm:"column(ts)"`
}

func (self *MemberGroupTable) TableName() string {
	return "tables_mgroup"
}

func (self *MemberGroupTable) TableEngine() string {
	return "INNODB"
}

type PostTable struct {
	Id      int64  `orm:"column(id);pk"`
	TblName string `orm:"column(tablename)"`
	Ts      int64  `orm:"column(ts)"`
}

func (self *PostTable) TableName() string {
	return "tables_post"
}

func (self *PostTable) TableEngine() string {
	return "INNODB"
}

type GROUP_STATUS int

const (
	GROUP_STATUS_OPENING    GROUP_STATUS = 1
	GROUP_STATUS_RECRUITING GROUP_STATUS = 2
	GROUP_STATUS_LOWMEMBER  GROUP_STATUS = 3
	GROUP_STATUS_CLOSED     GROUP_STATUS = 4
)

type GROUP_BELONG int

const (
	GROUP_BELONG_MEMBER   GROUP_BELONG = 1
	GROUP_BELONG_OFFICIAL GROUP_BELONG = 10
)

type GROUP_TYPE int

const (
	GROUP_TYPE_NORMAL GROUP_TYPE = 1
)

type Group struct {
	Id             int64        `orm:"column(id);pk"`
	MembersTableId int          `orm:"column(members_tableid)"`
	ThreadTableId  int          `orm:"column(thread_tableid)"`
	Name           string       `orm:"column(groupname)"`
	Description    string       `orm:"column(description)"`
	Uid            int64        `orm:"column(uid)"`
	CreateTime     int64        `orm:"column(createtime)"`
	Members        int          `orm:"column(members)"`
	Threads        int          `orm:"column(threads)"`
	Country        string       `orm:"column(country)"`
	City           string       `orm:"column(city)"`
	GameIds        string       `orm:"column(gameids)"`
	DisplarOrder   int          `orm:"column(displayorder)"`
	Status         GROUP_STATUS `orm:"column(status)"`
	Img            int64        `orm:"column(img)"`
	BgImg          int64        `orm:"column(bgimg)"`
	Belong         GROUP_BELONG `orm:"column(belong)"`
	Type           GROUP_TYPE   `orm:"column(type)"`
	Vitality       int          `orm:"column(vitality)"`
	SearchKeyword  string       `orm:"column(searchkeyword)"`
	LongiTude      float32      `orm:"column(longitude)"`
	LatiTude       float32      `orm:"column(latitude)"`
	Recommend      bool         `orm:"column(recommend)"`
	StartTime      int64        `orm:"column(starttime)"`
	EndTime        int64        `orm:"column(endtime)"`
	MinUsers       int          `orm:"column(min_users)"`
	OrderNo        string       `orm:"column(orderno)"`
	InviteUids     []int64      `orm:"-"`
}

func (self *Group) TableName() string {
	return "groups"
}

func (self *Group) TableEngine() string {
	return "INNODB"
}

func (self *Group) GameIDs() []int {
	gids := strings.Split(self.GameIds, ",")
	ids := []int{}
	for _, gid := range gids {
		id, err := strconv.Atoi(gid)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

type GroupMember struct {
	GroupId int64 `orm:"column(groupid)"`
	Uid     int64 `orm:"column(uid)"`
	Ts      int64 `orm:"column(ts)"`
}

type MemberGroup struct {
	Uid     int64 `orm:"column(uid)"`
	GroupId int64 `orm:"column(groupid)"`
	Ts      int64 `orm:"column(ts)"`
}

type InviteMember struct {
	Uid     int64
	Invited bool
	Joined  bool
}

type THREAD_STATUS int

const (
	THREAD_STATUS_OPENING  THREAD_STATUS = 1
	THREAD_STATUS_LOCKING  THREAD_STATUS = 2
	THREAD_STATUS_AUDITING THREAD_STATUS = 3
)

type Thread struct {
	Id           int64         `orm:"column(id);pk"`
	GroupId      int64         `orm:"column(groupid)"`
	PostTableId  int           `orm:"column(posts_tableid)"`
	Subject      string        `orm:"column(subject)"`
	AuthorId     int64         `orm:"column(authorid)"`
	Author       string        `orm:"column(author)"`
	DateLine     int64         `orm:"column(dateline)"`
	LastPost     int64         `orm:"column(lastpost)"`
	LastPoster   string        `orm:"column(lastposter)"`
	LastPostUid  int64         `orm:"column(lastpostuid)"`
	LastId       string        `orm:"column(lastid)"`
	Views        int           `orm:"column(views)"`
	Replies      int           `orm:"column(replies)"`
	Shares       int           `orm:"column(shares)"`
	Favorites    int           `orm:"column(favs)"`
	Status       THREAD_STATUS `orm:"column(status)"`
	Img          int64         `orm:"column(img)"`
	Closed       bool          `orm:"column(closed)"`
	Highlight    bool          `orm:"column(highlight)"`
	Heats        int           `orm:"column(heats)"`
	DisplayOrder int           `orm:"column(displayorder)"`
	HasVod       bool          `orm:"column(hasvod)"`
	Lordpid      string        `orm:"column(lordpid)"`
}

func (self *Thread) TableName() string {
	return "thread"
}

func (self *Thread) TableEngine() string {
	return "INNODB"
}

type FROM_DEVICE int

const (
	FROM_DEVICE_IPHONE  FROM_DEVICE = 1
	FROM_DEVICE_IPAD    FROM_DEVICE = 2
	FROM_DEVICE_ANDROID FROM_DEVICE = 3
	FROM_DEVICE_WEB     FROM_DEVICE = 4
	FROM_DEVICE_WPHONE  FROM_DEVICE = 5
	FROM_DEVICE_OTHER   FROM_DEVICE = 6
)

func GetFromDevice(name string) FROM_DEVICE {
	switch name {
	case "android":
		return FROM_DEVICE_ANDROID
	case "ios":
		return FROM_DEVICE_IPHONE
	case "ipad":
		return FROM_DEVICE_IPAD
	case "wphone":
		return FROM_DEVICE_WPHONE
	case "web":
		return FROM_DEVICE_WEB
	default:
		return FROM_DEVICE_OTHER
	}
}

type Post struct {
	Id            string      `orm:"column(id);pk"`
	ThreadId      int64       `orm:"column(tid)"`
	AuthorId      int64       `orm:"column(authorid)"`
	Subject       string      `orm:"column(subject)"`
	DateLine      int64       `orm:"column(dateline)"`
	Message       string      `orm:"column(message)"`
	Ip            string      `orm:"column(ip)"`
	Invisible     bool        `orm:"column(invisible)"`
	Ding          int         `orm:"column(ding)"`
	Cai           int         `orm:"column(cai)"`
	Position      int         `orm:"column(position)"`
	ReplyId       string      `orm:"column(replyid)"`
	ReplyUid      int64       `orm:"column(replyuid)"`
	ReplyPosition int         `orm:"column(replyposition)"`
	Img           int64       `orm:"column(img)"`
	Resources     string      `orm:"column(resources)"`
	LongiTude     float32     `orm:"column(longitude)"`
	LatiTude      float32     `orm:"column(latitude)"`
	FromDevice    FROM_DEVICE `orm:"column(fromdev)"`
	ImgIds        []int64     `orm:"-"`
	VodIds        []int64     `orm:"-"`
}

func (self *Post) TableName() string {
	return "post"
}

func (self *Post) TableEngine() string {
	return "INNODB"
}

type REPORT_CATEGORY int

const (
	REPORT_CATEGORY_POST  REPORT_CATEGORY = 1
	REPORT_CATEGORY_GROUP REPORT_CATEGORY = 2
)

type Report struct {
	Id      int64           `orm:"column(id);pk"`
	RefId   string          `orm:"column(refid)"`
	C       REPORT_CATEGORY `orm:"column(c)"`
	Ts      int64           `orm:"column(ts)"`
	RefTxt  string          `orm:"column(reftxt)"`
	PostUid int64           `orm:"column(postuid)"`
	Msg     string          `orm:"column(msg)"`
}

func (self *Report) TableName() string {
	return "reports"
}

func (self *Report) TableEngine() string {
	return "INNODB"
}

type MemberCount struct {
	Uid       int64 `orm:"column(uid);pk"`
	Groups    int   `orm:"column(groups)"`
	ToDings   int   `orm:"column(todings)"`
	ToCais    int   `orm:"column(tocais)"`
	FromDings int   `orm:"column(fromdings)"`
	FromCais  int   `orm:"column(fromcais)"`
	Posts     int   `orm:"column(posts)"`
	Threads   int   `orm:"column(threads)"`
	Joins     int   `orm:"column(joins)"`
	Shares    int   `orm:"column(shares)"`
	Reports   int   `orm:"column(reports)"`
	LastTs    int64 `orm:"column(lastts)"`
}

func (self *MemberCount) TableName() string {
	return "member_count"
}

func (self *MemberCount) TableEngine() string {
	return "INNODB"
}
