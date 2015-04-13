package main

import (
	"flag"
	"fmt"
	"libs/credits/proxy"
	"os"

	"github.com/thrift"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port]:")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var host string
	var port int
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost:19090", "Specify host and port")
	flag.IntVar(&port, "p", 19090, "Specify port")
	flag.Parse()

	NetworkAddr := host

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	handler := &proxy.CreditServiceProxy{}
	processor := proxy.NewCreditServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("credit service server in", NetworkAddr)
	server.Serve()
}
