package message

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type MSG_TYPE string

const (
	MSG_TYPE_COMMENT MSG_TYPE = "vod:comment"
	MSG_TYPE_VOD     MSG_TYPE = "share:vod"
	MSG_TYPE_TEXT    MSG_TYPE = "share:text"
	MSG_TYPE_PICS    MSG_TYPE = "share:pics"
)

type MsgData struct {
	Id        bson.ObjectId `bson:"_id"`
	FromUid   int64         `bson:"from_uid"`
	ToUid     int64         `bson:"to_uid"`
	MsgType   MSG_TYPE      `bson:"msg_type"`
	Text      string        `bson:"text"`
	RefId     string        `bson:"ref_id"` //关联id,如视频id等
	PostTime  time.Time     `bson:"post_time"`
	Timestamp int64         `bson:"timstamp"`
}

//type MsgTypeConvert interface {
//	MTConvert(obj interface{}) MSG_TYPE
//}

//var convertors map[string]MsgTypeConvert = make(map[string]MsgTypeConvert)

//func RegisterConverter(mod_name string, convertor MsgTypeConvert) {
//	if ok, _ := convertors[mod_name]; ok {
//		return
//	}
//	converters[mod_name] = converter
//}
