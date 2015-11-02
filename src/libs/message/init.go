package message

import (
	"fmt"
	"libs/message/service"
	"libs/stat"
	"libs/vars"

	"github.com/astaxie/beego"
	"github.com/thrift"
)

var use_ssdb_message_db, sys_msg_db, sys_msg_collection string //消息配置参数
var mbox_atmsg_length int
var sysMsgStorageConfig *MsgStorageConfig
var msgtype_storage_maps map[vars.MSG_TYPE]*MsgStorageConfig = make(map[vars.MSG_TYPE]*MsgStorageConfig)
var sys_message_service_run bool

const (
	SYS_MSG_COUNT_MODNAME       = "sys_msg_box_count_mod"
	SYS_MSG_NEWS_COUNT_MODENAME = "sys_msg_news_count_mod"
)

func init() {
	//初始化消息模块配置
	sys_msg_db = beego.AppConfig.String("sys.message.db")
	sys_msg_collection = beego.AppConfig.String("sys.msg.collection")
	use_ssdb_message_db = beego.AppConfig.String("ssdb.message.db")
	mbox_atmsg_length = beego.AppConfig.DefaultInt("mbox.atmsg.length", 200)
	sys_message_service_run = beego.AppConfig.DefaultBool("sys.message.service.run", false)
	sys_message_service_port := beego.AppConfig.DefaultInt("sys.message.service.port", 20001)
	initMsgSysConfig()

	if sys_message_service_run {
		go runMessageServer(sys_message_service_port)
	}

	//注册计数
	stat.RegisterUCountKey(SYS_MSG_COUNT_MODNAME, func(uid int64) string {
		return fmt.Sprintf("sys_msg_box_count:%d", uid)
	})
	stat.RegisterUCountKey(SYS_MSG_NEWS_COUNT_MODENAME, func(uid int64) string {
		return fmt.Sprintf("sys_msg_newalert:%d", uid)
	})
}

func runMessageServer(port int) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	host := fmt.Sprintf("0.0.0.0:%d", port)
	serverTransport, err := thrift.NewTServerSocket(host)
	if err != nil {
		panic(err)
	}

	handler := &MessageServiceImpl{}
	processor := service.NewMessageServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("message service server in " + host)
	server.Serve()
}

func initMsgSysConfig() {
	if len(sys_msg_db) == 0 {
		panic("未配置参数:sys.message.db参数")
	}
	if len(sys_msg_collection) == 0 {
		panic("未配置参数:sys.msg.collection参数")
	}
	sysMsgStorageConfig = &MsgStorageConfig{
		DbName:          sys_msg_db,
		TableName:       sys_msg_collection,
		CacheDb:         use_ssdb_message_db,
		MailboxSize:     mbox_atmsg_length,
		MailboxCountMod: SYS_MSG_COUNT_MODNAME,       //"sys_msg_box_count:%d",
		NewMsgCountMod:  SYS_MSG_NEWS_COUNT_MODENAME, //"sys_msg_newalert:%d",
	}

	RegisterMsgTypeMaps(vars.MSG_TYPE_SYS, sysMsgStorageConfig)
}

func RegisterMsgTypeMaps(msgType vars.MSG_TYPE, msc *MsgStorageConfig) {
	if _, ok := msgtype_storage_maps[msgType]; ok {
		return
	}
	msgtype_storage_maps[msgType] = msc
}

func getMsgStorageConfig(msgType vars.MSG_TYPE) *MsgStorageConfig {
	if v, ok := msgtype_storage_maps[msgType]; ok {
		return v
	}
	panic(fmt.Sprintf("不存在%s的消息存储配置", msgType))
}
