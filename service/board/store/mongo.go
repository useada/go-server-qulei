package store

import (
	"a.com/go-server/common/mongo"
)

type db struct {
	Pool *mongo.Pool
}

func NewMongoRepo(pool *mongo.Pool) Store {
	return &db{Pool: pool}
}
