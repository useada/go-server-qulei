package configor

import (
	"github.com/BurntSushi/toml"
)

func LoadConfig(path string, v interface{}) error {
	_, err := toml.DecodeFile(path, v)
	return err
}
