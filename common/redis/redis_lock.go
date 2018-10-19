package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
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
	conn := GetConn()
	if nil == conn {
		return "", errors.New("get redis conn failed")
	}
	defer conn.Close()

	ticket := fmt.Sprintf("%s%d", s, time.Now().UnixNano())
	res, err := conn.Do("SET", key(s), ticket, "PX", ttl, "NX")
	if err != nil {
		return "", err
	}
	if res != "OK" {
		return "", errors.New("Failed to acquire lock")
	}
	return ticket, nil
}

func UnLock(s, ticket string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	ret, err := UnlockScript.Do(conn, key(s), ticket)
	if ret == 0 {
		return errors.New("unlock script failed")
	}
	return err
}

func key(s string) string {
	return fmt.Sprintf("RedisLock|%s", s)
}
