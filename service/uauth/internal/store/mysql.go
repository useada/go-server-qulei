package store

import (
	"a.com/go-server/common/mysql"
)

type db struct {
	Pool *mysql.Pool
}

func NewMysqlRepo(pool *mysql.Pool) Store {
	return &db{Pool: pool}
}
