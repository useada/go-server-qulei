package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"a.com/go-server/common/tracing"
)

// --Bytes Int String
func GetInt(c context.Context, key string) (int, error) {
	return redis.Int(Get(c, key))
}

func GetInt64(c context.Context, key string) (int64, error) {
	return redis.Int64(Get(c, key))
}

func GetString(c context.Context, key string) (string, error) {
	return redis.String(Get(c, key))
}

func GetBytes(c context.Context, key string) ([]byte, error) {
	return redis.Bytes(Get(c, key))
}

func Get(c context.Context, key string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("GET", key)
		return err
	}
	return res, Doit(c, "get", handle)
}

func Set(c context.Context, key string, val interface{}, ttl int64) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{key, val}
		if ttl > 0 {
			args = append(args, "EX", ttl)
		}
		_, err = conn.Do("SET", args...)
		return err
	}
	return Doit(c, "set", handle)
}

func IncrBy(c context.Context, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("INCRBY", key, val))
		return err
	}
	return ret, Doit(c, "incrby", handle) // ret:INCRBY之后的值
}

func DecrBy(c context.Context, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("DECRBY", key, val))
		return err
	}
	return ret, Doit(c, "decrby", handle) // ret:DECRBY之后的值
}

func MGetBytes(c context.Context, keys []string) ([][]byte, error) {
	return redis.ByteSlices(MGet(c, keys))
}

func MGet(c context.Context, keys []string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		args := make([]interface{}, 0)
		for _, key := range keys {
			args = append(args, key)
		}
		res, err = conn.Do("MGET", args...)
		return err
	}
	return res, Doit(c, "mget", handle)
}

// args:  [key1, val1, key2, val2, ...]
func MSet(c context.Context, args []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("MSET", args...)
		return err
	}
	return Doit(c, "mset", handle)
}

// -- Hash
func HGetInt64(c context.Context, hkey string, key string) (int64, error) {
	return redis.Int64(HGet(c, hkey, key))
}

func HGetString(c context.Context, hkey string, key string) (string, error) {
	return redis.String(HGet(c, hkey, key))
}

func HGetStrings(c context.Context, hkey, key string) ([]string, error) {
	return redis.Strings(HGet(c, hkey, key))
}

func HGetBytes(c context.Context, hkey, key string) ([]byte, error) {
	return redis.Bytes(HGet(c, hkey, key))
}

func HGetByteSlices(c context.Context, hkey, key string) ([][]byte, error) {
	return redis.ByteSlices(HGet(c, hkey, key))
}

func HGet(c context.Context, hkey string, key string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("HGET", hkey, key)
		return err
	}
	return res, Doit(c, "hget", handle)
}

func HSet(c context.Context, hkey string, key string, val interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("HSET", hkey, key, val)
		return err
	}
	return Doit(c, "hset", handle)
}

func HDel(c context.Context, hkey, key string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("HDEL", hkey, key)
		return err
	}
	return Doit(c, "hdel", handle)
}

func HIncrBy(c context.Context, hkey string, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("HINCRBY", hkey, key, val))
		return err
	}
	return ret, Doit(c, "hincrby", handle)
}

func HMGetStrings(c context.Context, hkey string, keys []string) ([]string, error) {
	return redis.Strings(HMGet(c, hkey, keys))
}

func HMGetBytes(c context.Context, hkey string, keys []string) ([][]byte, error) {
	return redis.ByteSlices(HMGet(c, hkey, keys))
}

func HMGet(c context.Context, hkey string, keys []string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{hkey}
		for _, key := range keys {
			args = append(args, key)
		}
		res, err = conn.Do("HMGET", args...)
		return err
	}
	return res, Doit(c, "hmget", handle)
}

// vals:  [key1, val1, key2, val2, ...]
func HMSet(c context.Context, hkey string, vals []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{hkey}
		args = append(args, vals...)
		_, err = conn.Do("HMSET", args...)
		return err
	}
	return Doit(c, "hmset", handle)
}

// 效率原因, 不建议使用
func HGetAll(c context.Context, hkey string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("HGETALL", hkey)
		return err
	}
	return res, Doit(c, "hgetall", handle)
}

func HExists(c context.Context, hkey, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("HEXISTS", hkey, key))
		return err
	}
	return ret, Doit(c, "hexists", handle) // 0:不存在 1:存在 (err != nil, ret == 0)
}

// --Set
func SAdd(c context.Context, skey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("SADD", skey, member)
		return err
	}
	return Doit(c, "sadd", handle)
}

func SIsMember(c context.Context, skey, member string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("SISMEMBER", skey, member))
		return err
	}
	return ret, Doit(c, "sismember", handle) // 0:不存在 1:存在 (err != nil, ret == 0)
}

func SRem(c context.Context, skey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("SREM", skey, member)
		return err
	}
	return Doit(c, "srem", handle)
}

// --Sorted Set
func ZAdd(c context.Context, zkey string, score int64, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZADD", zkey, score, member)
		return err
	}
	return Doit(c, "zadd", handle)
}

// vals: [score1, member1, score2, member2, ...]
func ZMAdd(c context.Context, zkey string, vals []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey}
		args = append(args, vals...)
		_, err = conn.Do("ZADD", args...)
		return err
	}
	return Doit(c, "zmadd", handle)
}

func ZScore(c context.Context, zkey string, member string) (score int64, err error) {
	handle := func(conn redis.Conn) error {
		score, err = redis.Int64(conn.Do("ZSCORE", zkey, member))
		return err
	}
	return score, Doit(c, "zscore", handle)
}

func ZRangeByScore(c context.Context, zkey string, beg, end int64, limit int) (items []string, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey, beg, end}
		if limit > 0 {
			args = append(args, "LIMIT", 0, limit)
		}
		items, err = redis.Strings(conn.Do("ZRANGEBYSCORE", args...))
		return err
	}
	return items, Doit(c, "zrangebyscore", handle)
}

func ZRevRangeByScore(c context.Context, zkey string, beg, end int64, limit int) (items []string, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey, beg, end}
		if limit > 0 {
			args = append(args, "limit", 0, limit)
		}
		items, err = redis.Strings(conn.Do("ZREVRANGEBYSCORE", args...))
		return err
	}
	return items, Doit(c, "zrevrangebyscore", handle)
}

func ZRemRangeByScore(c context.Context, zkey string, beg, end int64) error {
	handle := func(conn redis.Conn) error {
		_, err := conn.Do("ZREMRANGEBYSCORE", zkey, beg, end)
		return err
	}
	return Doit(c, "zremrangebyscore", handle)
}

func ZRem(c context.Context, zkey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZREM", zkey, member)
		return err
	}
	return Doit(c, "zrem", handle)
}

func ZRemByScore(c context.Context, zkey string, min, max int64) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZREMRANGEBYSCORE", zkey, min, max)
		return err
	}
	return Doit(c, "zremrangebyscore", handle)
}

func ZCard(c context.Context, zkey string) (val int64, err error) {
	handle := func(conn redis.Conn) error {
		val, err = redis.Int64(conn.Do("ZCARD", zkey))
		return err
	}
	return val, Doit(c, "zcard", handle)
}

// --List
func LPopString(c context.Context, lkey string) (val string, err error) {
	return redis.String(LPop(c, lkey))
}

func LPop(c context.Context, lkey string) (val interface{}, err error) {
	handle := func(conn redis.Conn) error {
		val, err = conn.Do("LPOP", lkey)
		return err
	}
	return val, Doit(c, "lpop", handle)
}

func RPush(c context.Context, lkey string, val interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("RPUSH", lkey, val)
		return err
	}
	return Doit(c, "rpush", handle)
}

// --keys
func Exist(c context.Context, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("EXISTS", key))
		return err
	}
	return ret, Doit(c, "exists", handle) // 1:存在 0:不存在
}

func Expire(c context.Context, key string, ttl int) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("EXPIRE", key, ttl))
		return err
	}
	return ret, Doit(c, "expire", handle) // 1:成功 0:失败
}

func KeyTTL(c context.Context, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("TTL", key))
		return err
	}
	return ret, Doit(c, "ttl", handle) // -2:key不存在 -1:没有设置TTL num:剩余生存时间(秒)
}

func Delete(c context.Context, key string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("DEL", key)
		return err
	}
	return Doit(c, "del", handle)
}

// --

func Doit(c context.Context, cmd string, h func(redis.Conn) error) error {
	span := tracing.StartDBSpan(c, "redis", cmd)
	defer span.Finish()

	conn := gRedigo.Get()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()
	return h(conn)
}

var gRedigo *redis.Pool

type RedisConfigor struct {
	Host    string
	Auth    string
	Index   int
	MaxIdle int `toml:"max_idle"`
}

func Init(conf RedisConfigor) {
	fmt.Println("初始化Redis连接池")
	gRedigo = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Host)
			if nil != err {
				return nil, err
			}
			if conf.Auth != "" {
				if _, err := c.Do("AUTH", conf.Auth); err != nil {
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", conf.Index); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: 240 * time.Second,
	}
}
