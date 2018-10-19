package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	configor "a.com/go-server/common/configor"
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

func Get(key string) (interface{}, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return conn.Do("GET", key)
}

func Set(key string, val interface{}) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	return err
}

func SetTTL(key string, val interface{}, ttl int64) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("SET", key, val, "EX", ttl)
	return err
}

func Incr(key string, val int) (int, error) {
	conn := GetConn()
	if nil == conn {
		return 0, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Int(conn.Do("INCRBY", key, val))
}

func Decr(key string, val int) (int, error) {
	conn := GetConn()
	if nil == conn {
		return 0, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Int(conn.Do("DECRBY", key, val))
}

func MGetBytes(keys []string) ([][]byte, error) {
	return redis.ByteSlices(MGet(keys))
}

func MGet(keys []string) (interface{}, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	return conn.Do("MGET", args...)
}

func MSetBytes(keys []string, vals [][]byte) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	for idx, key := range keys {
		conn.Send("SET", key, vals[idx])
	}
	if err := conn.Flush(); err != nil {
		return err
	}

	for range keys {
		conn.Receive()
	}
	return nil
}

func MSetBytesTTL(keys []string, vals [][]byte, ttl int64) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	for idx, key := range keys {
		conn.Send("SET", key, vals[idx], "EX", ttl)
	}
	if err := conn.Flush(); err != nil {
		return err
	}

	for range keys {
		conn.Receive()
	}
	return nil
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

func HGet(hkey string, key string) (interface{}, error) {
	conn := GetConn()
	if nil == conn {
		return "", errors.New("get redis conn failed")
	}
	defer conn.Close()

	return conn.Do("HGET", hkey, key)
}

func HSet(hkey string, key string, val interface{}) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("HSET", hkey, key, val)
	return err
}

func HDel(hkey, key string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("HDEL", hkey, key)
	return err
}

func HIncr(hkey string, key string, val int) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("HINCRBY", hkey, key, val)
	return err
}

func HMGetStrings(hkey string, keys []string) ([]string, error) {
	return redis.Strings(HMGet(hkey, keys))
}

func HMGetBytes(hkey string, keys []string) ([][]byte, error) {
	return redis.ByteSlices(HMGet(hkey, keys))
}

func HMGet(hkey string, keys []string) (interface{}, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	args := []interface{}{hkey}
	for _, key := range keys {
		args = append(args, key)
	}
	return conn.Do("HMGET", args...)
}

// vals:  [key1 val1 key2 val2 ...]
func HMSet(hkey string, vals []interface{}) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	args := []interface{}{hkey}
	args = append(args, vals...)

	_, err := conn.Do("HMSET", args...)
	return err
}

func HGetAll(hkey string) (interface{}, error) {
	conn := GetConn()
	if nil == conn {
		return 0, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return conn.Do("HGETALL", hkey)
}

func HExists(hkey, key string) (bool, error) {
	conn := GetConn()
	if nil == conn {
		return false, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Bool(conn.Do("HEXISTS", hkey, key))
}

// --Set
func SAdd(skey, member string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("SADD", skey, member)
	return err
}

func SExist(skey, member string) bool {
	conn := GetConn()
	if nil == conn {
		return false
	}
	defer conn.Close()

	ret, _ := redis.Int(conn.Do("SISMEMBER", skey, member))
	return ret > 0
}

func SRem(skey, member string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("SREM", skey, member)
	return err
}

// --Sorted Set
func ZAdd(zkey string, score int64, member string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("ZADD", zkey, score, member)
	return err
}

// val: [score1, member1, score2, member2 ...]
func ZMAdd(zkey string, vals []interface{}) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	args := []interface{}{zkey}
	args = append(args, vals...)

	_, err := conn.Do("ZADD", args...)
	return err
}

func ZScore(zkey string, member string) (int64, error) {
	conn := GetConn()
	if nil == conn {
		return 0, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Int64(conn.Do("ZSCORE", zkey, member))
}

func ZRangeByScore(zkey string, beg, end int64) ([]string, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Strings(conn.Do("ZRANGEBYSCORE", zkey, beg, end))
}

func ZRangeByScoreLimit(zkey string, beg, end int64, limit int) ([]string, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Strings(conn.Do("ZRANGEBYSCORE",
		zkey, beg, end, "limit", 0, limit))
}

func ZRevRangeByScore(zkey string, beg, end int64) ([]string, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Strings(conn.Do("ZRANGEBYSCORE", zkey, beg, end))
}

func ZRevRangeByScoreLimit(zkey string, beg, end int64, limit int) ([]string, error) {
	conn := GetConn()
	if nil == conn {
		return nil, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Strings(conn.Do("ZREVRANGEBYSCORE",
		zkey, beg, end, "limit", 0, limit))
}

func ZRemRangeByScore(zkey string, beg, end int64) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("ZREMRANGEBYSCORE", zkey, beg, end)
	return err
}

func ZRem(zkey, member string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("ZREM", zkey, member)
	return err
}

// --List
func LPop(lkey string) (string, error) {
	conn := GetConn()
	if nil == conn {
		return "", errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.String(conn.Do("LPOP", lkey))
}

func RPush(lkey string, val string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("RPUSH", lkey, val)
	return err
}

// --keys
func Exist(key string) (int, error) {
	conn := GetConn()
	if nil == conn {
		return -1, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Int(conn.Do("EXISTS", key))
}

func Expire(key string, ttl int) (int, error) {
	conn := GetConn()
	if nil == conn {
		return -1, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Int(conn.Do("EXPIRE", key, ttl))
}

func KeyTTL(key string) (int, error) {
	conn := GetConn()
	if nil == conn {
		return -1, errors.New("get redis conn failed")
	}
	defer conn.Close()

	return redis.Int(conn.Do("TTL", key))
}

func Delete(key string) error {
	conn := GetConn()
	if nil == conn {
		return errors.New("get redis conn failed")
	}
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// --
func GetConn() redis.Conn {
	return gRedisPool.Get()
}

var gRedisPool *redis.Pool

func InitRedis(conf configor.RedisConfigor) {
	fmt.Println("初始化Redis连接池")
	gRedisPool = &redis.Pool{
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
