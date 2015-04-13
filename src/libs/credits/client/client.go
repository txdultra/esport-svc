package client

import (
	"libs/credits/proxy"
	"logs"

	"github.com/astaxie/beego"
	"github.com/thrift"
)

var creditServerHost string = beego.AppConfig.String("credit.host")

//记得关闭transport
func NewClient() (*proxy.CreditServiceClient, thrift.TTransport, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, _ := thrift.NewTSocket(creditServerHost)

	useTransport := transportFactory.GetTransport(transport)
	client := proxy.NewCreditServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		logs.Errorf("credit client transport open fail:%s", err.Error())
		return nil, nil, err
	}
	return client, transport, nil
}
