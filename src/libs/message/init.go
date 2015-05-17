package message

import (
	"fmt"

	"github.com/astaxie/beego"
)

var use_ssdb_message_db, sys_msg_db, sys_msg_collection string //消息配置参数
var mbox_atmsg_length int
var sysMsgStorageConfig *MsgStorageConfig
var msgtype_storage_maps map[MSG_TYPE]*MsgStorageConfig = make(map[MSG_TYPE]*MsgStorageConfig)

func init() {
	//初始化消息模块配置
	sys_msg_db = beego.AppConfig.String("sys.message.db")
	sys_msg_collection = beego.AppConfig.String("sys.msg.collection")
	use_ssdb_message_db = beego.AppConfig.String("ssdb.message.db")
	mbox_atmsg_length = beego.AppConfig.DefaultInt("mbox.atmsg.length", 200)
	initMsgSysConfig()
}

func initMsgSysConfig() {
	if len(sys_msg_db) == 0 {
		panic("未配置参数:sys.message.db")
	}
	if len(sys_msg_collection) == 0 {
		panic("未配置参数:sys.msg.collection")
	}
	sysMsgStorageConfig = &MsgStorageConfig{
		DbName:                sys_msg_db,
		TableName:             sys_msg_collection,
		CacheDb:               use_ssdb_message_db,
		MailboxSize:           mbox_atmsg_length,
		MailboxCountCacheName: "sys_msg_box_count:%d",
		NewMsgCountCacheName:  "sys_msg_newalert:%d",
	}

	RegisterMsgTypeMaps(MSG_TYPE_SYS, sysMsgStorageConfig)
}

func RegisterMsgTypeMaps(msgType MSG_TYPE, msc *MsgStorageConfig) {
	if _, ok := msgtype_storage_maps[msgType]; ok {
		return
	}
	msgtype_storage_maps[msgType] = msc
}

func getMsgStorageConfig(msgType MSG_TYPE) *MsgStorageConfig {
	if v, ok := msgtype_storage_maps[msgType]; ok {
		return v
	}
	panic(fmt.Sprintf("不存在%s的消息存储配置", msgType))
}
