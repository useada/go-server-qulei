package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"a.com/go-server/proto/constant"

	"a.com/go-server/service/board/model"
)

const (
	TTL_ZSET_CRITICAL = 1            // 2 * (redis read/write timeout 500ms)
	TTL_ZSET_KEY      = 3600 * 9     //
	TTL_HASH_KEY      = 3600 * 9 * 3 //
)

func (k *kv) InitComments(ctx context.Context, oid, cid string, items model.Comments, total bool) error {
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

	k.setComments(ctx, oid, hashArgs)

	if len(items) > 0 {
		return k.pushComments(ctx, oid, cid, items[0].CreatedAt, zsetArgs)
	}
	return k.pushComments(ctx, oid, cid, 0, zsetArgs)
}

func (k *kv) GetComments(ctx context.Context, oid string, ids []string) (model.Comments, error) {
	vals, err := k.Pool.HMGetBytes(ctx, k.genHashKey(oid), ids)
	if err != nil {
		return nil, err
	}

	items := make(model.Comments, 0)
	for _, val := range vals {
		item := model.Comment{}
		if err := json.Unmarshal(val, &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (k *kv) GetComment(ctx context.Context, oid, id string) (*model.Comment, error) {
	pitem := &model.Comment{}
	val, err := k.Pool.HGetBytes(ctx, k.genHashKey(oid), id)
	if err != nil {
		return pitem, err
	}
	return pitem, json.Unmarshal(val, pitem)
}

func (k *kv) SetComment(ctx context.Context, pitem *model.Comment) error {
	val, err := json.Marshal(pitem)
	if err != nil {
		return err
	}

	hkey := k.genHashKey(pitem.Oid)
	if err = k.Pool.HSet(ctx, hkey, pitem.ID, val); err == nil {
		k.Pool.Expire(ctx, hkey, TTL_HASH_KEY)
	}
	return err
}

func (k *kv) DelComment(ctx context.Context, oid, id string) error {
	if len(oid) == 0 || len(id) == 0 {
		return nil
	}
	return k.Pool.HDel(ctx, k.genHashKey(oid), id)
}

func (k *kv) RangeComments(ctx context.Context, oid, cid string, stamp int64, limit int) ([]string, error) {
	zkey := k.genZsetKey(oid + cid)
	if ok := k.ttlCommentKey(ctx, zkey); !ok {
		return nil, errors.New("zset key ttl failed")
	}

	ids, err := k.Pool.ZRevRangeByScore(ctx, zkey, stamp-1, constant.TIME_INF_MIN, limit)
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

func (k *kv) PushComment(ctx context.Context, pitem *model.Comment) error {
	k.SetComment(ctx, pitem)
	return k.pushComment(ctx, pitem)
}

func (k *kv) PopComment(ctx context.Context, oid, cid, id string) error {
	k.DelComment(ctx, oid, id)
	return k.popComment(ctx, oid, cid, id)
}

func (k *kv) setComments(ctx context.Context, oid string, vals []interface{}) error {
	hkey := k.genHashKey(oid)
	err := k.Pool.HMSet(ctx, hkey, vals)
	if err == nil {
		k.Pool.Expire(ctx, hkey, TTL_HASH_KEY)
	}
	return err
}

func (k *kv) pushComments(ctx context.Context, oid, cid string, stamp int64, vals []interface{}) error {
	zkey := k.genZsetKey(oid + cid)
	k.Pool.ZRemByScore(ctx, zkey, 0, stamp)

	err := k.Pool.ZMAdd(ctx, zkey, vals)
	if err != nil {
		k.Pool.Delete(ctx, zkey)
	} else {
		k.Pool.Expire(ctx, zkey, TTL_ZSET_KEY)
	}
	return err
}

func (k *kv) pushComment(ctx context.Context, pitem *model.Comment) error {
	zkey := k.genZsetKey(pitem.Oid + pitem.Cid)
	err := k.Pool.ZAdd(ctx, zkey, pitem.CreatedAt, pitem.ID)
	if err != nil {
		k.Pool.Delete(ctx, zkey)
	} else {
		k.Pool.Expire(ctx, zkey, TTL_ZSET_KEY)
	}
	return err
}

func (k *kv) popComment(ctx context.Context, oid, cid, id string) error {
	zkey := k.genZsetKey(oid + cid)
	if ok := k.ttlCommentKey(ctx, zkey); !ok {
		return errors.New("check zset key ttl error")
	}

	err := k.Pool.ZRem(ctx, zkey, id)
	if err != nil {
		k.Pool.Delete(ctx, zkey)
	} else {
		k.Pool.Expire(ctx, zkey, TTL_ZSET_KEY)
	}
	return err
}

func (k *kv) ttlCommentKey(ctx context.Context, zkey string) bool {
	val, err := k.Pool.KeyTTL(ctx, zkey)
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

func (k *kv) genHashKey(oid string) string {
	return fmt.Sprintf("BDCM|HASH|%s", oid)
}

func (k *kv) genZsetKey(oid string) string {
	return fmt.Sprintf("BDCM|ZSET|%s", oid)
}
