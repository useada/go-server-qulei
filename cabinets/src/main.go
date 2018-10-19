package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"a.com/server/mywork/common/configor"
	"a.com/server/mywork/common/consul"
	"a.com/server/mywork/common/locip"
	"a.com/server/mywork/common/minilog"
	"a.com/server/mywork/common/mysql"
)

var (
	Conf  Configor
	Log   *minilog.Logger
	LocIP string
)

func init() {
	if err := configor.LoadConfig("./conf/conf.toml", &Conf); err != nil {
		panic(err)
	}

	if err := mysql.InitMysql(Conf.Mysql); err != nil {
		panic(err)
	}

	var err error
	if LocIP, err = locip.GetLocalIP(); err != nil {
		panic(err)
	}

	InitS3Client(Conf.S3)

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
