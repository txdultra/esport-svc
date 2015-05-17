package groups

import (
	"dbs"
	"libs/message"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var group_setting_id int
var update_group_limit_seconds int64
var db_aliasname, db_name string
var use_ssdb_group_db, credit_service_host string
var group_pic_thumbnail_w, group_pic_middle_w int
var search_server, group_timerjob_host string
var search_port, search_timeout int
var app_task_run bool
var use_ssdb_message_db, group_msg_db, group_msg_collection string //消息配置参数
var mbox_atmsg_length int
var msgStorageConfig *message.MsgStorageConfig

func init() {
	//ssdb tag
	use_ssdb_group_db = beego.AppConfig.String("group.ssdb.db")
	//积分系统地址
	credit_service_host = beego.AppConfig.String("group.credit.host")

	//setting id
	group_setting_id, _ = beego.AppConfig.Int("group.setting.id")
	update_group_limit_seconds = beego.AppConfig.DefaultInt64("group.update.limit.seconds", 20)

	group_pic_thumbnail_w, _ = beego.AppConfig.Int("group.pic.thumbnail.w")
	if group_pic_thumbnail_w <= 0 {
		group_pic_thumbnail_w = 74 * 2
	}
	group_pic_middle_w, _ = beego.AppConfig.Int("group.pic.middle.w")
	if group_pic_middle_w <= 0 {
		group_pic_middle_w = 148 * 2
	}

	app_task_run = beego.AppConfig.DefaultBool("app.task.run.group", false)
	group_timerjob_host = beego.AppConfig.String("group.timerjob.host")

	//search
	search_server = beego.AppConfig.String("search.group.server")
	search_port, _ = beego.AppConfig.Int("search.group.port")
	search_timeout, _ = beego.AppConfig.Int("search.group.timeout")

	//初始化消息模块配置
	group_msg_db = beego.AppConfig.String("group.message.db")
	group_msg_collection = beego.AppConfig.String("group.msg.collection")
	use_ssdb_message_db = beego.AppConfig.String("group.ssdb.db") //独立缓存库
	mbox_atmsg_length = beego.AppConfig.DefaultInt("mbox.atmsg.length", 200)
	initMsgSysConfig()

	orm.RegisterModel(new(GroupCfg), new(Group), new(Thread), new(Post),
		new(Report), new(MemberCount),
		new(GroupMemberTable), new(MemberGroupTable), new(PostTable))
	register_db()

	//启动前加载
	beego.AddAPPStartHook(func() error {
		load_tbls() //加载分表数据
		if app_task_run {
			//定时工作
			tjInit()
		}
		return nil
	})
}

func initMsgSysConfig() {
	if len(group_msg_db) == 0 {
		panic("未配置参数:group.message.db")
	}
	if len(group_msg_collection) == 0 {
		panic("未配置参数:group.msg.collection")
	}
	msgStorageConfig = &message.MsgStorageConfig{
		DbName:                group_msg_db,
		TableName:             group_msg_collection,
		CacheDb:               use_ssdb_message_db,
		MailboxSize:           mbox_atmsg_length,
		MailboxCountCacheName: "group_msg_box_count:%d",
		NewMsgCountCacheName:  "group_msg_newalert:%d",
	}

	message.RegisterMsgTypeMaps(MSG_TYPE_MESSAGE, msgStorageConfig)
	message.RegisterMsgTypeMaps(MSG_TYPE_INVITED, msgStorageConfig)
}

func register_db() {
	db_aliasname = beego.AppConfig.String("group.db.aliasname")
	if len(db_aliasname) == 0 {
		return
	}
	db_user := beego.AppConfig.String("group.db.user")
	db_pwd := beego.AppConfig.String("group.db.pwd")
	db_host := beego.AppConfig.String("group.db.host")
	db_port, _ := beego.AppConfig.Int("group.db.port")
	db_name = beego.AppConfig.String("group.db.name")
	db_charset := beego.AppConfig.String("group.db.charset")
	db_protocol := beego.AppConfig.String("group.db.protocol")
	db_time_local := beego.AppConfig.String("group.db.time_local")
	db_maxconns, _ := beego.AppConfig.Int("group.db.maxconns")
	db_maxidels, _ := beego.AppConfig.Int("group.db.maxidles")
	if db_maxconns <= 0 {
		db_maxconns = 1500
	}
	if db_maxidels <= 0 {
		db_maxidels = 100
	}

	db_addr := db_host + ":" + strconv.Itoa(db_port)
	connection_url := db_user + ":" + db_pwd + "@" + db_protocol + "(" + db_addr + ")/" + db_name + "?charset=" + db_charset + "&loc=" + db_time_local
	dbs.LoadDb(db_aliasname, "mysql", connection_url, db_maxidels, db_maxconns)
}

func load_tbls() {
	tbl_mutex.Lock()
	defer tbl_mutex.Unlock()

	o := dbs.NewOrm(db_aliasname)

	var gmtbls []*GroupMemberTable
	o.QueryTable(&GroupMemberTable{}).All(&gmtbls)
	for _, tbl := range gmtbls {
		gmTbls[int(tbl.Id)] = tbl.TblName
	}

	var mgtbls []*MemberGroupTable
	o.QueryTable(&MemberGroupTable{}).All(&mgtbls)
	for _, tbl := range mgtbls {
		mgTbls[int(tbl.Id)] = tbl.TblName
	}

	var ptbls []*PostTable
	o.QueryTable(&PostTable{}).All(&ptbls)
	for _, tbl := range ptbls {
		postTbls[int(tbl.Id)] = tbl.TblName
	}
}
