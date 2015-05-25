package main

import (
	"bytes"
	_ "controllers"
	_ "docs"
	"errors"
	"fmt"
	"log"
	"net"
	//"strings"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/toolbox"
	"github.com/glycerine/go-capnproto"
	//"github.com/golang/groupcache"
	//"github.com/astaxie/beego/orm"
	"libs"
	//"libs/lives"
	//"libs/passport"

	//"libs/reptile"
	//"libs/collect"

	//"libs/credits/proxy"

	"libs/dlock"
	"libs/search"
	_ "libs/version"
	//"outobjs"
	//"modules/jobs"

	_ "net/http/pprof"
	_ "routers"

	"regexp"
	//"libs/comment"

	//"memcache"

	"runtime"
	"time"
	"utils"
	//"libs/share"

	//"libs/vod"
	//"libs/dlock"
	//"io/ioutil"
	//"os"
	//"path/filepath"
	//"encoding/base64"
	"flag"
	"libs/pushapi"
	"os"
	"os/exec"
	//"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yunge/sphinx"
	//"math/rand"
	//"logs"
	//"encoding/json"
	//"github.com/huichen/sego"

	//加载钩子程序

	_ "libs/utask/hook"

	cmd "github.com/Lupino/periodic/cmd/periodic/subcmd"
	pd "github.com/Lupino/periodic/driver"
	"github.com/Lupino/periodic/protocol"
)

func main() {

	fmt.Println(utils.UrlEncode("!@#$%^^&"))

	return

	watcherpath := "/watcher-test"
	watcher := dlock.NewWatcher()
	watcher.Write(watcherpath, []byte(fmt.Sprintf("new_%d", time.Now().Unix())))

	watcher.RegisterWatcher(watcherpath, func(data []byte) {
		fmt.Println("get rev :" + string(data))
	})

	go func() {
		for {
			time.Sleep(3 * time.Second)
			watcher.Write(watcherpath, []byte(fmt.Sprintf("new_%d", time.Now().Unix())))
		}
	}()

	time.Sleep(2 * time.Hour)

	return

	search_config := &search.SearchOptions{
		Host:    "192.168.0.79",
		Port:    9316,
		Timeout: 1000,
	}
	attrs := []string{"members", "threads"}
	val := []interface{}{uint64(1), 11, 13}
	values := [][]interface{}{
		val,
	}
	sph := search.NewSearcher(search_config)
	fmt.Println(sph.UpdateAttributes("group_idx", attrs, values))

	var filterRanges []search.FilterRangeInt
	filterRanges = append(filterRanges, search.FilterRangeInt{
		Attr:    "members",
		Min:     10,
		Max:     100,
		Exclude: false,
	})
	search_config.FilterRangeInt = filterRanges
	search_config.MaxMatches = 500
	search_config.Limit = 100
	sph2 := search.NewSearcher(search_config)
	fmt.Println(sph2.Query("", []string{"@weight DESC"}, "group_idx", "all"))
	return

	//	//	o := dbs.NewOrm("group_db")
	//	//	nums, _ := o.QueryTable("group").Exclude("status", 4).Count()
	//	//	fmt.Println(nums)
	//	fmt.Println(fmt.Sprintf("%03dA%d", 1, time.Now().UnixNano()/1000))
	//	//	ssdb.New("a").Zadd("aaa", "a", 1)
	//	//	ssdb.New("a").Zadd("aaa", "b", 1)
	//	//	fmt.Println(ssdb.New("a").Zcard("aaa"))

	//	ssdb.New("a").Hset("hset_test", "a_1", 1000)
	//	ssdb.New("a").Hset("hset_test", "a_2", 1001)

	//	kvs, _ := ssdb.New("a").Hgetall("hset_test")
	//	for i := 0; i < len(kvs); i += 2 {
	//		k := kvs[i]
	//		v, err := strconv.Atoi(kvs[i+1])
	//		if err != nil {
	//			continue
	//		}
	//		fmt.Println(k, ":", v)
	//	}
	//	val := groups.Post{
	//		ThreadId:   1,
	//		AuthorId:   1,
	//		Subject:    "32",
	//		Message:    "32",
	//		Ip:         "",
	//		FromDevice: groups.FROM_DEVICE_ANDROID,
	//		ImgIds:     []int64{1, 2, 3},
	//		ReplyId:    "",
	//		LongiTude:  1.12,
	//		LatiTude:   2.33,
	//	}
	//	var value interface{} = val
	//	var b bytes.Buffer
	//	encoder := gob.NewEncoder(&b)
	//	if err := encoder.Encode(value); err != nil {
	//		fmt.Printf("revel/cache: gob encoding '%s' failed: %s", value, err)
	//	}

	//cclient := ssdb.New("f")
	//	objs := cclient.MultiGet([]string{"group_post_001A1431080248022584", "group_post_001A1431079443586341"}, reflect.TypeOf(groups.Post{}))
	//	for _, ojb := range objs {
	//		fmt.Println(ojb.(*groups.Post))
	//	}
	//	cclient.Zadd("a123", "a", 10)
	//	cclient.Zadd("a123", "b", 11)
	//	cclient.Zadd("a123", "c", 12)
	//	cclient.Zadd("a123", "d", 13)
	//	cclient.Zadd("a123", "e", 14)

	//	fmt.Println(cclient.MultiZadd("b123", []interface{}{1, 2, 3}, []int64{100, 101, 102}))

	//	fmt.Println(cclient.MultiZget("a123", []interface{}{"b", "e", "e"}, reflect.TypeOf("")))
	//	fmt.Println(cclient.Zcard("group.invited_2_52_uids"))
	//	fmt.Println(cclient.Zscore("group.invited_2_52_uids", 123))

	//	objs, _ := cclient.Zscan("group.invited_2_52_uids", 0, 1<<40, 1<<10, reflect.TypeOf(int64(0))) //cclient.Zscan("b123", 0, 1<<32, 1<<10, reflect.TypeOf(int64(0)))
	//	fmt.Println(objs)
	//	for _, ojb := range objs {
	//		str := *(ojb.(*int64))
	//		fmt.Println(str)
	//	}

	//	//1431398695153046
	//	fmt.Println(1 << 53)

	type rtData struct {
		Id            uint64
		Pid           uint64
		Groupname     string
		Createtime    uint64
		Members       int
		Threads       int
		Gameids       string
		Displayorder  int
		Status        int
		Belong        int
		Type          int
		Vitality      int
		Searchkeyword string
		Longitude     float64
		Latitude      float64
		Recommend     int
		Starttime     uint64
		Endtime       uint64
	}

	//	rt := sphinx.NewClient(&sphinx.Options{
	//		Host:       "192.168.0.79",
	//		Port:       9316,
	//		SqlPort:    9317,
	//		Timeout:    5000,
	//		MaxMatches: 500,
	//	})
	//	err := rt.SetIndex("group_rt").Insert(&rtData{
	//		Id:            1,
	//		Pid:           1,
	//		Groupname:     "总共凤姐哦",
	//		Createtime:    uint64(time.Now().Unix()),
	//		Members:       0,
	//		Threads:       1,
	//		Gameids:       "1,2,3",
	//		Displayorder:  0,
	//		Status:        1,
	//		Belong:        1,
	//		Type:          1,
	//		Vitality:      0,
	//		Searchkeyword: "lol,英雄联盟",
	//		Longitude:     0,
	//		Latitude:      0,
	//		Recommend:     1,
	//		Starttime:     uint64(time.Now().Unix()),
	//		Endtime:       uint64(time.Now().Add(48 * time.Hour).Unix()),
	//	})
	//	fmt.Println(err)

	//	result, err := rt.SetIndex("group_rt").Execute("select * from group_rt")
	//	fmt.Println(result, err)

	//	res, err := rt.Query("共凤", "group_rt", "test rt insert")
	//	if err != nil {
	//		fmt.Println("TestInsert > %v\n", err)
	//		return
	//	}
	//	fmt.Println(res.TotalFound)
	var jobName = ""
	if false {
		if true {
			go func() {
				i := 0
				for {
					var s = capn.NewBuffer(nil)
					var job = pd.NewRootJob(s)
					jobName = utils.RandomStrings(5)
					job.SetName(jobName)
					job.SetFunc("test_func")
					job.SetArgs("")
					job.SetTimeout(1000)

					delay := 100
					job.SetSchedAt(int64(time.Now().Unix()) + int64(delay))
					cmd.SubmitJob("tcp://192.168.0.79:5000", job)
					i++
					if i > 5 {
						return
					}
				}
			}()

			time.Sleep(2 * time.Second)
			//return
		}

		//删除
		conn, err := net.Dial("tcp", "192.168.0.79:5000")
		if err != nil {
			fmt.Println(err)
			return
		}
		aconn := protocol.NewClientConn(conn)
		defer aconn.Close()
		err = aconn.Send(protocol.TYPE_CLIENT.Bytes())
		if err != nil {
			return
		}
		var s = capn.NewBuffer(nil)
		var job = pd.NewRootJob(s)
		job.SetName(jobName)
		job.SetFunc("test_func")
		buf := bytes.NewBuffer(nil)
		buf.Write([]byte(""))
		buf.Write(protocol.NULL_CHAR)
		buf.WriteByte(byte(protocol.REMOVE_JOB))
		buf.Write(protocol.NULL_CHAR)
		job.Segment.WriteTo(buf)
		err = aconn.Send(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		payload, err := aconn.Receive()
		if err != nil {
			log.Fatal(err)
		}
		_, cmd, _ := protocol.ParseCommand(payload)
		rlt := cmd.String()
		fmt.Println("jobName:", jobName, ":", rlt)
		//return

		//接收
		extraJob := func(payload []byte) (job pd.Job, jobHandle []byte, err error) {
			parts := bytes.SplitN(payload, protocol.NULL_CHAR, 4)
			if len(parts) != 4 {
				err = errors.New("Invalid payload " + string(payload))
				return
			}
			job, err = pd.ReadJob(parts[3])
			jobHandle = parts[2]
			return
		}

		conn, err = net.Dial("tcp", "192.168.0.79:5000")
		if err != nil {
			fmt.Println(err)
			return
		}
		pconn := protocol.NewClientConn(conn)
		defer pconn.Close()
		err = pconn.Send(protocol.TYPE_WORKER.Bytes())
		if err != nil {
			return
		}
		//var msgId = []byte(fmt.Sprintf("%d", time.Now().Unix()))
		buf = bytes.NewBuffer(nil)
		buf.Write([]byte(""))
		buf.Write(protocol.NULL_CHAR)
		buf.WriteByte(byte(protocol.CAN_DO))
		buf.Write(protocol.NULL_CHAR)
		buf.WriteString("test_func")
		err = pconn.Send(buf.Bytes())
		if err != nil {
			return
		}

		for {
			buf := bytes.NewBuffer(nil)
			buf.Write([]byte(""))
			buf.Write(protocol.NULL_CHAR)
			buf.Write(protocol.GRAB_JOB.Bytes())
			err = pconn.Send(buf.Bytes())
			if err != nil {
				fmt.Println(err)
				return
			}

			rdata, _ := pconn.Receive()
			fmt.Println("1_", rdata)

			job, jobHandle, _ := extraJob(rdata)
			fmt.Println("job name:", job.Name(), job.Id(), job.Args())

			buf = bytes.NewBuffer(nil)
			buf.Write([]byte(""))
			buf.Write(protocol.NULL_CHAR)
			buf.WriteByte(byte(protocol.WORK_DONE))
			buf.Write(protocol.NULL_CHAR)
			buf.Write(jobHandle)
			err = pconn.Send(buf.Bytes())
			if err != nil {
				return
			}
			fmt.Println(time.Now())
		}
		return
	}

	//	rep := reptile.NewYoukuReptileV2()
	//	fmt.Println(rep.Reptile("http://v.youku.com/v_show/id_XOTEzMzUxOTk2.html"))
	//	return

	//	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	//	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//	transport, _ := thrift.NewTSocket("192.168.0.79:19091")

	//	useTransport := transportFactory.GetTransport(transport)
	//	client := proxy.NewUserTaskServiceClientFactory(useTransport, protocolFactory)
	//	if err := transport.Open(); err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//	defer transport.Close()

	//	result, err := client.CreateGroup(&proxy.TaskGroup{
	//		GroupType:   proxy.TASK_GROUP_TYPE_LIST,
	//		TaskType:    proxy.TASK_TYPE_DUPLICATE,
	//		Name:        "测试重复3",
	//		Description: "测试重复描述3",
	//		Enabled:     true,
	//		StartTime:   time.Now().Unix(),
	//		EndTime:     time.Now().Add(time.Duration(30) * time.Duration(24) * time.Hour).Unix(),
	//	})

	//	fmt.Println(result, err)

	//	group, err := client.GetGroup(2)
	//	fmt.Println(group, err)

	//	result, err := client.CreateTask(&proxy.Task{
	//		TaskType:      proxy.TASK_TYPE_DUPLICATE,
	//		GroupId:       3,
	//		Name:          "任务222",
	//		Description:   "任务222描述",
	//		TaskLimits:    2,
	//		EventName:     "test",
	//		StartTime:     time.Now().Unix(),
	//		EndTime:       time.Now().Add(time.Duration(720) * time.Hour).Unix(),
	//		ResetTime:     300,
	//		Period:        2,
	//		PeriodType:    proxy.TASK_PERIOD_TYPE_DAY,
	//		Reward:        proxy.TASK_REWARD_TYPE_CREDIT,
	//		Prize:         10,
	//		ResetVar:      1,
	//		LastResetTime: 10,
	//	})
	//	fmt.Println(result, err)

	//	task, err := client.GetTask(result.Id)
	//	fmt.Println(task, err)

	//	missions, _ := client.GetMissions(1)
	//	for _, mission := range missions {
	//		fmt.Println(mission)
	//	}

	//	var uid int64 = 1
	//	eventName := "test"
	//	result, err := client.EventHandler(uid, eventName, 3)
	//	fmt.Println(result, err)

	//	missions, _ := client.GetMissions(uid)
	//	for _, mission := range missions {
	//		fmt.Println(mission)
	//	}

	//	group.Name = "测试修改1"
	//	group.Enabled = false
	//	group.StartTime = time.Now().Unix()
	//	result, err := client.UpdateGroup(group)
	//	fmt.Println(result, err)

	//	result, err := client.DeleteGroup(1)
	//	fmt.Println(result, err)
	//group, err := client.GetGroup(1)
	//fmt.Println(group, err)

	//	utask.TaskInit()

	//	tasker := utask.NewTasker()
	//	task := &utask.Task{
	//		TaskType:     utask.TASK_TYPE_DUPLICATE,
	//		GroupId:      1,
	//		Name:         "吃不下",
	//		Description:  "佛界哦我就覅哦额我",
	//		Icon:         0,
	//		TaskLimits:   3,
	//		ApplyPerm:    "",
	//		EventName:    "submit_share",
	//		StartTime:    time.Now().Unix(),
	//		EndTime:      time.Now().Add(720 * time.Hour).Unix(),
	//		RedoTime:     20,
	//		Period:       3,
	//		PeriodType:   utask.TASK_PERIOD_TYPE_DAY,
	//		Reward:       utask.TASK_REWARD_TYPE_CREDIT,
	//		Prize:        20,
	//		Applicants:   1,
	//		Achievers:    0,
	//		Version:      "v1.0",
	//		DisplayOrder: 10,
	//	}
	//	fmt.Println(tasker.CreateTask(task))

	//	fmt.Println("t1:", tasker.GetTask(task.TaskId))
	//	fmt.Println("gt1:", tasker.GetGroupTasks(task.GroupId))
	//	fmt.Println("et1:", tasker.GetEventNameTasks(task.EventName))

	//	task.Prize = 30
	//	fmt.Println(tasker.UpdateTask(task))

	//	fmt.Println("-----------------------------------")

	//	fmt.Println("t2:", tasker.GetTask(task.TaskId))
	//	fmt.Println("gt2:", tasker.GetGroupTasks(task.GroupId))
	//	fmt.Println("et2:", tasker.GetEventNameTasks(task.EventName))

	//	fmt.Println("------------------------------------")

	//	tasker.DeleteTask(1, 1)
	//	fmt.Println("t3:", tasker.GetTask(task.TaskId))
	//	fmt.Println("gt3:", tasker.GetGroupTasks(task.GroupId))
	//	fmt.Println("et3:", tasker.GetEventNameTasks(task.EventName))
	//	fmt.Println("g3:", tasker.GetGroup(task.GroupId))
	//	return

	//	group := &utask.TaskGroup{
	//		GroupId:     1,
	//		TaskType:    utask.TASK_TYPE_REDO,
	//		Name:        "测试1",
	//		Description: "测试1描述修",
	//		Enabled:     true,
	//		StartTime:   time.Now().Unix(),
	//		EndTime:     time.Now().Unix(),
	//	}
	//	fmt.Println(tasker.UpdateGroup(group))
	//fmt.Println(tasker.DeleteGroup(2))

	//	credit_service := credits.NewCreditService()
	//	oper := &credits.OperationCreditParameter{
	//		No:        "001A0000000001N1423977676060",
	//		Operation: credits.CREDIT_OPERATION_UNLOCK,
	//	}
	//	result := credit_service.OperCredit(oper)
	//	fmt.Println(result)
	//time.Sleep(60 * time.Second)
	//return

	//for i := 1; i < 2; i++ {
	//	ii := i
	//go func() {
	//	oper := &credits.OperationCreditParameter{
	//		Uid:       1,
	//		Points:    10,
	//		Desc:      fmt.Sprintf("批量锁减测试数据%d", 1),
	//		Operation: credits.CREDIT_OPERATION_LOCKINCR,
	//		Ref:       "aa",
	//		RefId:     "bb",
	//	}
	//	result := credit_service.OperCredit(oper)
	//	fmt.Println(result)
	//}()
	//}
	//	time.Sleep(60 * time.Second)
	//	return

	//c := make(chan string)
	//for i := 0; i < 10; i++ {
	//	time.AfterFunc(utils.StrToDuration(fmt.Sprintf("%ds", i)), func() {
	//		c <- fmt.Sprintf("%d", i)
	//	})
	//}
	//for i := 0; i < 200; i++ {
	//	go func() {
	//		locker := dlock.NewDistributedLock()
	//		dtlock, conn, err := locker.GetLock("/c_djq_u_12312879")
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		defer conn.Close()
	//		err = dtlock.Lock()
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		fmt.Println(time.Now().UnixNano(), ":", err)
	//		//<-c
	//		dtlock.Unlock()
	//		return
	//	}()
	//}

	//time.Sleep(100 * time.Second)
	//return

	//huomao := reptile.HuomaoLive{}
	//fmt.Println(huomao.GetStatus("http://www.huomaotv.com/live/13"))
	//return

	//sns := share.NewShareMsgService()
	//sn := &share.Share{
	//	Id:         1,
	//	Uid:        2,
	//	CreateTime: time.Now(),
	//	ShareType:  1,
	//}

	//for i := 0; i < 100; i++ {
	//	sn.Id = int64(i) + 10
	//	sns.SendMsgToQueue(sn, int64(i+1))
	//	time.Sleep(100 * time.Millisecond)
	//}

	//cc := make(chan string)

	//gp.RegisterGroupCache("share-group", func(_ groupcache.Context, key string, dest groupcache.Sink) error {
	//	fmt.Println(key)
	//	return nil
	//})

	//group := gp.GetGroupCache("share-group")

	//var str string
	//group.Get(nil, "abc", groupcache.StringSink(&str))
	//fmt.Println("111111", str, time.Now().UnixNano())
	//group.Get(nil, "abc", groupcache.StringSink(&str))
	//fmt.Println("222222", str, time.Now().UnixNano())
	//group.Get(nil, "abc", groupcache.StringSink(&str))
	//fmt.Println("333333", str, time.Now().UnixNano())
	//<-cc
	//redis.Set(nil, "fewfewe123", 1000322)
	//fmt.Println(redis.Get(nil, "fewfewe123"))

	//for i := 0; i < 2000; i++ {
	//	go httplib.Get("http://192.168.0.79:8080/v1/count/all?access_token=ZmNlv8e01NSCpbXl8pb53gMtpbZcpht8").String()
	//}

	//for i := 0; i < 1000; i++ {
	//	go httplib.Get("http://192.168.0.79:8080/v1/share/public_timeline?access_token=ZmNlv8e01NSCpbXl8pb53gMtpbZcpht8").String()
	//}

	//ssdb.New("c").MultiDel([]string{"abc111"})
	//kvs := map[string]int64{
	//	"1": 123,
	//	"2": 321,
	//	"3": 214,
	//}
	//fmt.Println(ssdb.New("c").Hmset("fjifewifow", kvs))
	//fmt.Println(ssdb.New("c").Hkeys("fjifewifow"))
	//ssdb.NewSSDBClient("c").Set("abc222", "abc")
	//ssdb.NewSSDBClient("c").Set("abc333", "abc")
	//fmt.Println(ssdb.NewSSDBClient("c").MultiGet([]string{"abc111", "abc222", "abc333"}, reflect.TypeOf("")))

	//arr := ssdb.NewSSDBClient("c").MultiGet([]string{"share_cmt:10221419239098307", "share_cmt:1022221419239098307"}, reflect.TypeOf(share.ShareComment{}))
	//for _, sc := range arr {
	//	_sc, ok := sc.(*share.ShareComment)
	//	if ok {
	//		fmt.Println(*_sc)
	//	}
	//}

	//fmt.Println(ssdb.NewSSDBClient("c").Zcard("share_cmts:7200000010"))
	//fmt.Println(ssdb.NewSSDBClient("c").Zrange("share_cmts:7200000010", 0, -1, reflect.TypeOf(share.ShareComment{})))

	//for i := 0; i < 1000; i++ {
	//	key := fmt.Sprintf("member_collect_box:%d", i)
	//	ssdb.NewSSDBClient("b").Zclear(key)
	//}

	//function calc_hash_tbl($u, $n = 256, $m = 16){
	//$h = sprintf("%u", crc32($u));
	//$h1 = intval($h / $n);
	//$h2 = $h1 % $n;
	//$h3 = base_convert($h2, 10, $m);
	//$h4 = sprintf("%02s", $h3);
	//return $h4;
	//}

	//scs := share.NewShareComments()
	//scs.Create(&share.ShareComment{
	//	Sid:     1,
	//	SUid:    10,
	//	Uid:     2,
	//	RUid:    0,
	//	Content: "减肥我就f",
	//})

	//scs.Delete(10221419144816116)

	if false {
		for i := 0; i < 1000; i++ {
			ii := i
			go func() {
				sc := sphinx.NewClient(&sphinx.Options{
					Host: "192.168.0.33",
					Port: 9312,
					//Socket:     "192.168.0.33:9312",
					Timeout:    10000,
					Offset:     0,
					Limit:      20,
					MaxMatches: 500,
					MatchMode:  sphinx.SPH_MATCH_ANY,
				})
				res, err := sc.Query("老鼠", "vods_idx", "")
				if err == nil {
					fmt.Println(res.TotalFound, "dddd:", ii)
				} else {
					fmt.Println(err)
				}
			}()
		}
	}

	//option := proxy.NewSearchOptions()
	//option.Host = "192.168.0.33"
	//option.Port = 9312
	//option.Timeout = 1000
	//option.Offset = 0
	//option.Limit = 20
	//option.MaxMatches = 500
	//option.MatchMode = proxy.SPH_MATCH_SPH_MATCH_ANY
	////option.RankMode = proxy.SPH_RANK_SPH_RANK_PROXIMITY_BM25
	//option.SortMode = proxy.SPH_SORT_SPH_SORT_EXTENDED
	//option.SortBy = "@weight DESC,post_time DESC"
	//option.FieldWeights = make(map[string]int32)
	//option.Filters = []*proxy.SearchFilter{}
	//option.FilterRangeInts = []*proxy.FilterRangeInt{}
	//option.FilterRangeFloats = []*proxy.FilterRangeFloat{}
	//option.Excerpts = &proxy.SearchExcerpt{}

	//for iii := 0; iii < 10000; iii++ {
	//	go func() {
	//		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	//		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//		transport, err := thrift.NewTSocketTimeout(net.JoinHostPort("127.0.0.1", "9300"))
	//		defer transport.Close()
	//		if err != nil {
	//			fmt.Fprintln(os.Stderr, "error resolving address:", err)
	//			//os.Exit(1)
	//			return
	//		}
	//		useTransport := transportFactory.GetTransport(transport)
	//		client := proxy.NewSearchServiceClientFactory(useTransport, protocolFactory)
	//		if err := transport.Open(); err != nil {
	//			fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9300", " ", err)
	//			//os.Exit(1)
	//			return
	//		}
	//		result, _ := client.Query(option, "老鼠", "vods_idx")
	//		fmt.Println(result)
	//	}()
	//}

	runes := []rune("总积分f")
	fmt.Println(runes[0] % 9)

	fmt.Println(12398124 % 9)

	fmt.Println(utils.IntToCh(1001))
	now := time.Now()
	fmt.Println(now)
	fmt.Println(now.UnixNano())
	convert := func(nano int64) time.Time {
		seconds := nano / 1000000000
		_nano := nano % 1000000000
		fmt.Println(seconds)
		fmt.Println(_nano)
		return time.Unix(seconds, _nano)
	}
	fmt.Println(convert(now.UnixNano()))

	MOBILE_AGENT_REGEX := `(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino`
	fmt.Println(regexp.MatchString(MOBILE_AGENT_REGEX, `Mozilla/5.0 (iphone; U; CPU iPhone OS 5_1_1 like Mac OS X; da-dk) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3`))

	//fmt.Println(reptile.Get_REP_SUPPORT("http://www.huomaotv.com/live/23"))

	//uids, _ := ssdb.NewSSDBClient(1).Zrevrange("friendship_follower_uid:9", 0, 100, reflect.TypeOf(int64(0)))
	//fmt.Println("++++++++++++++++++++++++++++++++++++++++$$$", len(uids))

	//ssdb.NewSSDBClient(1).Zremrangebyscore("friendship_follower_uid:9", float64(time.Now().Unix()), 1000000000000000000000000)

	//lps := &lives.LivePrograms{}
	//fmt.Println(lps.GetsByDate("<=", time.Now(), 0))

	//rediskey := "kkkkkkkkkkkkkkkkkkkkk"
	type AAA struct {
		Name string
		Age  int
	}
	//fmt.Println(redis.Del(rediskey))
	//fmt.Println(redis.ZCard(rediskey))
	//redis.ZAdd(rediskey, float64(time.Now().Unix()), AAA{Name: "a", Age: 1})
	//redis.ZAdd(rediskey, float64(time.Now().Add(10*time.Second).Unix()), AAA{Name: "b", Age: 1})
	//redis.ZAdd(rediskey, float64(time.Now().Add(20*time.Second).Unix()), AAA{Name: "c", Age: 1})
	//redis.ZAdd(rediskey, float64(time.Now().Add(30*time.Second).Unix()), AAA{Name: "d", Age: 1})
	//redis.ZAdd(rediskey, float64(time.Now().Add(40*time.Second).Unix()), AAA{Name: "e", Age: 1})

	//fmt.Println(redis.ZCard(rediskey))
	//tt, _ := redis.Zrevrangebyscore(rediskey, 1415688835, 0, 5, reflect.TypeOf(AAA{}))
	//for _, t := range tt {
	//	fmt.Println(t.(*AAA))
	//}

	//key1 := "abc_1"
	//key2 := "abc_2"
	//key3 := "abc_1_2"

	//redis.Del(key1)
	//redis.Del(key2)
	//redis.Del(nil, key3)

	messages := make(map[string]interface{})
	messages["title"] = "尼玛"
	messages["description"] = " 大B,收到没:" + time.Now().String()

	//pusher := pushapi.NewBaiduPusher("uhHQ6k7yEGg7Yb8DDaeKXvhd", "jtL3IroLeGjwIgaET1qTdYgtk8dx3zIu")
	//rlt, err := pusher.PushMsg("840282558323329949", pushapi.BAIDU_PUSH_TYPE_SINGLE, "5319908774962730520", "", pushapi.BAIDU_DEVICE_TYPE_IOS, pushapi.BAIDU_MSG_TYPE_NOTICE, messages, time.Now().Add(10*time.Minute))
	//fmt.Println(rlt, err)

	pusher := pushapi.NewBaiduPusher("uhHQ6k7yEGg7Yb8DDaeKXvhd", "jtL3IroLeGjwIgaET1qTdYgtk8dx3zIu")
	//	rlt, err := pusher.PushMsg("952532336315496486", pushapi.BAIDU_PUSH_TYPE_SINGLE, "4688375275015171206", "",
	//		pushapi.BAIDU_DEVICE_TYPE_IOS, pushapi.BAIDU_MSG_TYPE_NOTICE,
	//		messages, time.Now().Add(10*time.Minute), pushapi.BAIDU_DEPLOY_STATUS_DEV)
	//	fmt.Println("push:", rlt, err)
	fmt.Println(pusher)

	//mcs := &passport.MemberConfigs{}
	//mca := passport.NewMemberConfigAttrs()
	//mca.AllowPush = false
	//mcs.SetConfig(1, mca)
	//fmt.Println(mcs.GetConfig(1))

	//redis.ZAdd(key1, float64(time.Now().Unix()), AAA{Name: "a", Age: 1})
	//redis.ZAdd(key1, float64(time.Now().Add(10*time.Second).Unix()), AAA{Name: "b", Age: 1})

	//redis.ZAdd(key2, float64(time.Now().Unix()), AAA{Name: "c", Age: 1})
	//redis.ZAdd(key2, float64(time.Now().Add(10*time.Second).Unix()), AAA{Name: "d", Age: 1})

	//redis.ZUnionStore(nil, key1, 2, []string{key1, key2}, "weights", 0, 0)

	//fmt.Println(redis.ZCard(nil, key1))

	//zt, _ := redis.ZRevRangeByScore(nil, key1, reflect.TypeOf(AAA{}), float64(time.Now().Unix()), "-inf")
	//for _, t := range zt {
	//	fmt.Println(t.(*AAA))
	//}

	//fmt.Println(redis.ZScore(nil, key1, AAA{Name: "c", Age: 1}))

	//aaaa := func(args ...interface{}) {
	//	for i := range args {
	//		fmt.Println(i)
	//	}
	//}

	//ars := []interface{}{"a", "b", "c"}
	//aaaa(ars...)

	//ttt := time.Unix(0, 0)
	//fmt.Println(ttt)
	//if utils.IsZero(ttt) {
	//	fmt.Println("zero")
	//}

	//collect.DeleteAll("ff", "ff")
	//fmt.Println(redis.Zrevrange("redis.share.my.subscriptions:52", 0, 100, reflect.TypeOf(share.SimpleShare{})))
	//pst := &passport.FriendShips{}

	//fmt.Println(pst.FollowerIdsP(97, 500, 1))
	//run := false
	//if run {
	//	pms := &lives.LivePrograms{}
	//	pid, _ := pms.Create(lives.LiveProgram{
	//		Title:     "测试三",
	//		SubTitle:  "测试三二级",
	//		Date:      time.Now(),
	//		StartTime: time.Now().Add(1 * time.Minute),
	//		EndTime:   time.Now().Add(30 * time.Minute),
	//		MatchId:   1,
	//		PostTime:  time.Now(),
	//		PostUid:   1,
	//		GameId:    1,
	//		ChannelId: 1,
	//	})
	//	psms := &lives.LiveSubPrograms{}
	//	psid, _ := psms.Create(&lives.LiveSubProgram{
	//		ProgramId: pid,
	//		GameId:    1,
	//		Vs1Name:   "尼玛",
	//		Vs1Img:    1,
	//		Vs1Uid:    1,
	//		Vs2Name:   "我操",
	//		Vs2Img:    2,
	//		Vs2Uid:    2,
	//		ViewType:  lives.LIVE_SUBPROGRAM_VIEW_VS,
	//		Title:     "",
	//		Img:       0,
	//		StartTime: time.Now().Add(2 * time.Minute),
	//		EndTime:   time.Now().Add(10 * time.Minute),
	//		PostTime:  time.Now(),
	//		PostUid:   1,
	//	})
	//	//fmt.Println(psid)

	//	notices := &lives.ProgramNoticer{}
	//	err := notices.SubscribeNotice(3, pid, []int64{psid})
	//	fmt.Println(err)
	//}

	//fmt.Println(utils.Md5(fmt.Sprintf("%s_%s_%s_%s", "LAv5VYFhVqMm3ub6unQR4ClURqsJVW1c", "neotv", "1414057911", "34326H4L937zaR3nfE")))

	//org := &lives.LiveOrgs{}
	//fmt.Println(org.DeleteStream(1))

	//auth := new(passport.OpenIDOAuth)
	//auth.AuthType = byte(1)
	//auth.Uid = 1
	//auth.AuthIdentifier = "fesf"
	//auth.AuthEmail, auth.AuthEmailVerified = "", ""
	//auth.AuthPreferredUserName, auth.AuthProviderName = "", ""
	//auth.AuthToken = "fejiojiofoe"
	//auth.AuthDate = time.Now()
	//auth.OpenIDUid = "fesf"
	//auth_provider := passport.NewOpenIDAuth()
	//_, auth_err := auth_provider.Create(auth)
	//fmt.Println(auth_err)

	//loc, _ := time.LoadLocation("Asia/Shanghai")
	//t, err := time.ParseInLocation("2006-01-02 15:04:05", "2012-10-01 11:13:23", loc)
	//fmt.Println(t, err, time.Now())

	//url := "http://hdla.douyutv.com/live/22391r7kK4aD2Cn7.flv?wsSecret=4985c7f068a19bd425977e5dca789a7b&wsTime=1413458645"
	//rep, _ := httplib.Get(url).Response()
	//fmt.Println(rep.Request.URL)

	//lps := &lives.LivePers{}
	//fmt.Println(lps.ListForAdmin("", 0, 2, 5))

	//redis.Hset("test_abcddd123", "a", 321)
	//var val int
	//redis.Hget("test_abcddd123", "a", &val)
	//redis.Hclear("test_abcddd123")
	//redis.Lpush("test_kkkk", "ccc")
	//fmt.Println(redis.Llen("test_kkkk"))
	//fmt.Println(val)

	//redis.Zadd("test_0001", 123, float64(1))
	//redis.Zadd("test_0001", 321, float64(2))
	//fmt.Println(redis.Zcard("test_0001"))
	//fmt.Println(redis.Zexists("test_0001", 321))

	//fmt.Println(redis.MultiZadd("test_0001", []interface{}{11, 22, 33}, []float64{1, 2, 3}))
	//fmt.Println(redis.MultiZdel("test_0001", []interface{}{11, 2}))
	//fmt.Println(redis.Zcard("test_0001"))
	//fmt.Println(redis.Zexists("test_0001", 22))

	//passport.IncrAtCount(1, 2)
	//passport.IncrAtCount(1, 3)
	//passport.IncrAtCount(1, 4)
	//passport.IncrAtCount(1, 3)
	//passport.IncrAtCount(1, 3)
	//passport.IncrAtCount(1, 5)
	//passport.IncrAtCount(1, 6)

	//pers := lives.LivePers{}
	//p := pers.Get(14)
	//p.Name = "尼玛直播3"
	//pers.Update(*p)
	//pers.UpdateGames(14, []int{1, 2})
	//abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=~!@#$%^&*()_+[]{}\\|:;\"\'<>,.?/

	//ms := passport.NewMemberProvider()
	//i, ids := ms.Query("老党", 1, 20, "any", nil, nil)
	//fmt.Println(i, ids)

	//sp := search.NewSearcher(&search.SearchOptions{
	//	Host:       "192.168.15.129",
	//	Port:       9312,
	//	Timeout:    10000,
	//	Offset:     0,
	//	Limit:      20,
	//	MaxMatches: 100,
	//})
	//fmt.Println(sp.VideoQuery("老鼠", "any"))

	//fmt.Println(regexp.MatchString(`^.*(\d[a-zA-Z]|[a-zA-Z]\d).*$`, "fjw2iwo"))
	//for _, _s := range lst {
	//	fmt.Println(*(_s.(*A)))
	//}

	//f := func(source string, salt string) string {
	//	return utils.Md5(utils.Md5(source) + salt)
	//}

	//salt := utils.RandomStrings(5)
	//fmt.Println(f("123123", salt), salt)

	//t := time.Now().AddDate(0, 0, -10)
	//fmt.Println(utils.FriendTime(t))

	//m := cache.NewMemcachedCache([]string{"127.0.0.1:12000"}, 10*time.Second, 50*time.Millisecond)
	//type ppp struct {
	//	Name  string
	//	Value int
	//}
	//p := ppp{
	//	Name:  "abc",
	//	Value: 123,
	//}
	//m.Set("ddd", p, 30*time.Second)

	//cc, _ := memcache.Connect("127.0.0.1:12000")
	//stored, err := cc.Set("ddd", 0, uint64(60*60), []byte("ddddddd"))
	//if err == nil && stored == false {
	//	fmt.Println("stored fail")
	//	return
	//}
	//cs := 1
	//for j := 0; j < 1000; j++ {
	//	for i := 0; i < 1000; i++ {
	//		go func() {
	//			//_p := ppp{}
	//			//err := m.Get("ddd", &_p)
	//			//if err != nil {
	//			//	cs++
	//			//	fmt.Println(err)
	//			//	fmt.Println(cs)
	//			//} else {
	//			//	fmt.Println(_p)
	//			//}

	//			_c, _ := memcache.Connect("127.0.0.1:12000")
	//			v, err := _c.Get("ddd")
	//			if err != nil {
	//				cs++
	//				fmt.Println(err)
	//				fmt.Println(cs)
	//			} else {
	//				var contain interface{}
	//				if len(v) > 0 {
	//					contain = string(v[0].Value)
	//				} else {
	//					contain = nil
	//				}
	//				fmt.Println(contain)
	//			}
	//			_c.Close()
	//		}()
	//	}
	//	time.Sleep(2 * time.Second)
	//}

	//fmt.Println(time.Unix(1406592000, 0))
	//fmt.Println(time.Now().UnixNano())
	//ps := make(map[string]interface{})
	//ps["ref_id"] = int64(11)
	//ps["from_id"] = int64(20)
	//ps["content"] = "jfoiewjfoiw从joie我jo分级网e"
	//ps["ip"] = "127.0.0.1"
	//ps["reply_id"] = "53d61519abbe091bb8000001"
	//ps["longitude"] = 0.021
	//ps["latitude"] = 2.131
	//cmter := comment.NewCommenter("vod")
	//err := cmter.Create(ps, func(tos []string) []int64 { return []int64{} }, false)
	//fmt.Println(err)

	//_, cms := cmter.Gets(11, 1, 10)
	//for _, v := range cms {
	//	fmt.Println(v)
	//}

	//fmt.Println(utils.TotalPages(10, 2))

	//pattern := `@[^@\s]*?[:：，,.。 ]`
	//re := regexp.MustCompile(pattern)
	//us := re.FindAllString("fjioe@Jiof @iojo @jio_few:@jiofew @上方_揭瓦 ", -1)
	//fmt.Println(us)
	//uss := utils.ExtractMentions("fjioe@Jiof @iojo @jio_few:@jiofew @上方_揭瓦")
	//fmt.Println(uss)

	//pattern := `@[^@]+?(?=[\s:：(),.。@])`
	//rxp := rubex.MustCompile(pattern)
	//result := rxp.FindAllString("fjioe@Jiof @iojo@jiofew:@jiofew")
	//fmt.Println(result)
	fmt.Println("current server version:", libs.APP_SERVER_VER)
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
