package main

/*
import (
	"fmt"
	//"labix.org/v2/mgo/bson"
	"libs"
	//"sync"
	"testing"
	//"time"
)


func TestProgramNotice(t *testing.T) {
	n := &libs.ProgramNotices{}
	err := n.SubscribeNotice(2, 1, []int64{1, 2})
	fmt.Println(err)

	nsn := libs.NewSubsNotice("notices", "programs")
	counts, subs := nsn.GetSubscribes(1, 10, map[string]interface{}{"subs": map[string]interface{}{
		"obj_id":    3,
		"is_notice": true,
	},
	})
	fmt.Println(counts)
	for _, v := range subs {
		fmt.Println(v)
	}
}



func TestSubscribe(t *testing.T) {
	//n := libs.NewSubsNotice("notices", "programs")
	//for i := 1; i < 1000; i++ {
	//	sub := libs.SubscribeNotice{
	//		Pid:    1,
	//		FromId: int64(i),
	//		RefIds: []int64{1, 2, 3},
	//	}
	//	n.Subscribe(sub)
	//}
	//fmt.Println("subscribe completed...")

	//lsp := &libs.LiveSubPrograms{}
	//sprm := lsp.Get(1)
	//fmt.Println(sprm)

	//dur := time.Now().Sub(sprm.StartTime)
	//_dur := sprm.StartTime.Sub(time.Now())
	//fmt.Println(dur.Minutes())
	//fmt.Println(_dur.Minutes())
	//if dur.Minutes() < 3.0 {
	//	fmt.Println("过于接近开始时间")
	//}

	//n := libs.NewSubsNotice("notices", "programs")
	//_, ss := n.GetSubscribers(1, 100, map[string]interface{}{"refids": 1})
	//for _, v := range ss {
	//	fmt.Println(v)
	//}

	pn := &libs.ProgramNotices{}
	err := pn.StartNoticeTimer(1)
	fmt.Println(err)
	c := make(chan bool)
	<-c
}

*/
