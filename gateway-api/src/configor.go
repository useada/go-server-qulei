package main

import (
	configor "a.com/server/mywork/common/configor"
)

type Configor struct {
	Server configor.ServerConfigor
	Grpc   GrpcClients
}

type GrpcClients struct {
	Consul   string
	Services []string
}
