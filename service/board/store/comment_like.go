package store

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/service/board/model"
)

func (d *db) ListUserCommLikes(ctx context.Context, uid string) (model.CommentLikes, error) {
	items := make(model.CommentLikes, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).Sort("-created_at").Limit(500).All(&items)
	}
	return items, d.Pool.Doit(ctx, "BoardComment", "like", handle)
}

func (d *db) NewCommLike(ctx context.Context, pitem *model.CommentLike) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return d.Pool.Doit(ctx, "BoardComment", "like", handle)
}

func (d *db) DelCommLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return d.Pool.Doit(ctx, "BoardComment", "like", handle)
}
