package main

import (
	configor "a.com/server/mywork/common/configor"
)

type Configor struct {
	Server configor.ServerConfigor
	Consul configor.ConsulConfigor
	S3     S3Configor
	Mysql  configor.MysqlConfigor
}

type S3Configor struct {
	Region string
	ACL    string
	Bucket string
}
