package hook

import (
	"fmt"
	"libs/utask/proxy"
	"logs"

	"libs/hook"

	"github.com/thrift"
)

func init() {
	ac := &AsyncClient{}
	hook.RegisterHook("upload_avatar", "utask_upload_avatar", ac)
	hook.RegisterHook("fans_count", "utask_fans_count", ac)
	hook.RegisterHook("gz_count", "utask_gz_count", ac)
	hook.RegisterHook("friend_count", "utask_friend_count", ac)
	hook.RegisterHook("create_share", "utask_create_share", ac)
	hook.RegisterHook("vod_comment", "utask_vod_comment", ac)
	hook.RegisterHook("share_qt_comment", "utask_share_qt_comment", ac)
	hook.RegisterHook("subscr_vuser", "utask_subscr_vuser", ac)
	hook.RegisterHook("fx_vod", "utask_fx_vod", ac)
	hook.RegisterHook("live_vod_view_time", "utask_live_vod_view_time", ac)
	hook.RegisterHook("live_sub_program", "utask_live_sub_program", ac)
	hook.RegisterHook("group_modify_bgimg", "utask_group_modify_bgimg", ac)
	hook.RegisterHook("group_new_thread", "utask_group_new_thread", ac)
	hook.RegisterHook("group_new_post", "utask_group_new_post", ac)
	hook.RegisterHook("group_ding_post", "utask_group_ding_post", ac)
	hook.RegisterHook("share_weixin", "utask_share_outside", ac)
	fmt.Println("utask hook registed")
}

type AsyncClient struct{}

func (c *AsyncClient) Do(event string, args ...interface{}) {
	if len(args) < 2 {
		return
	}
	uid, ok1 := args[0].(int64)
	n, ok2 := args[1].(int)
	if !ok1 || !ok2 {
		return
	}

	go func() {
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		transport, _ := thrift.NewTSocket(proxy.UTaskServerHost)

		useTransport := transportFactory.GetTransport(transport)
		client := proxy.NewUserTaskServiceClientFactory(useTransport, protocolFactory)
		if err := transport.Open(); err != nil {
			logs.Errorf("utask async commit transport open fail:%s", err.Error())
			return
		}
		defer transport.Close()
		client.EventHandler(uid, event, int32(n))
	}()
}
