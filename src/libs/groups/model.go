package groups

type GroupCfg struct {
	Id                    int64   `orm:"column(id);pk"`
	GroupNameLen          int     `orm:"column(groupname_len)"`
	GroupDescLen          int     `orm:"column(groupdesc_len)"`
	CreateGroupBasePoint  int     `orm:"column(creategroup_basepoint)"`
	CreateGroupRate       float32 `orm:"column(creategroup_rate)"`
	CreateGroupMinUsers   int     `orm:"column(creategroup_minusers)"`
	CreateGroupRecruitDay int     `orm:"column(creategroup_recruitday)"`
	CreateGroupMaxCount   int     `orm:"column(creategroup_maxcount)"`
	CreateGroupClause     string  `orm:"column(creategroup_clause)"`
	ReportOptions         string  `orm:"column(report_options)"`
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
	Ts      int    `orm:"column(ts)"`
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
	Ts      int    `orm:"column(ts)"`
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
	Ts      int    `orm:"column(ts)"`
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
}

func (self *Group) TableName() string {
	return "group"
}

func (self *Group) TableEngine() string {
	return "INNODB"
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

type THREAD_STATUS int

const (
	THREAD_STATUS_OPENING  THREAD_STATUS = 1
	THREAD_STATUS_LOCKING  THREAD_STATUS = 2
	THREAD_STATUS_AUDITING THREAD_STATUS = 3
)

type Thread struct {
	Id           int64         `orm:"column(id);pk"`
	PostTableId  int           `orm:"column(posts_tableid)"`
	Subject      string        `orm:"column(subject)"`
	AuthorId     int64         `orm:"column(authorid)"`
	Author       string        `orm:"column(author)"`
	DateLine     int           `orm:"column(dateline)"`
	LastPost     int64         `orm:"column(lastpost)"`
	LastPoster   string        `orm:"column(lastposter)"`
	LastPostUid  int64         `orm:"column(lastpostuid)"`
	LastId       int64         `orm:"column(lastid)"`
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
	Lordpid      int64         `orm:"column(lordpid)"`
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
)

type Post struct {
	Id            int64       `orm:"column(id);pk"`
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
	ReplyId       int64       `orm:"column(replyid)"`
	ReplyUid      int64       `orm:"column(replyuid)"`
	ReplyPosition int         `orm:"column(replyposition)"`
	Img           int64       `orm:"column(img)"`
	Resources     string      `orm:"column(resources)"`
	LongiTude     float32     `orm:"column(longitude)"`
	LatiTude      float32     `orm:"column(latitude)"`
	FromDevice    FROM_DEVICE `orm:"column(fromdev)"`
}

func (self *Post) TableName() string {
	return "post"
}

func (self *Post) TableEngine() string {
	return "INNODB"
}

//type GroupRecruit struct {
//	GroupId    int64  `orm:"column(groupid)"`
//	CreateTime int    `orm:"column(createtime)"`
//	Uid        int64  `orm:"column(uid)"`
//	StartTime  int    `orm:"column(starttime)"`
//	EndTime    int    `orm:"column(endtime)"`
//	MinUsers   int    `orm:"column(min_users)"`
//	OrderNo    string `orm:"column(orderno)"`
//}

//func (self *GroupRecruit) TableName() string {
//	return "group_recruit_temp"
//}

//func (self *GroupRecruit) TableEngine() string {
//	return "INNODB"
//}

//type GroupReadyStop struct {
//	GroupId    int64 `orm:"column(groupid)"`
//	CreateTime int   `orm:"column(createtime)"`
//	Uid        int64 `orm:"column(uid)"`
//	StartTime  int   `orm:"column(starttime)"`
//	EndTime    int   `orm:"column(endtime)"`
//	MinUsers   int   `orm:"column(min_users)"`
//}

//func (self *GroupReadyStop) TableName() string {
//	return "group_readystop_temp"
//}

//func (self *GroupReadyStop) TableEngine() string {
//	return "INNODB"
//}

type REPORT_CATEGORY int

const (
	REPORT_CATEGORY_POST  REPORT_CATEGORY = 1
	REPORT_CATEGORY_GROUP REPORT_CATEGORY = 2
)

type Report struct {
	Id      int64           `orm:"column(id);pk"`
	RefId   int64           `orm:"column(refid)"`
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
