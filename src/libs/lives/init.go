package lives

import (
	"fmt"
	"libs"
	"libs/collect"
	"sync"
	"time"
	"utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var once sync.Once
var search_live_server, search_program_server string
var search_live_timeout, search_program_timeout, search_live_port, search_program_port int
var event_program_key, use_ssdb_live_db string
var program_lock_dur_str time.Duration = utils.StrToDuration("300s")
var app_task_run_live bool
var push_baidu_apikey, push_baidu_secret string

const (
	COLLECT_ORG_MOD = "j_live"
	COLLECT_PER_MOD = "p_live"
)

func init() {
	orm.RegisterModel(
		//直播
		new(LivePerson), new(LiveChannel), new(LiveStream), new(LiveProgram), new(LiveSubProgram),
		new(LiveSearchData), new(LiveProgramSearchData),
	)
	program_notice_db = beego.AppConfig.String("notice.db")
	program_notice_collection = beego.AppConfig.String("notice.program.collection")
	program_notice_orginal_collection = beego.AppConfig.String("notice.program.orginal.collection")
	search_live_server = beego.AppConfig.String("search.live.server")
	search_live_port, _ = beego.AppConfig.Int("search.live.port")
	search_live_timeout, _ = beego.AppConfig.Int("search.live.timeout")
	search_program_server = beego.AppConfig.String("search.program.server")
	search_program_port, _ = beego.AppConfig.Int("search.program.port")
	search_program_timeout, _ = beego.AppConfig.Int("search.program.timeout")
	//baidu 云推送
	push_baidu_apikey = beego.AppConfig.String("push.baidu.apikey")
	push_baidu_secret = beego.AppConfig.String("push.baidu.secret")
	//
	event_program_key = beego.AppConfig.String("push.class.program")
	use_ssdb_live_db = beego.AppConfig.String("ssdb.live.db")
	//锁定时间段
	_lock_dur_str := beego.AppConfig.String("update.locker.program.time")
	if _lock_dur_str != "" {
		program_lock_dur_str = utils.StrToDuration(_lock_dur_str)
	}
	//task
	app_task_run_live, _ = beego.AppConfig.Bool("app.task.run.live")

	//检验参数
	if len(program_notice_db) == 0 {
		panic("未配置notice.db参数")
	}
	if len(program_notice_collection) == 0 {
		panic("未配置notice.program.collection参数")
	}
	if len(search_live_server) == 0 {
		panic("未配置search.live.server参数")
	}
	if len(push_baidu_apikey) == 0 || len(push_baidu_secret) == 0 {
		panic("未配置push.baidu.apikey或push.baidu.secret参数")
	}
	if len(event_program_key) == 0 {
		panic("未配置push.class.program参数")
	}
	if len(use_ssdb_live_db) == 0 {
		panic("未配置ssdb.live.db参数")
	}

	//启动前加载
	beego.AddAPPStartHook(func() error {
		once.Do(func() {
			if app_task_run_live { //开启后台任务
				fmt.Println("program notice init...")
				pnotice := &ProgramNoticer{}
				pnotice.loadNoticeTimers()
				fmt.Println("live personal init...")
				live := &LivePers{}
				live.init()
			}
		})
		return nil
	})

	//注册收藏功能
	collect.RegisterCompler(COLLECT_ORG_MOD, &LiveOrgs{})
	collect.RegisterCompler(COLLECT_PER_MOD, &LivePers{})

	//注册推送通知后台处理线程
	libs.RegisterMsqMsgrocessTasker("live_program_push_processer", &ProgramNoticeMsqProcesser{})
}
