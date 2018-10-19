package main

import (
	configor "a.com/server/mywork/common/configor"
)

type Configor struct {
	Server  configor.ServerConfigor
	Consul  configor.ConsulConfigor
	Elastic ElasticConfigor
}

type ElasticConfigor struct {
	Hosts []string
	Auth  string
}
