package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"a.com/go-server/service/board/internal/model"
)

func (k *kv) ListUserLikes(ctx context.Context, uid string) (model.Likes, error) {
	items := model.Likes{}
	data, err := k.Pool.GetBytes(ctx, k.genUserLikeKey(uid))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (k *kv) NewUserLikes(ctx context.Context, uid string, items model.Likes) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return k.Pool.Set(ctx, k.genUserLikeKey(uid), data, 3600*24)
}

func (k *kv) DelUserLikes(ctx context.Context, uid string) error {
	return k.Pool.Delete(ctx, k.genUserLikeKey(uid))
}

func (k *kv) genUserLikeKey(uid string) string {
	return fmt.Sprintf("BDUL|%s", uid)
}
