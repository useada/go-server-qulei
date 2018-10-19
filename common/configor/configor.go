package configor

import (
	"github.com/BurntSushi/toml"
)

type ServerConfigor struct {
	Name    string
	Host    string
	Port    int
	LogPath string
}

type ConsulConfigor struct {
	Host string
}

type MysqlConfigor struct {
	Host     string
	Auth     string
	MaxIdle  int `toml:"max_idle"`
	MaxOpen  int `toml:"max_open"`
	Database []string
}

type MongoConfigor struct {
	Host     string
	Auth     string
	Database []string
}

type RedisConfigor struct {
	Host    string
	Auth    string
	Index   int
	MaxIdle int `toml:"max_idle"`
}

func LoadConfig(path string, v interface{}) error {
	_, err := toml.DecodeFile(path, v)
	return err
}
