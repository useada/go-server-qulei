package main

import (
	configor "a.com/go-server/common/configor"
)

type Configor struct {
	Server configor.ServerConfigor
	Grpc   GrpcClients
}

type GrpcClients struct {
	Consul   string
	Services []string
}
