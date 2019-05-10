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
	"a.com/go-server/common/mysql"
)

type Configor struct {
	Server configor.ServerConfigor
	Logger configor.LoggerConfigor
	Consul configor.ConsulConfigor
	S3     S3Configor
	Mysql  configor.MysqlConfigor
}

type S3Configor struct {
	Region string
	ACL    string
	Bucket string
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

	if err := mysql.Init(Conf.Mysql); err != nil {
		panic(err)
	}

	var err error
	if LocIP, err = locip.GetLocalIP(); err != nil {
		panic(err)
	}

	InitS3Client(Conf.S3)

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
