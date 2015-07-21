package main

import (
	_ "controllers"
	_ "dbs"
	_ "docs"
	"fmt"
	_ "logs"

	_ "github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/toolbox"
	//"github.com/astaxie/beego/httplib"
	//_ "github.com/astaxie/beego/orm"
	_ "libs"
	_ "libs/lives"
	_ "libs/passport"
	//"libs/reptile"
	//"libs/collect"
	_ "libs/version"
	//"libs/vod"
	//"outobjs"
	//"modules/jobs"
	_ "net/http/pprof"
	_ "routers"
	//"utils/cache"
	//"regexp"
	//"libs/comment"
	//"memcache"
	"runtime"
	//"time"
	//"utils"
	//"libs/search"
	//"io/ioutil"
	//"os"
	//"path/filepath"
	//"encoding/base64"
	//"utils/redis"
	"flag"
	_ "libs/utask/hook"
	"os"
	"os/exec"
)

func main() {

	//vs := &vod.Vods{}
	//page := 1
	//for {
	//	_, vods := vs.DbQuery(nil, page, 100)
	//	if len(vods) == 0 {
	//		return
	//	}
	//	for _, v := range vods {
	//		vc := vs.GetCount(v.Id)
	//		if vc != nil {
	//			rint := vs.RangeViews(v.GameId)
	//			vc.Views += rint
	//			vs.UpdateCount(vc)
	//		}
	//	}
	//	page++
	//}

	//o := dbs.NewDefaultOrm()
	//mps := passport.NewMemberProvider()
	//var members []*passport.Member
	//o.QueryTable(&passport.Member{}).OrderBy("create_time").All(&members)
	//for _, m := range members {
	//	mgs := mps.MemberGames(m.Uid)
	//	gids := []int{}
	//	for _, mg := range mgs {
	//		gids = append(gids, mg.GameId)
	//	}
	//	mps.UpdateMemberGids(m.Uid, gids)
	//}

	runtime.GOMAXPROCS(runtime.NumCPU())
	beego.Run()
	fmt.Println("Server Running...")
}

var godaemon = flag.Bool("d", false, "run app as a daemon with -d=true")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if *godaemon {
		args := os.Args[1:]
		i := 0
		for ; i < len(args); i++ {
			if args[i] == "-d=true" {
				args[i] = "-d=false"
				break
			}
		}
		cmd := exec.Command(os.Args[0], args...)
		cmd.Start()
		fmt.Println("[PID]", cmd.Process.Pid)
		os.Exit(0)
	}
}
