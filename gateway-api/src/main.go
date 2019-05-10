package main

import (
	"go.uber.org/zap"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/logger"
	"a.com/go-server/gclient"
)

type Configor struct {
	Server configor.ServerConfigor
	Logger configor.LoggerConfigor
	Grpc   GrpcClients
}

type GrpcClients struct {
	Consul   string
	Services []string
}

var (
	Conf Configor
	Log  *zap.SugaredLogger
)

func init() {
	if err := configor.LoadConfig("./conf/conf.toml", &Conf); err != nil {
		panic(err)
	}

	Log = logger.InitLogger(Conf.Logger)

	if err := gclient.Init(Conf.Grpc.Consul, Conf.Grpc.Services); err != nil {
		panic(err)
	}
	Log.Info(Conf.Grpc.Consul, Conf.Grpc.Services)
}

func main() {
	router := Router()

	router.Run(Conf.Server.Host)
}
