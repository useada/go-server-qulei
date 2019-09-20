package main

import (
	"fmt"
	"net"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/consul"
	"a.com/go-server/common/locip"
	"a.com/go-server/common/logger"
	"a.com/go-server/common/mysql"
	"a.com/go-server/common/tracing"

	"a.com/go-server/service/uploader/internal/cloud"
	"a.com/go-server/service/uploader/internal/handler"
	"a.com/go-server/service/uploader/internal/store"
)

type Config struct {
	Server   ServerConfig
	Discover DiscoverConfig
	Logger   logger.Config
	S3       cloud.Config
	Mysql    []mysql.Config
}

type ServerConfig struct {
	Name string
	Host string
	Port int
}

type DiscoverConfig struct {
	Addr string
}

var (
	Conf  Config
	LocIP string
)

func init() {
	if err := configor.LoadConfig("./configs/conf.toml", &Conf); err != nil {
		panic(err)
	}

	var err error
	if LocIP, err = locip.GetLocalIP(); err != nil {
		panic(err)
	}

	if _, err := tracing.InitTracing(Conf.Server.Name); err != nil {
		panic(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", Conf.Server.Host)
	if err != nil {
		fmt.Println("failed to listen:", err)
		panic(err)
	}

	grpcSvr := grpc.NewServer(
		grpc.UnaryInterceptor(tracing.GrpcServerInterceptor(opentracing.GlobalTracer())),
	)

	handler.RegisterHandler(grpcSvr,
		store.NewMysqlRepo(mysql.NewPool(Conf.Mysql)),
		cloud.NewS3Repo(cloud.NewS3Client(Conf.S3)),
		logger.InitLogger(Conf.Logger))

	if err := consul.NewRegister(Conf.Discover.Addr).
		Registe(Conf.Server.Name, LocIP, Conf.Server.Port); err != nil {
		panic(err)
	}
	consul.RegisterGrpcHealth(grpcSvr)

	if err := grpcSvr.Serve(listener); err != nil {
		fmt.Println("failed to serve:", err)
		panic(err)
	}
}
