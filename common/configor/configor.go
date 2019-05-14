package configor

import (
	"github.com/BurntSushi/toml"
)

type ServerConfigor struct {
	Name string
	Host string
	Port int
}

type ConsulConfigor struct {
	Host string
}

type MysqlNode struct {
	Host    string
	Auth    string
	MaxIdle int `toml:"max_idle"`
	MaxOpen int `toml:"max_open"`
}
type MysqlConfigor struct {
	Name   string
	Option string
	Master MysqlNode
	Slave  []MysqlNode
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

type LoggerConfigor struct {
	FilePath   string `toml:"file_path"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Level      int
	Compress   bool
}

func LoadConfig(path string, v interface{}) error {
	_, err := toml.DecodeFile(path, v)
	return err
}
