package cache

import (
	"a.com/go-server/common/redis"
)

type kv struct {
	Pool *redis.Pool
}

func NewRedisRepo(pool *redis.Pool) Cache {
	return &kv{Pool: pool}
}
