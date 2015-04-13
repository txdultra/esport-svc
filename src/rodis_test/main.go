package main

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"sync"
	"time"
	"utils"
)

func main() {
	wait := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		go func() {
			wait.Add(1)
			defer wait.Done()
			data, _ := httplib.Get("http://121.40.203.148:8080/v1/share/timeline?access_token=VVx5TX6Me4hkkeQNdL5forBUL6txuFzq&uid=22").Bytes()
			json, err := utils.NewJson(data)
			if err != nil {
				fmt.Println("read fail:" + err.Error())
				//continue
				return
			}
			j, ok := json.CheckGet("lists")
			if ok {
				arr, _ := j.Array()
				if len(arr) > 0 {
					fmt.Println(fmt.Sprintf("read succsss, length:%d,%d", len(arr), time.Now().Unix()))
				} else {
					fmt.Println(fmt.Sprintf("read fail, length:%d", 0))
				}
			} else {
				fmt.Println("read fail:" + err.Error())
			}
		}()
	}
	wait.Wait()
}
