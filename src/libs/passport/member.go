package passport

import (
	"encoding/json"
	"libs/vars"
	"logs"
	"strconv"
	"strings"
	"time"
)

type MEMBER_REGISTER_MODE int

const (
	MEMBER_REGISTER_OPENID MEMBER_REGISTER_MODE = 1
	MEMBER_REGISTER_EMAIL  MEMBER_REGISTER_MODE = 2
	MEMBER_REGISTER_MOBILE MEMBER_REGISTER_MODE = 4
	MEMBER_REGISTER_NORMAL MEMBER_REGISTER_MODE = 8
)

type Member struct {
	Uid              int64                `orm:"auto;pk"`
	UserName         string               `orm:"column(user_name)"`
	NickName         string               `orm:"column(nick_name)"`
	Email            string               `orm:"column(email)"`
	Password         string               `orm:"column(password)"`
	Salt             string               `orm:"column(salt)"`
	MobileIdentifier string               `orm:"size(32)"`
	MemberIdentifier string               `orm:"size(32)"`
	CreateTime       int64                `orm:"column(create_time)"`
	CreateIP         int64                `orm:"column(create_ip)"`
	Avatar           int64                `orm:"column(avatar)"`
	RegMode          MEMBER_REGISTER_MODE `orm:"column(reg_type)"`
	PushId           string               `orm:"column(push_id)"`
	PushChannelId    string               `orm:"column(channel_id)"`
	PushProxy        int                  `orm:"column(push_pxy)"`
	DeviceType       vars.CLIENT_OS       `orm:"column(device_type)"`
	Certified        bool                 `orm:"column(certified)"`
	CertifiedReason  string               `orm:"column(certified_reason)"`
	Gids             string               `orm:"column(gids)"`
}

func (self *Member) TableName() string {
	return "common_member"
}

func (self *Member) TableEngine() string {
	return "INNODB"
}

func (self *Member) ConvertToGids(gameIds []int) string {
	_gameids := []string{}
	for _, _gid := range gameIds {
		_gameids = append(_gameids, strconv.Itoa(_gid))
	}
	return strings.Join(_gameids, ",")
}

func (self *Member) GameIds() []int {
	gids := []int{}
	gameids := strings.Split(self.Gids, ",")
	for _, str := range gameids {
		_id, err := strconv.Atoi(str)
		if err != nil {
			logs.Errorf("member method GameIds transfor id fail:%s", err.Error())
			continue
		}
		if _id <= 0 {
			continue
		}
		gids = append(gids, _id)
	}
	return gids
}

type MemberAvatar struct {
	Uid    int64         `orm:"column(uid)"`
	FileId int64         `orm:"column(fid)"`
	Size   vars.PIC_SIZE `orm:"column(size)"`
	Width  int           `orm:"column(w)"`
	Height int           `orm:"column(h)"`
	Ts     int64         `orm:"column(ts)"`
}

type MemberState struct {
	Uid           int64 `orm:"pk;unique"`
	Logins        int
	LastLoginIP   int64 `orm:"column(last_login_ip)"`
	LastLoginTime int64
	LongiTude     float32 `orm:"column(longitude)"`
	LatiTude      float32 `orm:"column(latitude)"`
	Vods          int
	Fans          int
	Friends       int
	Notes         int
	Count1        int `orm:"column(count_1)"`
	Count2        int `orm:"column(count_2)"`
	Count3        int `orm:"column(count_3)"`
	Count4        int `orm:"column(count_4)"`
	Count5        int `orm:"column(count_5)"`
}

func (self *MemberState) TableName() string {
	return "common_member_states"
}

func (self *MemberState) TableEngine() string {
	return "INNODB"
}

type MemberNickName struct {
	Uid      int64     `orm:"pk"`
	NickName string    `orm:"unique"`
	PostTime time.Time `orm:"auto_now"`
}

func (self *MemberNickName) TableName() string {
	return "common_member_nicknames"
}

func (self *MemberNickName) TableEngine() string {
	return "INNODB"
}

type MemberProfile struct {
	Uid            int64 `orm:"pk"`
	BackgroundImg  int64 `orm:"column(bg)"`
	Description    string
	Gender         string
	RealName       string `orm:"column(realname)"`
	BirthYear      int    `orm:"column(birthyear)"`
	BirthMonth     int    `orm:"column(birthmonth)"`
	BirthDay       int    `orm:"column(birthday)"`
	Mobile         string
	FoundPwdMobile string `orm:"column(repwd_mobile)"`
	IDCard         string `orm:"column(idcard)"`
	QQ             string `orm:"column(qq)"`
	Alipay         string
	Youku          string
	Bio            string
	Field1         string `orm:"column(field1)"`
	Field2         string `orm:"column(field2)"`
	Field3         string `orm:"column(field3)"`
	Field4         string `orm:"column(field4)"`
	Field5         string `orm:"column(field5)"`
}

func (self *MemberProfile) TableName() string {
	return "common_member_profile"
}

func (self *MemberProfile) TableEngine() string {
	return "INNODB"
}

type LoginsOver struct {
	Id       int
	UserName string `orm:column(username)`
	Logins   int
	PostTime time.Time
	Ip       int64
}

func (self *LoginsOver) TableName() string {
	return "common_logins_over"
}

func (self *LoginsOver) TableEngine() string {
	return "INNODB"
}

//type MemberFriend struct {
//	Id        int64
//	SourceUid int64
//	TargetUid int64
//	PostTime  time.Time
//}

//func (self *MemberFriend) TableName() string {
//	return "common_member_friends"
//}

//func (self *MemberFriend) TableEngine() string {
//	return "INNODB"
//}

// 多字段唯一键
//func (self *MemberFriend) TableUnique() [][]string {
//	return [][]string{
//		[]string{"source_uid", "target_uid"},
//	}
//}

type MemberRelation struct {
	Uid                  int64
	Following            bool
	FollowedBy           bool
	NotificationsEnabled bool
}

type STMemberRelation struct {
	Source MemberRelation
	Target MemberRelation
}

type MemberGame struct {
	Uid      int64     `orm:"column(uid)"`
	GameId   int       `orm:"column(gid)"`
	PostTime time.Time `orm:"column(post_time)"`
}

type ManageMemberMap struct {
	Id       int64
	FromUid  int64  `orm:"column(from_uid)"`
	ToUid    int64  `orm:"column(to_uid)"`
	Platform string `orm:"column(plat)"`
}

func (self *ManageMemberMap) TableName() string {
	return "common_manage_members"
}

func (self *ManageMemberMap) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (self *ManageMemberMap) TableUnique() [][]string {
	return [][]string{
		[]string{"from_uid", "plat"},
	}
}

type MemberConfig struct {
	Uid     int64 `orm:"pk"`
	Setting string
}

func (self *MemberConfig) TableName() string {
	return "common_member_config"
}

func (self *MemberConfig) TableEngine() string {
	return "INNODB"
}

func (self *MemberConfig) MemberConfigAttrs() *MemberConfigAttrs {
	if len(self.Setting) == 0 {
		return NewMemberConfigAttrs()
	}
	var msa MemberConfigAttrs
	err := json.Unmarshal([]byte(self.Setting), &msa)
	if err != nil {
		return NewMemberConfigAttrs()
	}
	return &msa
}

type MemberConfigAttrs struct {
	AllowPush bool `json:"allow_push"`
}

func NewMemberConfigAttrs() *MemberConfigAttrs {
	//默认配置
	return &MemberConfigAttrs{
		AllowPush: true,
	}
}
