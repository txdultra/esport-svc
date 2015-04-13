package comment

import (
	"labix.org/v2/mgo/bson"
	"time"
)

//type AtUser struct {
//	ScreenName string `bson:"screen_name"`
//	Uid        int64  `bson:"uid"`
//}

type Comment struct {
	ID         bson.ObjectId    `bson:"_id"`
	Pid        string           `bson:"-"`
	RefId      int64            `bson:"ref_id"`
	FromId     int64            `bson:"from_id"`
	Position   int              `bson:"position"`
	Title      string           `bson:"title"`
	Text       string           `bson:"text"`
	WbText     string           `bson:"wb_text"`
	Ats        map[string]int64 `bson:"ats"`
	Reply      *ReplyComment    `bson:"reply"`
	PostTime   time.Time        `bson:"post_time"`
	ModifyTime time.Time        `bson:"modify_time"`
	AllowReply bool             `bson:"allow_reply"`
	InVisible  bool             `bson:"invisible"`
	IP         string           `bson:"ip"`
	Area       string           `bson:"area"`
	Longitude  float64          `bson:"longitude"`
	Latitude   float64          `bson:"latitude"`
	Nano       int64            `bson:"nano"`
}

type ReplyComment struct {
	ID        string           `bson:"id"`
	FromId    int64            `bson:"from_id"`
	Position  int              `bson:"position"`
	Text      string           `bson:"text"`
	WbText    string           `bson:"wb_text"`
	Ats       map[string]int64 `bson:"ats"`
	PostTime  time.Time        `bson:"post_time"`
	IP        string           `bson:"ip"`
	Longitude float64          `bson:"longitude"`
	Latitude  float64          `bson:"latitude"`
}
