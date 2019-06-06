package mysql

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"a.com/go-server/common/tracing"
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

func (c *Client) Doit(ctx context.Context, h func(*gorm.DB) error) error {
	if c == nil {
		return errors.New("mysql instance is nil")
	}

	span := tracing.StartDBSpan(ctx, "mysql", "do")
	defer span.Finish()

	return h(c.DB)
}
