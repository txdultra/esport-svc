package main

import (
	_ "dbs"
	//"docs"
	"fmt"
	_ "logs"

	_ "github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/toolbox"
	//"github.com/astaxie/beego/httplib"
	//_ "github.com/astaxie/beego/orm"
	_ "libs"
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
	"os"
	"os/exec"
)

func main() {

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
