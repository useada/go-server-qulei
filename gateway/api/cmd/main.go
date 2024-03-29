package main

import (
	"a.com/go-server/common/configor"
	"a.com/go-server/common/logger"
	"a.com/go-server/common/tracing"
	"a.com/go-server/gclient"

	"a.com/go-server/gateway/api/internal/handler"
	"a.com/go-server/gateway/api/internal/router"
)

type Config struct {
	Server ServerConfig
	Grpc   gclient.Config
	Logger logger.Config
}

type ServerConfig struct {
	Name string
	Host string
	Port int
}

var (
	Conf Config
)

func init() {
	if err := configor.LoadConfig("./configs/conf.toml", &Conf); err != nil {
		panic(err)
	}

	if _, err := tracing.InitTracing(Conf.Server.Name); err != nil {
		panic(err)
	}
}

func main() {
	handle := handler.NewApiHandler(gclient.NewGrpcClient(Conf.Grpc),
		logger.InitLogger(Conf.Logger))
	router.BindRouter(handle).Run(Conf.Server.Host)
}
