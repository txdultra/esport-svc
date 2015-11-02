package share

import (
	"dbs"
	"fmt"
	"libs/message"
	"libs/passport"
	"libs/stat"
	"strconv"

	"github.com/astaxie/beego"
)

var sns_share_db, sns_share_collection, use_ssdb_share_db, use_ssdb_cmt_db, use_ssdb_notice_db string
var use_ssdb_message_db, sns_msg_db, sns_msg_collection string //消息配置参数
var mbox_share_length, mbox_subscription_length, mbox_atmsg_length int
var sns_atcount_cache_length int
var sns_share_pic_thumbnail_w, sns_share_pic_thumbnail_h, sns_share_pic_middle_w, sns_share_pic_middle_h int
var sns_notice_db_insert_queue_open, sns_msg_db_insert_queue_open bool
var msq_db_batch_queue_name string
var msgStorageConfig *message.MsgStorageConfig

const (
	share_db  = "user_share"
	share_tbl = "user_share"

	//subcur_db  = "user_subcur"
	//subcur_tbl = "user_subcur"

	msg_db  = "user_msg"
	msg_tbl = "user_msg"

	share_cmt_db  = "user_share_cmt"
	share_cmt_tbl = "share_comment"

	share_notice_db  = "user_share_notice"
	share_notice_tbl = "share_notice"
)

func init() {
	//orm.RegisterModel(
	//	new(ShareViewPicture),
	//)
	mbox_share_length, _ = beego.AppConfig.Int("mbox.share.length")
	sns_atcount_cache_length, _ = beego.AppConfig.Int("sns.atcount.cache.length")
	if mbox_share_length <= 0 {
		mbox_share_length = 100
	}
	if sns_atcount_cache_length <= 0 {
		sns_atcount_cache_length = 10
	}
	mbox_subscription_length, _ = beego.AppConfig.Int("mbox.subscription.length")
	if mbox_subscription_length <= 0 {
		mbox_subscription_length = 100
	}

	sns_share_pic_thumbnail_w, _ = beego.AppConfig.Int("sns.share.pic.thumbnail.w")
	if sns_share_pic_thumbnail_w <= 0 {
		sns_share_pic_thumbnail_w = 74 * 2
	}
	sns_share_pic_middle_w, _ = beego.AppConfig.Int("sns.share.pic.middle.w")
	if sns_share_pic_middle_w <= 0 {
		sns_share_pic_middle_w = 148 * 2
	}

	use_ssdb_share_db = beego.AppConfig.String("ssdb.share.db")
	use_ssdb_cmt_db = beego.AppConfig.String("ssdb.share.cmt.db")
	use_ssdb_notice_db = beego.AppConfig.String("ssdb.share.notice.db")
	sns_notice_db_insert_queue_open, _ = beego.AppConfig.Bool("sns.notice.db_insert_queue.open")
	sns_msg_db_insert_queue_open, _ = beego.AppConfig.Bool("sns.msg.db_insert_queue.open")
	msq_db_batch_queue_name = beego.AppConfig.String("msq.db_batch.queue_name")
	if len(use_ssdb_share_db) == 0 {
		panic("未配置参数:ssdb.share.db参数")
	}
	if len(use_ssdb_cmt_db) == 0 {
		panic("未配置参数:ssdb.share.cmt.db参数")
	}
	if len(use_ssdb_notice_db) == 0 {
		panic("未配置参数:ssdb.share.notice.db参数")
	}
	if len(msq_db_batch_queue_name) == 0 {
		panic("未配置参数:msq.db_batch.queue_name参数")
	}

	//注册关注事件
	passport.RegisterFriendEvent("friend_share_msg_event", &ShareMsgs{})
	passport.RegisterFriendEvent("friend_share_subcur_event", &ShareVodSubcurs{})
	passport.RegisterUnFriendEvent("unfriend_share_msg_event", &ShareMsgs{})
	passport.RegisterUnFriendEvent("unfriend_share_subcur_event", &ShareVodSubcurs{})

	//注册评论事件
	RegisterShareCommentEvent("share_commented_notice_event", &ShareNotices{})

	//注册用户计数器
	stat.RegisterUCountKey(SHARE_MY_VOD_SUBSCRIPTIONS_COUNT_MODNAME, func(uid int64) string {
		return fmt.Sprintf("share_my_subscriptions_counts:%d", uid)
	})
	stat.RegisterUCountKey(SHARE_MY_NEW_NOTICES_COUNT_MODNAME, func(uid int64) string {
		return fmt.Sprintf("share_new_notices_counts:%d", uid)
	})
	stat.RegisterUCountKey(SHARE_LASTNEW_MSG_MODNAME, func(uid int64) string {
		return fmt.Sprintf("share_member_newmsg_counts:%d", uid)
	})
	stat.RegisterUCountKey(SHARE_MEMBER_ATMSG_BOX_COUNT_MODNAME, func(uid int64) string {
		return fmt.Sprintf("member_atmsg_box_count:%d", uid)
	})
	stat.RegisterUCountKey(SHARE_MEMBER_ATMSG_NEW_COUNT_MODNAME, func(uid int64) string {
		return fmt.Sprintf("member_atmsg_newalert:%d", uid)
	})

	//db
	register_share_db()
	//register_subsur_db()
	register_msg_db()
	register_share_cmt_db()
	register_share_notice_db()

	//初始化分享中的资源处理方法
	initShareKindFuncs()

	//初始化消息模块配置
	sns_msg_db = beego.AppConfig.String("sns.share.db")
	sns_msg_collection = beego.AppConfig.String("sns.atmsg.collection")
	use_ssdb_message_db = beego.AppConfig.String("ssdb.message.db")
	mbox_atmsg_length = beego.AppConfig.DefaultInt("mbox.atmsg.length", 200)
	if len(sns_msg_db) == 0 {
		panic("未配置参数:sns.share.db参数")
	}
	if len(sns_msg_collection) == 0 {
		panic("未配置参数:sns.atmsg.collection参数")
	}
	if len(use_ssdb_message_db) == 0 {
		panic("未配置参数:ssdb.message.db参数")
	}
	initMsgSysConfig()

	//register share events
	RegisterShareNotifyEvent("share_msg", &ShareMsgs{})
	RegisterShareNotifyEvent("share_vod_subcur", &ShareVodSubcurs{})

	//register share save event
	RegisterShareSaveEvent("share_msg_save", &ShareMsgs{})
	RegisterShareSaveEvent("share_notice_event", &ShareNotices{})
	RegisterShareRevokeEvent("share_msg_revoke", &ShareMsgs{})
}

func initMsgSysConfig() {
	if len(sns_msg_db) == 0 {
		panic("未配置参数:sns.share.db")
	}
	if len(sns_msg_collection) == 0 {
		panic("未配置参数:sns.atmsg.collection")
	}
	msgStorageConfig = &message.MsgStorageConfig{
		DbName:          sns_msg_db,
		TableName:       sns_msg_collection,
		CacheDb:         use_ssdb_message_db,
		MailboxSize:     mbox_atmsg_length,
		MailboxCountMod: SHARE_MEMBER_ATMSG_BOX_COUNT_MODNAME,
		NewMsgCountMod:  SHARE_MEMBER_ATMSG_NEW_COUNT_MODNAME,
	}

	message.RegisterMsgTypeMaps(MSG_TYPE_VOD, msgStorageConfig)
	message.RegisterMsgTypeMaps(MSG_TYPE_TEXT, msgStorageConfig)
	message.RegisterMsgTypeMaps(MSG_TYPE_PICS, msgStorageConfig)
}

func GetMsgConfig() *message.MsgStorageConfig {
	return msgStorageConfig
}

func register_share_notice_db() {
	db_user := beego.AppConfig.String("db.share.notice.user")
	db_pwd := beego.AppConfig.String("db.share.notice.pwd")
	db_host := beego.AppConfig.String("db.share.notice.host")
	db_port, _ := beego.AppConfig.Int("db.share.notice.port")
	db_name := beego.AppConfig.String("db.share.notice.name")
	db_charset := beego.AppConfig.String("db.share.notice.charset")
	db_protocol := beego.AppConfig.String("db.share.notice.protocol")
	db_time_local := beego.AppConfig.String("db.share.notice.time_local")
	db_maxconns, _ := beego.AppConfig.Int("db.share.notice.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("db.share.notice.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(share_notice_db, "mysql", connection_url, db_maxidels, db_maxconns)
}

func register_share_cmt_db() {
	db_user := beego.AppConfig.String("db.share.cmt.user")
	db_pwd := beego.AppConfig.String("db.share.cmt.pwd")
	db_host := beego.AppConfig.String("db.share.cmt.host")
	db_port, _ := beego.AppConfig.Int("db.share.cmt.port")
	db_name := beego.AppConfig.String("db.share.cmt.name")
	db_charset := beego.AppConfig.String("db.share.cmt.charset")
	db_protocol := beego.AppConfig.String("db.share.cmt.protocol")
	db_time_local := beego.AppConfig.String("db.share.cmt.time_local")
	db_maxconns, _ := beego.AppConfig.Int("db.share.cmt.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("db.share.cmt.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(share_cmt_db, "mysql", connection_url, db_maxidels, db_maxconns)
}

func register_msg_db() {
	db_user := beego.AppConfig.String("db.share.msg.user")
	db_pwd := beego.AppConfig.String("db.share.msg.pwd")
	db_host := beego.AppConfig.String("db.share.msg.host")
	db_port, _ := beego.AppConfig.Int("db.share.msg.port")
	db_name := beego.AppConfig.String("db.share.msg.name")
	db_charset := beego.AppConfig.String("db.share.msg.charset")
	db_protocol := beego.AppConfig.String("db.share.msg.protocol")
	db_time_local := beego.AppConfig.String("db.share.msg.time_local")
	db_maxconns, _ := beego.AppConfig.Int("db.share.msg.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("db.share.msg.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(msg_db, "mysql", connection_url, db_maxidels, db_maxconns)
}

func register_share_db() {
	db_user := beego.AppConfig.String("db.share.user")
	db_pwd := beego.AppConfig.String("db.share.pwd")
	db_host := beego.AppConfig.String("db.share.host")
	db_port, _ := beego.AppConfig.Int("db.share.port")
	db_name := beego.AppConfig.String("db.share.name")
	db_charset := beego.AppConfig.String("db.share.charset")
	db_protocol := beego.AppConfig.String("db.share.protocol")
	db_time_local := beego.AppConfig.String("db.share.time_local")
	db_maxconns, _ := beego.AppConfig.Int("db.share.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("db.share.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(share_db, "mysql", connection_url, db_maxidels, db_maxconns)
}

//func register_subsur_db() {
//	db_user := beego.AppConfig.String("db.subsur.user")
//	db_pwd := beego.AppConfig.String("db.subsur.pwd")
//	db_host := beego.AppConfig.String("db.subsur.host")
//	db_port, _ := beego.AppConfig.Int("db.subsur.port")
//	db_name := beego.AppConfig.String("db.subsur.name")
//	db_charset := beego.AppConfig.String("db.subsur.charset")
//	db_protocol := beego.AppConfig.String("db.subsur.protocol")
//	db_time_local := beego.AppConfig.String("db.subsur.time_local")
//	db_maxconns, _ := beego.AppConfig.Int("db.subsur.maxconns")
//	db_maxidels, _ := beego.AppConfig.Int("db.subsur.maxidles")
//	if db_maxconns <= 0 {
//		db_maxconns = 1500
//	}
//	if db_maxidels <= 0 {
//		db_maxidels = 100
//	}

//	db_addr := db_host + ":" + strconv.Itoa(db_port)
//	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
//	dbs.LoadDb(subcur_db, "mysql", connection_url, db_maxidels, db_maxconns)
//}
