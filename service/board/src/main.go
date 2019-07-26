package main

import (
	"fmt"
	"net"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/consul"
	"a.com/go-server/common/locip"
	"a.com/go-server/common/logger"
	"a.com/go-server/common/mongo"
	"a.com/go-server/common/redis"
	"a.com/go-server/common/tracing"
)

type Config struct {
	Server    ServerConfig
	Discovery DiscoveryConfig
	Logger    logger.Config
	Redis     redis.Config
	Mongo     mongo.Config
}

type ServerConfig struct {
	Name string
	Host string
	Port int
}

type DiscoveryConfig struct {
	Addr string
}

var (
	Conf  Config
	Log   *zap.SugaredLogger
	LocIP string
)

func init() {
	if err := configor.LoadConfig("./conf/conf.toml", &Conf); err != nil {
		panic(err)
	}

	Log = logger.InitLogger(Conf.Logger)

	if _, err := tracing.InitTracing(Conf.Server.Name); err != nil {
		panic(err)
	}

	if err := mongo.Init(Conf.Mongo); err != nil {
		panic(err)
	}

	redis.Init(Conf.Redis)

	var err error
	if LocIP, err = locip.GetLocalIP(); err != nil {
		panic(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", Conf.Server.Host)
	if err != nil {
		fmt.Println("failed to listen:", err)
		panic(err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(tracing.GrpcServerInterceptor(opentracing.GlobalTracer())),
	)
	RegisterHandler(server)

	if err := consul.NewRegister(Conf.Discovery.Addr).
		Registe(Conf.Server.Name, LocIP, Conf.Server.Port); err != nil {
		panic(err)
	}
	consul.RegisterGrpcHealth(server)

	if err := server.Serve(listener); err != nil {
		fmt.Println("failed to serve:", err)
		panic(err)
	}
}
