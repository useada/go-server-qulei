package main

import (
	"fmt"

	"a.com/server/mywork/common/configor"
	"a.com/server/mywork/common/minilog"
	"a.com/server/mywork/gclient"
)

var (
	Conf Configor
	Log  *minilog.Logger
)

func init() {
	if err := configor.LoadConfig("./conf/conf.toml", &Conf); err != nil {
		panic(err)
	}

	if err := gclient.InitGrpcs(Conf.Grpc.Consul, Conf.Grpc.Services); err != nil {
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
