package message

import (
	"github.com/astaxie/beego"
)

var mbox_atmsg_length int
var sns_atmsg_db, sns_atmsg_collection, use_ssdb_message_db string

func init() {
	mbox_atmsg_length, _ = beego.AppConfig.Int("mbox.atmsg.length")
	if mbox_atmsg_length <= 0 {
		mbox_atmsg_length = 200
	}
	sns_atmsg_db = beego.AppConfig.String("sns.share.db")
	sns_atmsg_collection = beego.AppConfig.String("sns.atmsg.collection")
	use_ssdb_message_db = beego.AppConfig.String("ssdb.message.db")
}
