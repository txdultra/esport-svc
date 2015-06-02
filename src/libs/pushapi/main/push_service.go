package main

import (
	"flag"
	"fmt"
	"libs/pushapi"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/thrift"
)

var push_baidu_apikey, push_baidu_secret string
var bidupusher *pushapi.BaiduPusher

func init() {
	//baidu 云推送
	push_baidu_apikey = beego.AppConfig.String("push.baidu.apikey")
	push_baidu_secret = beego.AppConfig.String("push.baidu.secret")
	if len(push_baidu_apikey) > 0 && len(push_baidu_secret) > 0 {
		bidupusher = pushapi.NewBaiduPusher(push_baidu_apikey, push_baidu_secret)
	}
}

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port]:")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var host string
	var port int
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost:19092", "Specify host and port")
	flag.IntVar(&port, "p", 19092, "Specify port")
	flag.Parse()

	NetworkAddr := host

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	handler := &PushServiceImpl{}
	processor := pushapi.NewPushServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("push service server in", NetworkAddr)
	server.Serve()
}

type PushServiceImpl struct{}

func (p *PushServiceImpl) getBaiduPushType(pushType pushapi.PUSH_TYPE) pushapi.BAIDU_PUSH_TYPE {
	switch pushType {
	case pushapi.PUSH_TYPE_ALL:
		return pushapi.BAIDU_PUSH_TYPE_ALL
	case pushapi.PUSH_TYPE_GROUP:
		return pushapi.BAIDU_PUSH_TYPE_GROUP
	case pushapi.PUSH_TYPE_SINGLE:
		return pushapi.BAIDU_PUSH_TYPE_SINGLE
	default:
		return 0
	}
}

func (p *PushServiceImpl) getBaiduDeviceType(deviceType pushapi.DEVICE_TYPE) pushapi.BAIDU_DEVICE_TYPE {
	switch deviceType {
	case pushapi.DEVICE_TYPE_ANDROID:
		return pushapi.BAIDU_DEVICE_TYPE_ANDROID
	case pushapi.DEVICE_TYPE_BROWER:
		return pushapi.BAIDU_DEVICE_TYPE_BROWER
	case pushapi.DEVICE_TYPE_IOS:
		return pushapi.BAIDU_DEVICE_TYPE_IOS
	case pushapi.DEVICE_TYPE_PC:
		return pushapi.BAIDU_DEVICE_TYPE_PC
	case pushapi.DEVICE_TYPE_WP:
		return pushapi.BAIDU_DEVICE_TYPE_WP
	default:
		return 0
	}
}

func (p *PushServiceImpl) getBaiduMsgType(msgType pushapi.MSG_TYPE) pushapi.BAIDU_MSG_TYPE {
	switch msgType {
	case pushapi.MSG_TYPE_MSG:
		return pushapi.BAIDU_MSG_TYPE_MSG
	case pushapi.MSG_TYPE_NOTICE:
		return pushapi.BAIDU_MSG_TYPE_NOTICE
	default:
		return 0
	}
}

func (p *PushServiceImpl) getBaiduDeployStatus(status pushapi.DEPLOY_STATUS) pushapi.BAIDU_DEPLOY_STATUS {
	switch status {
	case pushapi.DEPLOY_STATUS_DEV:
		return pushapi.BAIDU_DEPLOY_STATUS_DEV
	case pushapi.DEPLOY_STATUS_REL:
		return pushapi.BAIDU_DEPLOY_STATUS_REL
	default:
		return 0
	}
}

func (p PushServiceImpl) PushMsg(data *pushapi.PushData) (r int64, err error) {
	if bidupusher == nil {
		panic("未配置百度推送apikey,secret参数")
	}
	pushType := p.getBaiduPushType(data.PushType)
	if int(pushType) == 0 {
		return 0, fmt.Errorf("PUSH_TYPE类型错误")
	}
	deviceType := p.getBaiduDeviceType(data.DeviceType)
	if int(deviceType) == 0 {
		return 0, fmt.Errorf("DEVICE_TYPE类型错误")
	}
	msgType := p.getBaiduMsgType(data.MsgType)
	if int(msgType) == 0 {
		return 0, fmt.Errorf("MSG_TYPE类型错误")
	}
	status := p.getBaiduDeployStatus(data.Status)
	if int(status) == 0 {
		return 0, fmt.Errorf("DEPLOY_STATUS类型错误")
	}

	msgs := make(map[string]interface{})
	for k, v := range data.Messages {
		msgs[k] = v
	}

	return bidupusher.PushMsg(
		data.Uid,
		pushType,
		data.ChannelId,
		data.Tag,
		deviceType,
		msgType,
		msgs,
		time.Unix(int64(data.Expries), 0),
		status,
		data.Args_...,
	)
}
