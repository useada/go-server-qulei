package mysql

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func Master(dbname string) *Client {
	instance, ok := gInstance[dbname]
	if !ok {
		return nil
	}
	return instance.getMaster()
}

func Slave(dbname string) *Client {
	instance, ok := gInstance[dbname]
	if !ok {
		return nil
	}
	return instance.getSlave()
}

type Client struct {
	*gorm.DB
}

func (c *Client) Doit(h func(*gorm.DB) error) error {
	if c == nil {
		return errors.New("mysql instance is nil")
	}
	return h(c.DB)
}
