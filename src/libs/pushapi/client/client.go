package client

import (
	"libs/pushapi"
	"logs"

	"github.com/astaxie/beego"
	"github.com/thrift"
)

//默认积分系统地址
var pushServerHost string = beego.AppConfig.String("push.host")

//记得关闭transport
func NewClient(host string) (*pushapi.PushServiceClient, thrift.TTransport, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	var transport *thrift.TSocket
	if len(host) > 0 {
		transport, _ = thrift.NewTSocket(host)
	} else {
		transport, _ = thrift.NewTSocket(pushServerHost)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := pushapi.NewPushServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		logs.Errorf("push client transport open fail:%s", err.Error())
		return nil, nil, err
	}
	return client, transport, nil
}

func Push(data *pushapi.PushData) error {
	client, transport, err := NewClient("")
	if err != nil {
		return err
	}
	defer transport.Close()
	_, err = client.PushMsg(data)
	return err
}
