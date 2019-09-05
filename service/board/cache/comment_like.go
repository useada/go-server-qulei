package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"a.com/go-server/service/board/model"
)

func (k *kv) ListUserCommLikes(ctx context.Context, uid string) (model.CommentLikes, error) {
	items := model.CommentLikes{}
	data, err := k.Pool.GetBytes(ctx, k.genUserCommLikeKey(uid))
	if err != nil {
		return items, err
	}
	return items, json.Unmarshal(data, &items)
}

func (k *kv) NewUserCommLikes(ctx context.Context, uid string, items model.CommentLikes) error {
	if len(items) == 0 {
		return nil
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return k.Pool.Set(ctx, k.genUserCommLikeKey(uid), data, 3600*24)
}

func (k *kv) DelUserCommLikes(ctx context.Context, uid string) error {
	return k.Pool.Delete(ctx, k.genUserCommLikeKey(uid))
}

func (k *kv) genUserCommLikeKey(uid string) string {
	return fmt.Sprintf("BDUCL|%s", uid)
}
