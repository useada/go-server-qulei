package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"a.com/go-server/service/board/model"
)

func (k *kv) ListUserCommLikes(ctx context.Context, id string) (model.CommentLikes, error) {
	items := model.CommentLikes{}
	data, err := k.Pool.GetBytes(ctx, k.genUserCommLikeKey(id))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (k *kv) NewUserCommLikes(ctx context.Context, id string, items model.CommentLikes) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return k.Pool.Set(ctx, k.genUserCommLikeKey(id), data, 3600*24)
}

func (k *kv) DelUserCommLikes(ctx context.Context, id string) error {
	return k.Pool.Delete(ctx, k.genUserCommLikeKey(id))
}

func (k *kv) genUserCommLikeKey(id string) string {
	return fmt.Sprintf("BDUCL|%s", id)
}
