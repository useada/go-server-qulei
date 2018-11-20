package main

import (
	"fmt"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/minilog"
	"a.com/go-server/gclient"
)

type Configor struct {
	Server configor.ServerConfigor
	Grpc   GrpcClients
}

type GrpcClients struct {
	Consul   string
	Services []string
}

var (
	Conf Configor
	Log  *minilog.Logger
)

func init() {
	if err := configor.LoadConfig("./conf/conf.toml", &Conf); err != nil {
		panic(err)
	}

	if err := gclient.Init(Conf.Grpc.Consul, Conf.Grpc.Services); err != nil {
		panic(err)
	}
	fmt.Println(Conf.Grpc.Consul, Conf.Grpc.Services)

	Log = minilog.NewLogger(Conf.Server.LogPath, "api", 5000)
	Log.WithFileLine("FATAL", "DEBG")
}

func main() {
	router := Router()

	router.Run(Conf.Server.Host)
}
