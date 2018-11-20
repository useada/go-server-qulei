package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"a.com/go-server/common/configor"
)

// --Bytes Int String
func GetInt(key string) (int, error) {
	return redis.Int(Get(key))
}

func GetInt64(key string) (int64, error) {
	return redis.Int64(Get(key))
}

func GetString(key string) (string, error) {
	return redis.String(Get(key))
}

func GetBytes(key string) ([]byte, error) {
	return redis.Bytes(Get(key))
}

func Get(key string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("GET", key)
		return err
	}
	return res, Doit(handle)
}

func Set(key string, val interface{}, ttl int64) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{key, val}
		if ttl > 0 {
			args = append(args, "EX", ttl)
		}
		_, err = conn.Do("SET", args)
		return err
	}
	return Doit(handle)
}

func IncrBy(key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("INCRBY", key, val))
		return err
	}
	return ret, Doit(handle) // ret:INCRBY之后的值
}

func DecrBy(key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("DECRBY", key, val))
		return err
	}
	return ret, Doit(handle) // ret:DECRBY之后的值
}

func MGetBytes(keys []string) ([][]byte, error) {
	return redis.ByteSlices(MGet(keys))
}

func MGet(keys []string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		args := make([]interface{}, 0)
		for _, key := range keys {
			args = append(args, key)
		}
		res, err = conn.Do("MGET", args...)
		return err
	}
	return res, Doit(handle)
}

// args:  [key1, val1, key2, val2, ...]
func MSet(args []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("MSET", args)
		return err
	}
	return Doit(handle)
}

// -- Hash
func HGetInt64(hkey string, key string) (int64, error) {
	return redis.Int64(HGet(hkey, key))
}

func HGetString(hkey string, key string) (string, error) {
	return redis.String(HGet(hkey, key))
}

func HGetStrings(hkey, key string) ([]string, error) {
	return redis.Strings(HGet(hkey, key))
}

func HGetBytes(hkey, key string) ([]byte, error) {
	return redis.Bytes(HGet(hkey, key))
}

func HGetByteSlices(hkey, key string) ([][]byte, error) {
	return redis.ByteSlices(HGet(hkey, key))
}

func HGet(hkey string, key string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("HGET", hkey, key)
		return err
	}
	return res, Doit(handle)
}

func HSet(hkey string, key string, val interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("HSET", hkey, key, val)
		return err
	}
	return Doit(handle)
}

func HDel(hkey, key string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("HDEL", hkey, key)
		return err
	}
	return Doit(handle)
}

func HIncrBy(hkey string, key string, val int64) (ret int64, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int64(conn.Do("HINCRBY", hkey, key, val))
		return err
	}
	return ret, Doit(handle)
}

func HMGetStrings(hkey string, keys []string) ([]string, error) {
	return redis.Strings(HMGet(hkey, keys))
}

func HMGetBytes(hkey string, keys []string) ([][]byte, error) {
	return redis.ByteSlices(HMGet(hkey, keys))
}

func HMGet(hkey string, keys []string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{hkey}
		for _, key := range keys {
			args = append(args, key)
		}
		res, err = conn.Do("HMGET", args...)
		return err
	}
	return res, Doit(handle)
}

// vals:  [key1, val1, key2, val2, ...]
func HMSet(hkey string, vals []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{hkey}
		args = append(args, vals...)
		_, err = conn.Do("HMSET", args...)
		return err
	}
	return Doit(handle)
}

// 效率原因, 不建议使用
func HGetAll(hkey string) (res interface{}, err error) {
	handle := func(conn redis.Conn) error {
		res, err = conn.Do("HGETALL", hkey)
		return err
	}
	return res, Doit(handle)
}

func HExists(hkey, key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("HEXISTS", hkey, key))
		return err
	}
	return ret, Doit(handle) // 0:不存在 1:存在 (err != nil, ret == 0)
}

// --Set
func SAdd(skey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("SADD", skey, member)
		return err
	}
	return Doit(handle)
}

func SIsMember(skey, member string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("SISMEMBER", skey, member))
		return err
	}
	return ret, Doit(handle) // 0:不存在 1:存在 (err != nil, ret == 0)
}

func SRem(skey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("SREM", skey, member)
		return err
	}
	return Doit(handle)
}

// --Sorted Set
func ZAdd(zkey string, score int64, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZADD", zkey, score, member)
		return err
	}
	return Doit(handle)
}

// vals: [score1, member1, score2, member2, ...]
func ZMAdd(zkey string, vals []interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey}
		args = append(args, vals...)
		_, err = conn.Do("ZADD", args...)
		return err
	}
	return Doit(handle)
}

func ZScore(zkey string, member string) (score int64, err error) {
	handle := func(conn redis.Conn) error {
		score, err = redis.Int64(conn.Do("ZSCORE", zkey, member))
		return err
	}
	return score, Doit(handle)
}

func ZRangeByScore(zkey string, beg, end int64, limit int) (items []string, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey, beg, end}
		if limit > 0 {
			args = append(args, "limit", 0, limit)
		}
		items, err = redis.Strings(conn.Do("ZRANGEBYSCORE", args))
		return err
	}
	return items, Doit(handle)
}

func ZRevRangeByScore(zkey string, beg, end int64, limit int) (items []string, err error) {
	handle := func(conn redis.Conn) error {
		args := []interface{}{zkey, beg, end}
		if limit > 0 {
			args = append(args, "limit", 0, limit)
		}
		items, err = redis.Strings(conn.Do("ZREVRANGEBYSCORE", args))
		return err
	}
	return items, Doit(handle)
}

func ZRemRangeByScore(zkey string, beg, end int64) error {
	handle := func(conn redis.Conn) error {
		_, err := conn.Do("ZREMRANGEBYSCORE", zkey, beg, end)
		return err
	}
	return Doit(handle)
}

func ZRem(zkey, member string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("ZREM", zkey, member)
		return err
	}
	return Doit(handle)
}

// --List
func LPopString(lkey string) (val string, err error) {
	return redis.String(LPop(lkey))
}

func LPop(lkey string) (val interface{}, err error) {
	handle := func(conn redis.Conn) error {
		val, err = conn.Do("LPOP", lkey)
		return err
	}
	return val, Doit(handle)
}

func RPush(lkey string, val interface{}) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("RPUSH", lkey, val)
		return err
	}
	return Doit(handle)
}

// --keys
func Exist(key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("EXISTS", key))
		return err
	}
	return ret, Doit(handle) // 1:存在 0:不存在
}

func Expire(key string, ttl int) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("EXPIRE", key, ttl))
		return err
	}
	return ret, Doit(handle) // 1:成功 0:失败
}

func KeyTTL(key string) (ret int, err error) {
	handle := func(conn redis.Conn) error {
		ret, err = redis.Int(conn.Do("TTL", key))
		return err
	}
	return ret, Doit(handle) // -2:key不存在 -1:没有设置TTL num:剩余生存时间(秒)
}

func Delete(key string) (err error) {
	handle := func(conn redis.Conn) error {
		_, err = conn.Do("DEL", key)
		return err
	}
	return Doit(handle)
}

// --

func Doit(h func(redis.Conn) error) error {
	conn := gRedigo.Get()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()
	return h(conn)
}

var gRedigo *redis.Pool

func Init(conf configor.RedisConfigor) {
	fmt.Println("初始化Redis连接池")
	gRedigo = &redis.Pool{
		MaxIdle:   conf.MaxIdle,
		MaxActive: conf.MaxIdle * 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Host,
				redis.DialConnectTimeout(time.Duration(500)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(500)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(500)*time.Millisecond))
			if nil != err {
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
		},
	}
}
