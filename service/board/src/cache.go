package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"a.com/go-server/common/redis"
	"a.com/go-server/proto/constant"
)

type CacheHandle struct {
}

var Cache = &CacheHandle{}

func (c *CacheHandle) InitComms(ctx context.Context, oid, cid string, items CommentModels, total bool) error {
	zsetArgs := make([]interface{}, 0)
	hashArgs := make([]interface{}, 0)
	for _, item := range items {
		if data, err := json.Marshal(item); err == nil {
			hashArgs = append(hashArgs, item.ID)
			hashArgs = append(hashArgs, data)
		}
		zsetArgs = append(zsetArgs, item.CreatedAt)
		zsetArgs = append(zsetArgs, item.ID)
	}
	if total { // 全部缓存
		zsetArgs = append(zsetArgs, 0)
		zsetArgs = append(zsetArgs, "GUARD")
	}

	var firststamp int64 // 最旧的一条数据的时间戳
	if len(items) > 0 {
		firststamp = items[0].CreatedAt
	}
	c.InitHashComms(ctx, oid, hashArgs)
	return c.InitZsetComms(ctx, oid, cid, firststamp, zsetArgs)
}

func (c *CacheHandle) PushComment(ctx context.Context, pitem *CommentModel) error {
	c.SetHashComm(ctx, pitem)
	return c.PushZsetComm(ctx, pitem)
}

func (c *CacheHandle) PopComment(ctx context.Context, oid, cid, id string) error {
	c.DelHashComm(ctx, oid, id)
	return c.PopZsetComm(ctx, oid, cid, id)
}

func (c *CacheHandle) MutiGetSummary(ctx context.Context, oids []string) (SummaryModels, error) {
	keys := make([]string, 0)
	for _, oid := range oids {
		keys = append(keys, c.KeySummary(oid))
	}
	vals, err := redis.MGetBytes(ctx, keys)
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

func (c *CacheHandle) NewSummary(ctx context.Context, pitem *SummaryModel) error {
	data, err := json.Marshal(pitem)
	if err != nil {
		return err
	}
	return redis.Set(ctx, c.KeySummary(pitem.ID), data, 3600)
}

func (c *CacheHandle) DelSummary(ctx context.Context, oid string) error {
	return redis.Delete(ctx, c.KeySummary(oid))
}

func (c *CacheHandle) ListUserCommLikes(ctx context.Context, uid string) (CommentLikeModels, error) {
	items := CommentLikeModels{}
	data, err := redis.GetBytes(ctx, c.KeyUserCommLikes(uid))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (c *CacheHandle) NewUserCommLikes(ctx context.Context, uid string, items CommentLikeModels) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return redis.Set(ctx, c.KeyUserCommLikes(uid), data, 3600*24)
}

func (c *CacheHandle) DelUserCommLikes(ctx context.Context, uid string) error {
	return redis.Delete(ctx, c.KeyUserCommLikes(uid))
}

func (c *CacheHandle) ListUserLikes(ctx context.Context, uid string) (LikeModels, error) {
	items := LikeModels{}
	data, err := redis.GetBytes(ctx, c.KeyUserLikes(uid))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (c *CacheHandle) NewUserLikes(ctx context.Context, uid string, items LikeModels) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return redis.Set(ctx, c.KeyUserLikes(uid), data, 3600*24)
}

func (c *CacheHandle) DelUserLikes(ctx context.Context, uid string) error {
	return redis.Delete(ctx, c.KeyUserLikes(uid))
}

// -- Hash

func (c *CacheHandle) InitHashComms(ctx context.Context, oid string, vals []interface{}) error {
	hkey := c.KeyHashComms(oid)
	err := redis.HMSet(ctx, hkey, vals)
	if err == nil {
		redis.Expire(ctx, hkey, TTL_HASH_KEY)
	}
	return err
}

func (c *CacheHandle) GetHashComm(ctx context.Context, oid, id string) (*CommentModel, error) {
	pitem := &CommentModel{}
	val, err := redis.HGetBytes(ctx, c.KeyHashComms(oid), id)
	if err != nil {
		return pitem, err
	}
	return pitem, json.Unmarshal(val, pitem)
}

func (c *CacheHandle) MutiGetHashComms(ctx context.Context, oid string, ids []string) (CommentModels, error) {
	vals, err := redis.HMGetBytes(ctx, c.KeyHashComms(oid), ids)
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

func (c *CacheHandle) SetHashComm(ctx context.Context, pitem *CommentModel) error {
	val, err := json.Marshal(pitem)
	if err != nil {
		return err
	}

	hkey := c.KeyHashComms(pitem.Oid)
	if err = redis.HSet(ctx, hkey, pitem.ID, val); err == nil {
		redis.Expire(ctx, hkey, TTL_HASH_KEY)
	}
	return err
}

func (c *CacheHandle) DelHashComm(ctx context.Context, oid, id string) error {
	if len(oid) == 0 || len(id) == 0 {
		return nil
	}
	return redis.HDel(ctx, c.KeyHashComms(oid), id)
}

// -- ZSET

func (c *CacheHandle) InitZsetComms(ctx context.Context, oid, cid string, stamp int64, vals []interface{}) error {
	zkey := c.KeyZsetComms(oid + cid)
	redis.ZRemByScore(ctx, zkey, 0, stamp)

	err := redis.ZMAdd(ctx, zkey, vals)
	if err != nil {
		redis.Delete(ctx, zkey)
	} else {
		redis.Expire(ctx, zkey, TTL_ZSET_KEY)
	}
	return err
}

func (c *CacheHandle) ListZsetComms(ctx context.Context, oid, cid string, stamp int64, limit int) ([]string, error) {
	zkey := c.KeyZsetComms(oid + cid)
	if ok := c.CheckZsetCommsKey(ctx, zkey); !ok {
		return nil, errors.New("zset key ttl failed")
	}

	ids, err := redis.ZRevRangeByScore(ctx, zkey, stamp-1, constant.TIME_INF_MIN, limit)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 { // 如果曾经cache过且没有过期，至少会包含GUARD
		return nil, errors.New("cache zset empty")
	}

	if len(ids) < limit && ids[len(ids)-1] != "GUARD" {
		return nil, errors.New("some data in database")
	}

	if ids[len(ids)-1] == "GUARD" {
		return ids[0 : len(ids)-1], nil
	}
	return ids, nil
}

func (c *CacheHandle) PushZsetComm(ctx context.Context, pitem *CommentModel) error {
	zkey := c.KeyZsetComms(pitem.Oid + pitem.Cid)
	err := redis.ZAdd(ctx, zkey, pitem.CreatedAt, pitem.ID)
	if err != nil {
		redis.Delete(ctx, zkey)
	} else {
		redis.Expire(ctx, zkey, TTL_ZSET_KEY)
	}
	return err
}

func (c *CacheHandle) PopZsetComm(ctx context.Context, oid, cid, id string) error {
	zkey := c.KeyZsetComms(oid + cid)
	if ok := c.CheckZsetCommsKey(ctx, zkey); !ok {
		return errors.New("check zset key ttl error")
	}

	err := redis.ZRem(ctx, zkey, id)
	if err != nil {
		redis.Delete(ctx, zkey)
	} else {
		redis.Expire(ctx, zkey, TTL_ZSET_KEY)
	}
	return err
}

func (c *CacheHandle) CheckZsetCommsKey(ctx context.Context, zkey string) bool {
	val, err := redis.KeyTTL(ctx, zkey)
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
