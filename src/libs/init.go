package libs

import (
	"logs"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//"strings"
	//"reflect"
	//"fmt"

	"sync"
	"time"
	"utils"
)

var file_provider_name, msq_connection_url string
var file_connect_timeout, file_rwconnect_timeout time.Duration
var open_distributed, app_task_run_msq bool
var once sync.Once
var msg_process_task map[string]IMsqMsgProcessTasker = make(map[string]IMsqMsgProcessTasker)

func init() {
	orm.RegisterModel(
		new(Game), new(Match), new(File), new(Msqtor), new(Smiley),
		//组件
		new(PushState), new(Recommend),
	)
	MSQ_USE_DRIVER = beego.AppConfig.String("msq.driver")
	//RegisterMsqProcessHandler("msq_test", reflect.TypeOf(TestHandlerMsg{}))

	//注册groupcache策略
	//initToGroupCache()

	file_provider_name = beego.AppConfig.String("file.provider.name")
	file_connect_timeout = utils.StrToDuration(beego.AppConfig.String("file.connect.timeout"))
	file_rwconnect_timeout = utils.StrToDuration(beego.AppConfig.String("file.rwconnect.timeout"))

	open_distributed, _ = beego.AppConfig.Bool("app.distributed")
	msq_connection_url = beego.AppConfig.String("msq.rabbitmq.addr")

	//task
	app_task_run_msq, _ = beego.AppConfig.Bool("app.task.run.msq")

	beego.AddAPPStartHook(func() error {
		once.Do(func() {
			if app_task_run_msq { //开启后台任务
				startMsqMsgProcessTaskers()
			}
		})
		return nil
	})
}

func OpenDistributed() bool {
	return open_distributed
}

////////////////////////////////////////////////////////////////////////////////
//队列功能
func MsqConnectionUrl() string {
	return msq_connection_url
}

func RegisterMsqMsgrocessTasker(queueName string, tasker IMsqMsgProcessTasker) {
	if _, ok := msg_process_task[queueName]; !ok {
		msg_process_task[queueName] = tasker
	}
}

func startMsqMsgProcessTaskers() {
	for _, tasker := range msg_process_task {
		go runMsqMsgProcessTasker(tasker)
	}
}

func runMsqMsgProcessTasker(tasker IMsqMsgProcessTasker) {
	for {
		c, err := tasker.Run()
		if err != nil {
			logs.Errorf("启动后台消息处理线程失败:%v", err)
			time.Sleep(1 * time.Minute)
			continue
		}
		for {
			msg := <-c
			if msg == MSQ_RECEIVE_CLOSED {
				time.Sleep(10 * time.Second) //停10秒后重启
				break
			}
		}
	}
}
