package vod

import (
	"fmt"
	"libs/collect"
	"libs/stat"
	"strconv"
	"sync"
	"time"
	"utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

var VOD_RE_REPTILE_TIMEOUT time.Duration                             //视频地址抓取过期时间
var reptileTasks map[int64]*time.Timer = make(map[int64]*time.Timer) //过期抓取任务
var taskLocker *sync.Mutex = new(sync.Mutex)

var issusLocker *sync.Mutex = new(sync.Mutex)
var repChans map[int64]chan string = make(map[int64]chan string)
var once sync.Once

var search_vod_server, use_ssdb_vod_db string
var search_vod_port, search_vod_timeout int
var app_task_run_vod bool

const (
	MOD_NAME                     string = "vod"
	VOD_FLVS_COLLECTION_NAME     string = "flvs"
	VOD_RECOMMEND_CATEGORTY_NAME string = "vod_home"

	collect_vod_mod = "vod"
)

func init() {
	orm.RegisterModel(
		//视频
		new(Video), //new(VideoOpt), new(VideoFlv),
		new(VideoPlaylist), new(VideoPlaylistVod),
		new(VodUcenter), new(VideoCount),
	)
	///////////////////////////////////////////////////////////
	rp_timeout := "3h" //默认超时时间
	rp_tmp := beego.AppConfig.String("vod.reptile.timeout")
	if len(rp_tmp) > 0 {
		rp_timeout = rp_tmp
	}
	VOD_RE_REPTILE_TIMEOUT = utils.StrToDuration(rp_timeout)
	///////////////////////////////////////////////////////////

	//counter seed
	vod_view_counter_seed, _ = beego.AppConfig.Int("vod.views.counter.seed")
	if vod_view_counter_seed <= 0 {
		vod_view_counter_seed = 300 //默认300
	}

	//search
	search_vod_server = beego.AppConfig.String("search.vod.server")
	search_vod_port, _ = beego.AppConfig.Int("search.vod.port")
	search_vod_timeout, _ = beego.AppConfig.Int("search.vod.timeout")

	//task
	app_task_run_vod, _ = beego.AppConfig.Bool("app.task.run.vod")

	//
	use_ssdb_vod_db = beego.AppConfig.String("ssdb.vod.db")

	//启动前加载
	beego.AddAPPStartHook(func() error {
		once.Do(func() {
			if app_task_run_vod {
				fmt.Println("user vodcenter reptile module init...")
				tskName := "user_vod_center_reptile_task"
				spec := "0 */30 * * * *"
				_spec := beego.AppConfig.String("reptile.ucenter.interval.spec")
				if len(_spec) > 0 {
					spec = _spec
				}
				toolbox.AddTask(tskName, toolbox.NewTask(tskName, spec, func() error {
					fmt.Println(tskName + " run...")
					vss := &VodUcenterReptile{}
					vss.ReptileTask()
					return nil
				}))
			}
		})
		return nil
	})
	//注册计数器
	stat.RegisterCounter(MOD_NAME, &Vods{})

	//注册收藏功能
	collect.RegisterCompler(collect_vod_mod, &Vods{})
}

func vod_chan_sync_key(videoId int64) string {
	return fmt.Sprintf("vod_rep_" + strconv.Itoa(int(videoId)))
}
