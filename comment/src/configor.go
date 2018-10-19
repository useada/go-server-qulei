package main

import (
	configor "a.com/go-server/common/configor"
)

type Configor struct {
	Server configor.ServerConfigor
	Consul configor.ConsulConfigor
	Redis  configor.RedisConfigor
	Mongo  configor.MongoConfigor
}
