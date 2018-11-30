package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"a.com/go-server/common/redis"
)

type CacheHandle struct {
}

var Cache = &CacheHandle{}

func (c *CacheHandle) InitComms(oid, cid string,
	items CommentModels, total bool) error {
	zsetArgs := make([]interface{}, 0)
	hashArgs := make([]interface{}, 0)
	for _, item := range items {
		if data, err := json.Marshal(item); err == nil {
			hashArgs = append(hashArgs, item.Id)
			hashArgs = append(hashArgs, data)
		}
		zsetArgs = append(zsetArgs, item.CreatedAt)
		zsetArgs = append(zsetArgs, item.Id)
	}
	if total { // 全部缓存
		zsetArgs = append(zsetArgs, 0)
		zsetArgs = append(zsetArgs, "GUARD")
	}

	var firststamp int64 // 最旧的一条数据的时间戳
	if len(items) > 0 {
		firststamp = items[0].CreatedAt
	}
	c.InitHashComms(oid, hashArgs)
	return c.InitZsetComms(oid, cid, firststamp, zsetArgs)
}

func (c *CacheHandle) PushComment(pitem *CommentModel) error {
	c.SetHashComm(pitem)
	return c.PushZsetComm(pitem)
}

func (c *CacheHandle) PopComment(oid, cid, id string) error {
	c.DelHashComm(oid, id)
	return c.PopZsetComm(oid, cid, id)
}

func (c *CacheHandle) MutiGetSummary(oids []string) (SummaryModels, error) {
	keys := make([]string, 0)
	for _, oid := range oids {
		keys = append(keys, c.KeySummary(oid))
	}
	vals, err := redis.MGetBytes(keys)
	if err != nil {
		return nil, err
	}

	items := make(SummaryModels, 0)
	for _, val := range vals {
		item := SummaryModel{}
		if err = json.Unmarshal(val, &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, err
}

func (c *CacheHandle) NewSummary(pitem *SummaryModel) error {
	data, err := json.Marshal(pitem)
	if err != nil {
		return err
	}
	return redis.Set(c.KeySummary(pitem.Id), data, 3600)
}

func (c *CacheHandle) DelSummary(oid string) error {
	return redis.Delete(c.KeySummary(oid))
}

func (c *CacheHandle) ListUserCommLikes(uid string) (CommentLikeModels, error) {
	items := CommentLikeModels{}
	data, err := redis.GetBytes(c.KeyUserCommLikes(uid))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (c *CacheHandle) NewUserCommLikes(uid string,
	items CommentLikeModels) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return redis.Set(c.KeyUserCommLikes(uid), data, 3600*24)
}

func (c *CacheHandle) DelUserCommLikes(uid string) error {
	return redis.Delete(c.KeyUserCommLikes(uid))
}

func (c *CacheHandle) ListUserLikes(uid string) (LikeModels, error) {
	items := LikeModels{}
	data, err := redis.GetBytes(c.KeyUserLikes(uid))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (c *CacheHandle) NewUserLikes(uid string, items LikeModels) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return redis.Set(c.KeyUserLikes(uid), data, 3600*24)
}

func (c *CacheHandle) DelUserLikes(uid string) error {
	return redis.Delete(c.KeyUserLikes(uid))
}

// -- Hash

func (c *CacheHandle) InitHashComms(oid string, vals []interface{}) error {
	hkey := c.KeyHashComms(oid)
	err := redis.HMSet(hkey, vals)
	if err == nil {
		redis.Expire(hkey, TTL_HASH_KEY)
	}
	return err
}

func (c *CacheHandle) GetHashComm(oid, id string) (*CommentModel, error) {
	pitem := &CommentModel{}
	val, err := redis.HGetBytes(c.KeyHashComms(oid), id)
	if err != nil {
		return pitem, err
	}
	return pitem, json.Unmarshal(val, pitem)
}

func (c *CacheHandle) MutiGetHashComms(oid string,
	ids []string) (CommentModels, error) {
	vals, err := redis.HMGetBytes(c.KeyHashComms(oid), ids)
	if err != nil {
		return nil, err
	}

	items := make(CommentModels, 0)
	for _, val := range vals {
		item := CommentModel{}
		if err := json.Unmarshal(val, &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (c *CacheHandle) SetHashComm(pitem *CommentModel) error {
	val, err := json.Marshal(pitem)
	if err != nil {
		return err
	}

	hkey := c.KeyHashComms(pitem.Oid)
	if err = redis.HSet(hkey, pitem.Id, val); err == nil {
		redis.Expire(hkey, TTL_HASH_KEY)
	}
	return err
}

func (c *CacheHandle) DelHashComm(oid, id string) error {
	if len(oid) == 0 || len(id) == 0 {
		return nil
	}
	return redis.HDel(c.KeyHashComms(oid), id)
}

// -- ZSET

func (c *CacheHandle) InitZsetComms(oid, cid string,
	stamp int64, vals []interface{}) error {
	zkey := c.KeyZsetComms(oid + cid)
	redis.ZRemByScore(zkey, 0, stamp)

	err := redis.ZMAdd(zkey, vals)
	if err != nil {
		redis.Delete(zkey)
	} else {
		redis.Expire(zkey, TTL_ZSET_KEY)
	}
	return err
}

func (c *CacheHandle) ListZsetComms(oid, cid, direction string,
	stamp int64, limit int) (ids []string, err error) {
	zkey := c.KeyZsetComms(oid + cid)
	if ok := c.CheckZsetCommsKey(zkey); !ok {
		return nil, errors.New("zset key ttl failed")
	}

	if count, _ := redis.ZCard(zkey); count == 0 {
		return nil, errors.New("cache empty")
	}

	if direction == "gt" || direction == "gte" {
		return c.ListZsetCommsNewer(zkey, stamp+1, TIMESTAMP_INF, limit)
	}
	return c.ListZsetCommsOlder(zkey, stamp-1, 0, limit)
}

func (c *CacheHandle) ListZsetCommsNewer(zkey string,
	beg, end int64, limit int) ([]string, error) {
	ids, err := redis.ZRangeByScore(zkey, beg+1, end, limit)
	if err != nil {
		return ids, err
	}

	if ids[0] == "GUARD" {
		return ids[1:], nil
	}
	return ids, err
}

func (c *CacheHandle) ListZsetCommsOlder(zkey string,
	end, beg int64, limit int) ([]string, error) {
	if end == 0 {
		end = TIMESTAMP_INF
	}
	ids, err := redis.ZRevRangeByScore(zkey, end-1, beg, limit)
	if err != nil {
		return ids, err
	}

	for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
		ids[i], ids[j] = ids[j], ids[i]
	}

	if len(ids) < limit && ids[0] != "GUARD" {
		return ids, errors.New("some data in database")
	}
	if ids[0] == "GUARD" {
		return ids[1:], nil
	}
	return ids, err
}

func (c *CacheHandle) PushZsetComm(pitem *CommentModel) error {
	zkey := c.KeyZsetComms(pitem.Oid + pitem.Cid)
	err := redis.ZAdd(zkey, pitem.CreatedAt, pitem.Id)
	if err != nil {
		redis.Delete(zkey)
	} else {
		redis.Expire(zkey, TTL_ZSET_KEY)
	}
	return err
}

func (c *CacheHandle) PopZsetComm(oid, cid, id string) error {
	zkey := c.KeyZsetComms(oid + cid)
	if ok := c.CheckZsetCommsKey(zkey); !ok {
		return errors.New("check zset key ttl error")
	}

	err := redis.ZRem(zkey, id)
	if err != nil {
		redis.Delete(zkey)
	} else {
		redis.Expire(zkey, TTL_ZSET_KEY)
	}
	return err
}

func (c *CacheHandle) CheckZsetCommsKey(zkey string) bool {
	val, err := redis.KeyTTL(zkey)
	if err != nil {
		return false
	}

	if val == -1 { // val: -1 永不过期
		return true
	}
	if val <= TTL_ZSET_CRITICAL { // 即将超时 / val: -2 过期或不存在
		return false
	}
	return true
}

// -- KEY

func (c *CacheHandle) KeySummary(oid string) string {
	return fmt.Sprintf("BDSUM|%s", oid)
}

func (c *CacheHandle) KeyHashComms(oid string) string {
	return fmt.Sprintf("BDCM|HASH|%s", oid)
}

func (c *CacheHandle) KeyZsetComms(oid string) string {
	return fmt.Sprintf("BDCM|ZSET|%s", oid)
}

func (c *CacheHandle) KeyUserCommLikes(uid string) string {
	return fmt.Sprintf("BDUCL|%s", uid)
}

func (c *CacheHandle) KeyUserLikes(uid string) string {
	return fmt.Sprintf("BDUL|%s", uid)
}
