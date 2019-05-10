package main

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/consul"
	"a.com/go-server/common/locip"
	"a.com/go-server/common/logger"
)

type Configor struct {
	Server  configor.ServerConfigor
	Logger  configor.LoggerConfigor
	Consul  configor.ConsulConfigor
	Elastic ElasticConfigor
}

type ElasticConfigor struct {
	Hosts []string
	Auth  string
}

var (
	Conf  Configor
	Log   *zap.SugaredLogger
	LocIP string
)

func init() {
	if err := configor.LoadConfig("./conf/conf.toml", &Conf); err != nil {
		panic(err)
	}

	if err := InitElasticClient(Conf.Elastic); err != nil {
		panic(err)
	}

	var err error
	if LocIP, err = locip.GetLocalIP(); err != nil {
		panic(err)
	}

	Log = logger.InitLogger(Conf.Logger)
}

func main() {
	listener, err := net.Listen("tcp", Conf.Server.Host)
	if err != nil {
		fmt.Println("failed to listen:", err)
		panic(err)
	}
	server := grpc.NewServer()
	RegisterHandler(server)

	if err := consul.NewConsulRegister(Conf.Consul).
		Register(Conf.Server.Name, LocIP, Conf.Server.Port); err != nil {
		panic(err)
	}
	consul.RegisterGrpcHealth(server)

	if err := server.Serve(listener); err != nil {
		fmt.Println("failed to serve:", err)
		panic(err)
	}
}
