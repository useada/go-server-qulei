package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/consul"
	"a.com/go-server/common/locip"
	"a.com/go-server/common/minilog"
)

type Configor struct {
	Server  configor.ServerConfigor
	Consul  configor.ConsulConfigor
	Elastic ElasticConfigor
}

type ElasticConfigor struct {
	Hosts []string
	Auth  string
}

var (
	Conf  Configor
	Log   *minilog.Logger
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

	Log = minilog.NewLogger(Conf.Server.LogPath, "server", 5000)
	Log.WithFileLine("FATAL", "DEBG")
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
