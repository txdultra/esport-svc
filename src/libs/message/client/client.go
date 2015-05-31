package client

import (
	"libs/message/service"
	"logs"

	"github.com/thrift"
)

//记得关闭transport
func NewClient(serverHost string) (*service.MessageServiceClient, thrift.TTransport, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, _ := thrift.NewTSocket(serverHost)

	useTransport := transportFactory.GetTransport(transport)
	client := service.NewMessageServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		logs.Errorf("message client transport open fail:%s", err.Error())
		return nil, nil, err
	}
	return client, transport, nil
}
