package message

import (
	"libs/vars"
	"time"

	"labix.org/v2/mgo/bson"
)

type MsgStorageConfig struct {
	DbName                string
	TableName             string
	CacheDb               string
	MailboxSize           int
	MailboxCountCacheName string
	NewMsgCountCacheName  string
}

type MsgData struct {
	Id        bson.ObjectId `bson:"_id"`
	FromUid   int64         `bson:"from_uid"`
	ToUid     int64         `bson:"to_uid"`
	MsgType   vars.MSG_TYPE `bson:"msg_type"`
	Text      string        `bson:"text"`
	RefId     string        `bson:"ref_id"` //关联id,如视频id等
	PostTime  time.Time     `bson:"post_time"`
	Timestamp int64         `bson:"timstamp"`
}
