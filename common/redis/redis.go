package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"a.com/go-server/common/tracing"
)

type Pool struct {
	Connections *redis.Pool
}

// -- Bytes Int String

// GetInt ...
func (p *Pool) GetInt(c context.Context, key string) (int, error) {
	return redis.Int(p.Get(c, key))
}

func (p *Pool) GetInt64(c context.Context, key string) (int64, error) {
	return redis.Int64(p.Get(c, key))
}

func (p *Pool) GetString(c context.Context, key string) (string, error) {
	return redis.String(p.Get(c, key))
}

func (p *Pool) GetBytes(c context.Context, key string) ([]byte, error) {
	return redis.Bytes(p.Get(c, key))
}

func (p *Pool) Get(c context.Context, key string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("GET", key)
		return err
	}
	return res, p.Doit(c, "get", handle)
}

func (p *Pool) Set(c context.Context, key string, val interface{}, ttl int64) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{key, val}
		if ttl > 0 {
			args = append(args, "EX", ttl)
		}
		_, err = conn.Do("SET", args...)
		return err
	}
	return p.Doit(c, "set", handle)
}

func (p *Pool) IncrBy(c context.Context, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("INCRBY", key, val))
		return err
	}
	return ret, p.Doit(c, "incrby", handle) // ret:INCRBY之后的值
}

func (p *Pool) DecrBy(c context.Context, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("DECRBY", key, val))
		return err
	}
	return ret, p.Doit(c, "decrby", handle) // ret:DECRBY之后的值
}

func (p *Pool) MGetBytes(c context.Context, keys []string) ([][]byte, error) {
	return redis.ByteSlices(p.MGet(c, keys))
}

func (p *Pool) MGet(c context.Context, keys []string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		args := make([]interface{}, 0)
		for _, key := range keys {
			args = append(args, key)
		}
		res, err = conn.Do("MGET", args...)
		return err
	}
	return res, p.Doit(c, "mget", handle)
}

// MSet args:  [key1, val1, key2, val2, ...]
func (p *Pool) MSet(c context.Context, args []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("MSET", args...)
		return err
	}
	return p.Doit(c, "mset", handle)
}

// -- Hash

// HGetInt64 ...
func (p *Pool) HGetInt64(c context.Context, hkey string, key string) (int64, error) {
	return redis.Int64(p.HGet(c, hkey, key))
}

func (p *Pool) HGetString(c context.Context, hkey string, key string) (string, error) {
	return redis.String(p.HGet(c, hkey, key))
}

func (p *Pool) HGetStrings(c context.Context, hkey, key string) ([]string, error) {
	return redis.Strings(p.HGet(c, hkey, key))
}

func (p *Pool) HGetBytes(c context.Context, hkey, key string) ([]byte, error) {
	return redis.Bytes(p.HGet(c, hkey, key))
}

func (p *Pool) HGetByteSlices(c context.Context, hkey, key string) ([][]byte, error) {
	return redis.ByteSlices(p.HGet(c, hkey, key))
}

func (p *Pool) HGet(c context.Context, hkey string, key string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("HGET", hkey, key)
		return err
	}
	return res, p.Doit(c, "hget", handle)
}

func (p *Pool) HSet(c context.Context, hkey string, key string, val interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("HSET", hkey, key, val)
		return err
	}
	return p.Doit(c, "hset", handle)
}

func (p *Pool) HDel(c context.Context, hkey, key string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("HDEL", hkey, key)
		return err
	}
	return p.Doit(c, "hdel", handle)
}

func (p *Pool) HIncrBy(c context.Context, hkey string, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("HINCRBY", hkey, key, val))
		return err
	}
	return ret, p.Doit(c, "hincrby", handle)
}

func (p *Pool) HMGetStrings(c context.Context, hkey string, keys []string) ([]string, error) {
	return redis.Strings(p.HMGet(c, hkey, keys))
}

func (p *Pool) HMGetBytes(c context.Context, hkey string, keys []string) ([][]byte, error) {
	return redis.ByteSlices(p.HMGet(c, hkey, keys))
}

func (p *Pool) HMGet(c context.Context, hkey string, keys []string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{hkey}
		for _, key := range keys {
			args = append(args, key)
		}
		res, err = conn.Do("HMGET", args...)
		return err
	}
	return res, p.Doit(c, "hmget", handle)
}

// HMSet vals:  [key1, val1, key2, val2, ...]
func (p *Pool) HMSet(c context.Context, hkey string, vals []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{hkey}
		args = append(args, vals...)
		_, err = conn.Do("HMSET", args...)
		return err
	}
	return p.Doit(c, "hmset", handle)
}

// HGetAll 效率原因, 不建议使用
func (p *Pool) HGetAll(c context.Context, hkey string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("HGETALL", hkey)
		return err
	}
	return res, p.Doit(c, "hgetall", handle)
}

func (p *Pool) HExists(c context.Context, hkey, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("HEXISTS", hkey, key))
		return err
	}
	return ret, p.Doit(c, "hexists", handle) // 0:不存在 1:存在 (err != nil, ret == 0)
}

// --Set

// SAdd ...
func (p *Pool) SAdd(c context.Context, skey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("SADD", skey, member)
		return err
	}
	return p.Doit(c, "sadd", handle)
}

func (p *Pool) SIsMember(c context.Context, skey, member string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("SISMEMBER", skey, member))
		return err
	}
	return ret, p.Doit(c, "sismember", handle) // 0:不存在 1:存在 (err != nil, ret == 0)
}

func (p *Pool) SRem(c context.Context, skey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("SREM", skey, member)
		return err
	}
	return p.Doit(c, "srem", handle)
}

// --Sorted Set

// ZAdd ...
func (p *Pool) ZAdd(c context.Context, zkey string, score int64, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZADD", zkey, score, member)
		return err
	}
	return p.Doit(c, "zadd", handle)
}

// ZMAdd vals: [score1, member1, score2, member2, ...]
func (p *Pool) ZMAdd(c context.Context, zkey string, vals []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey}
		args = append(args, vals...)
		_, err = conn.Do("ZADD", args...)
		return err
	}
	return p.Doit(c, "zmadd", handle)
}

func (p *Pool) ZScore(c context.Context, zkey string, member string) (score int64, err error) {
	handle := func(conn redis.Conn) error {
		score, err = redis.Int64(conn.Do("ZSCORE", zkey, member))
		return err
	}
	return score, p.Doit(c, "zscore", handle)
}

func (p *Pool) ZRangeByScore(c context.Context, zkey string, beg, end int64, limit int) (items []string, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey, beg, end}
		if limit > 0 {
			args = append(args, "LIMIT", 0, limit)
		}
		items, err = redis.Strings(conn.Do("ZRANGEBYSCORE", args...))
		return err
	}
	return items, p.Doit(c, "zrangebyscore", handle)
}

func (p *Pool) ZRevRangeByScore(c context.Context, zkey string, beg, end int64, limit int) (items []string, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey, beg, end}
		if limit > 0 {
			args = append(args, "limit", 0, limit)
		}
		items, err = redis.Strings(conn.Do("ZREVRANGEBYSCORE", args...))
		return err
	}
	return items, p.Doit(c, "zrevrangebyscore", handle)
}

func (p *Pool) ZRemRangeByScore(c context.Context, zkey string, beg, end int64) error {
	handle := func(conn redis.Conn) error {
		_, err := conn.Do("ZREMRANGEBYSCORE", zkey, beg, end)
		return err
	}
	return p.Doit(c, "zremrangebyscore", handle)
}

func (p *Pool) ZRem(c context.Context, zkey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZREM", zkey, member)
		return err
	}
	return p.Doit(c, "zrem", handle)
}

func (p *Pool) ZRemByScore(c context.Context, zkey string, min, max int64) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZREMRANGEBYSCORE", zkey, min, max)
		return err
	}
	return p.Doit(c, "zremrangebyscore", handle)
}

func (p *Pool) ZCard(c context.Context, zkey string) (val int64, err error) {
	handle := func(conn redis.Conn) error {
		val, err = redis.Int64(conn.Do("ZCARD", zkey))
		return err
	}
	return val, p.Doit(c, "zcard", handle)
}

// --List

// LPopString ...
func (p *Pool) LPopString(c context.Context, lkey string) (val string, err error) {
	return redis.String(p.LPop(c, lkey))
}

func (p *Pool) LPop(c context.Context, lkey string) (val interface{}, err error) {
	handle := func(conn redis.Conn) error {
		val, err = conn.Do("LPOP", lkey)
		return err
	}
	return val, p.Doit(c, "lpop", handle)
}

func (p *Pool) RPush(c context.Context, lkey string, val interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("RPUSH", lkey, val)
		return err
	}
	return p.Doit(c, "rpush", handle)
}

// --keys

// Exist ...
func (p *Pool) Exist(c context.Context, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("EXISTS", key))
		return err
	}
	return ret, p.Doit(c, "exists", handle) // 1:存在 0:不存在
}

func (p *Pool) Expire(c context.Context, key string, ttl int) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("EXPIRE", key, ttl))
		return err
	}
	return ret, p.Doit(c, "expire", handle) // 1:成功 0:失败
}

func (p *Pool) KeyTTL(c context.Context, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("TTL", key))
		return err
	}
	return ret, p.Doit(c, "ttl", handle) // -2:key不存在 -1:没有设置TTL num:剩余生存时间(秒)
}

func (p *Pool) Delete(c context.Context, key string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("DEL", key)
		return err
	}
	return p.Doit(c, "del", handle)
}

// --

func (p *Pool) Doit(c context.Context, cmd string, h func(redis.Conn) error) error {
	span := tracing.StartDBSpan(c, "redis", cmd)
	defer span.Finish()

	conn := p.Connections.Get()
	defer conn.Close()

	return h(conn)
}

type Config struct {
	Host    string
	Auth    string
	Index   int
	MaxIdle int `toml:"max_idle"`
}

func NewPool(conf Config) *Pool {
	pool := &Pool{
		Connections: &redis.Pool{
			MaxIdle:      conf.MaxIdle,
			MaxActive:    conf.MaxIdle * 3,
			Dial:         dialfunc(conf),
			TestOnBorrow: testfunc(),
			Wait:         true,
			IdleTimeout:  240 * time.Second,
		},
	}

	fmt.Println("初始化Redis连接池 FINISH")
	return pool
}

func dialfunc(conf Config) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", conf.Host,
			redis.DialConnectTimeout(time.Duration(100)*time.Millisecond),
			redis.DialReadTimeout(time.Duration(300)*time.Millisecond),
			redis.DialWriteTimeout(time.Duration(300)*time.Millisecond))
		if err != nil {
			return nil, err
		}
		if conf.Auth != "" {
			if _, err := c.Do("AUTH", conf.Auth); err != nil {
				c.Close()
				return nil, err
			}
		}
		if _, err := c.Do("SELECT", conf.Index); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	}
}

func testfunc() func(c redis.Conn, t time.Time) error {
	return func(c redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}

		_, err := c.Do("PING")
		return err
	}
}
