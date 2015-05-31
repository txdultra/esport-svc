package share

import (
	"libs/vars"
	"strconv"
	"strings"
	//"labix.org/v2/mgo/bson"

	"logs"
	"time"
)

const (
	MSG_TYPE_VOD  vars.MSG_TYPE = "share:vod"
	MSG_TYPE_TEXT vars.MSG_TYPE = "share:text"
	MSG_TYPE_PICS vars.MSG_TYPE = "share:pics"
)

type SHARE_KIND int

const (
	SHARE_KIND_EMPTY SHARE_KIND = 0
	SHARE_KIND_TXT   SHARE_KIND = 1
	SHARE_KIND_VOD   SHARE_KIND = 2
	SHARE_KIND_PIC   SHARE_KIND = 4
	SHARE_KIND_URL   SHARE_KIND = 8
)

type SHARE_TYPE int8

const (
	SHARE_TYPE_ORIGINAL SHARE_TYPE = 0
	SHARE_TYPE_TRANSFER SHARE_TYPE = 1
	SHARE_TYPE_COMMENT  SHARE_TYPE = 2
)

var NeedSharePicSizes []vars.PIC_SIZE = []vars.PIC_SIZE{
	vars.PIC_SIZE_ORIGINAL,
	vars.PIC_SIZE_THUMBNAIL,
	vars.PIC_SIZE_MIDDLE,
}

type ShareSubcur struct {
	Uid  int64
	Sid  int64
	SUid int64
	Ts   int64
	St   int
}

type Share struct {
	Id               int64
	Uid              int64
	ShareType        int        `orm:"column(st)"`
	Type             SHARE_TYPE `orm:"column(t)"`
	Source           string     `orm:"column(source)"`
	Geo              string
	Text             string `orm:"column(txt)"`
	CreateTime       time.Time
	Ts               int64
	Resources        string `orm:"column(res)"`
	CommentedCount   int    `orm:"column(cmted_count)"`
	CommentCount     int    `orm:"column(cmt_count)"`
	TransferredCount int    `orm:"column(transferred_count)"`
	TransferCount    int    `orm:"column(transfer_count)"`
	AttitudedCount   int    `orm:"column(attituded_count)"`
	RefUids          string `orm:"column(ref_uids)"`
}

func (self *Share) GetRefUids() []int64 {
	uid_maps := make(map[int64]string)
	uids := []int64{}
	if len(self.RefUids) == 0 {
		return uids
	}
	uidarr := strings.Split(self.RefUids, ",")
	for _, str := range uidarr {
		_uid, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			logs.Errorf("share method getRefUids transfor uid fail:%s", err.Error())
			continue
		}
		if _uid <= 0 {
			continue
		}
		if _, ok := uid_maps[_uid]; ok {
			continue
		}
		uid_maps[_uid] = str
		uids = append(uids, _uid)
	}
	return uids
}

func (self *Share) TableName() string {
	return "share_threads"
}

func (self *Share) TableEngine() string {
	return "INNODB"
}

type ShareResource struct {
	Id   string
	Kind SHARE_KIND
}

type ShareViewPicture struct {
	ParentFileId int64 `orm:"column(pfid)"`
	FileId       int64 `orm:"column(fid)"`
	Ts           int64
	PicSize      vars.PIC_SIZE `orm:"column(ps)"`
}

func (self *ShareViewPicture) TableName() string {
	return "share_view_pics"
}

func (self *ShareViewPicture) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (self *ShareViewPicture) TableUnique() [][]string {
	return [][]string{
		[]string{"pfid", "ps"},
	}
}

type SHARE_COMMENT_TYPE int8

const (
	SHARE_COMMENT_TYPE_SINGLE SHARE_COMMENT_TYPE = 0
	SHARE_COMMENT_TYPE_REPLY  SHARE_COMMENT_TYPE = 1
)

type ShareComment struct {
	Id         int64              `orm:"column(id)"`
	Sid        int64              `orm:"column(sid)"`
	SUid       int64              `orm:"column(suid)"`
	SUNickname string             `orm:"column(su_nickname)"`
	Uid        int64              `orm:"column(uid)"`
	UNickname  string             `orm:"column(u_nickname)"`
	RUid       int64              `orm:"column(ruid)"` //回复对象uid
	RUNickname string             `orm:"column(ru_nickname)"`
	T          SHARE_COMMENT_TYPE `orm:"column(t)"`
	Content    string             `orm:"column(content)"`
	Ts         int64              `orm:"column(ts)"`
	Ex1        string             `orm:"column(ex1)"`
	Ex2        string             `orm:"column(ex2)"`
	Ex3        string             `orm:"column(ex3)"`
}

type SHARE_NOTICE_TYPE int8

const (
	SHARE_NOTICE_TYPE_TD SHARE_NOTICE_TYPE = 0 //提到了你
	SHARE_NOTICE_TYPE_PL SHARE_NOTICE_TYPE = 1 //评论
	SHARE_NOTICE_TYPE_HF SHARE_NOTICE_TYPE = 2 //回复
)

type ShareNotice struct {
	Id         string            `orm:"column(id)"`
	Sid        int64             `orm:"column(sid)"`
	Uid        int64             `orm:"column(uid)"`
	LUid       int64             `orm:"column(luid)"`
	LUNickname string            `orm:"column(lu_nickname)"`
	RUid       int64             `orm:"column(ruid)"`
	RUNickname string            `orm:"column(ru_nickname)"`
	Content    string            `orm:"column(content)"`
	Pic        int64             `orm:"column(pic)"`
	Ts         int64             `orm:"column(ts)"`
	T          SHARE_NOTICE_TYPE `orm:"column(t)"`
	ST         int               `orm:"column(st)"`
}

//发送到消息队列使用
type dbMsg struct {
	DbName string `json:"db_name"`
	Sql    string `json:"sql"`
}

//RefId      string        `bson:"ref_id"`
//ExtIds     []string      `bson:"ext_ids"`

//type SimpleShare struct {
//	RelId     string
//	ShareType int
//}
