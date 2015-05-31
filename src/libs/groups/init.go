package groups

import (
	"dbs"
	"fmt"
	"libs/dlock"
	"libs/message"
	"libs/share"
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
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

//分布式监视者
var watcher = dlock.NewWatcher()

const (
	watcher_path = "/group_config"
)

func init() {
	//ssdb tag
	use_ssdb_group_db = beego.AppConfig.String("group.ssdb.db")
	if len(use_ssdb_group_db) == 0 {
		panic("未配置group.ssdb.db参数")
	}
	//积分系统地址
	credit_service_host = beego.AppConfig.String("group.credit.host")
	if len(credit_service_host) == 0 {
		panic("未配置group.credit.host参数")
	}

	//setting id
	group_setting_id = beego.AppConfig.DefaultInt("group.setting.id", 1)
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
	if len(group_timerjob_host) == 0 {
		panic("未配置group.timerjob.host参数")
	}

	//search
	search_server = beego.AppConfig.String("search.group.server")
	search_port, _ = beego.AppConfig.Int("search.group.port")
	search_timeout, _ = beego.AppConfig.Int("search.group.timeout")
	if len(search_server) == 0 {
		panic("未配置小组使用的搜索服务地址")
	}

	orm.RegisterModel(new(GroupCfg), new(Group), new(Thread), new(Post),
		new(Report), new(MemberCount),
		new(GroupMemberTable), new(MemberGroupTable), new(PostTable))
	register_db()

	//初始化消息模块配置
	group_msg_db = beego.AppConfig.String("group.message.db")
	group_msg_collection = beego.AppConfig.String("group.msg.collection")
	use_ssdb_message_db = beego.AppConfig.String("group.ssdb.db") //独立缓存库
	mbox_atmsg_length = beego.AppConfig.DefaultInt("mbox.atmsg.length", 200)
	if len(group_msg_db) == 0 {
		panic("未配置group.message.db参数")
	}
	if len(group_msg_collection) == 0 {
		panic("未配置group.msg.collection参数")
	}
	if len(use_ssdb_message_db) == 0 {
		panic("未配置group.ssdb.db参数")
	}

	initMsgSysConfig()

	//启动前加载
	beego.AddAPPStartHook(func() error {
		//加载分表数据
		load_tbls()
		//注册分享模块功能
		register_share_mod()

		if app_task_run {
			//定时工作
			tjInit()
			runGroupCountUpdateService()
			//注册监视者通知
			watcher.RegisterWatcher(watcher_path, func(data []byte) {
				ResetDefaultCfg()
			})
			//活跃度统计
			tskName := "group_vitality_task"
			spec := "0 0 0/7 * * *"
			_spec := beego.AppConfig.String("group.vitality.interval.spec")
			if len(_spec) > 0 {
				spec = _spec
			}
			toolbox.AddTask(tskName, toolbox.NewTask(tskName, spec, func() error {
				gpVitality()
				return nil
			}))
		}
		return nil
	})
}

func register_share_mod() {
	ts := NewThreadService(GetDefaultCfg())
	//注册资源转换方法s
	share.RegisterShareTransformResourceFuncs(SHARE_KIND_GROUP, func(id string, args ...string) string {
		return fmt.Sprintf("[gpt:%s]", id)
	})
	share.RegisterShareResourceTransformFuncs(SHARE_KIND_GROUP, func(resource string) (share.SHARE_KIND, string, error) {
		if ok, _ := regexp.MatchString(`\[gpt:(\d+)\]`, resource); ok {
			rep := regexp.MustCompile(`\[gpt:(\d+)\]`)
			arr := rep.FindStringSubmatch(resource)
			return SHARE_KIND_GROUP, arr[1], nil
		}
		return share.SHARE_KIND_EMPTY, "", fmt.Errorf("未匹配GROUP_KIND")
	})
	share.RegisterResPicFileIdFuncs(SHARE_KIND_GROUP, func(res *share.ShareResource) int64 {
		_id, _ := strconv.ParseInt(res.Id, 10, 64)
		thread := ts.Get(_id)
		if thread == nil {
			return 0
		}
		return thread.Img
	})
	share.RegisterResToOutputProxyObjectFuncs(SHARE_KIND_GROUP, func(res *share.ShareResource) *share.ResOutputProxyObject {
		_id, _ := strconv.ParseInt(res.Id, 10, 64)
		thread := ts.Get(_id)
		if thread == nil {
			return nil
		}
		return &share.ResOutputProxyObject{
			Id:           res.Id,
			Title:        thread.Subject,
			Content:      "",
			ThumbnailPic: thread.Img,
			Uid:          thread.AuthorId,
		}
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
	message.RegisterMsgTypeMaps(MSG_TYPE_REPLY, msgStorageConfig)
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
