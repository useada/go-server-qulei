package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var UnlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

func TryLock(s string, ttl int64) (string, error) {
	ticket := fmt.Sprintf("%s%d", s, time.Now().UnixNano())
	handle := func(conn redis.Conn) error {
		res, err := conn.Do("SET", key(s), ticket, "PX", ttl, "NX")
		if res != "OK" {
			return errors.New("Failed to acquire lock")
		}
		return err
	}
	return ticket, Doit(handle)
}

func UnLock(s, ticket string) error {
	handle := func(conn redis.Conn) error {
		ret, err := UnlockScript.Do(conn, key(s), ticket)
		if ret == 0 {
			return errors.New("unlock script failed")
		}
		return err
	}
	return Doit(handle)
}

func key(s string) string {
	return fmt.Sprintf("RedisLock|%s", s)
}
