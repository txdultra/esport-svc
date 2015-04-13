package client

import (
	"libs/utask/proxy"
	"logs"

	"github.com/thrift"
)

//记得关闭transport
func NewClient() (*proxy.UserTaskServiceClient, thrift.TTransport, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, _ := thrift.NewTSocket(proxy.UTaskServerHost)

	useTransport := transportFactory.GetTransport(transport)
	client := proxy.NewUserTaskServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		logs.Errorf("utask client transport open fail:%s", err.Error())
		return nil, nil, err
	}
	return client, transport, nil
}
