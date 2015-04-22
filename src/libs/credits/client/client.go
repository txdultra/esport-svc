package client

import (
	"libs/credits/proxy"
	"logs"

	"github.com/astaxie/beego"
	"github.com/thrift"
)

//默认积分系统地址
var creditServerHost string = beego.AppConfig.String("credit.host")

//记得关闭transport
func NewClient(host string) (*proxy.CreditServiceClient, thrift.TTransport, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	var transport *thrift.TSocket
	if len(host) > 0 {
		transport, _ = thrift.NewTSocket(host)
	} else {
		transport, _ = thrift.NewTSocket(creditServerHost)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := proxy.NewCreditServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		logs.Errorf("credit client transport open fail:%s", err.Error())
		return nil, nil, err
	}
	return client, transport, nil
}
