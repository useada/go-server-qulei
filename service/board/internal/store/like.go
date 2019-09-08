package store

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/service/board/internal/model"
)

func (d *db) ListLikes(ctx context.Context, oid string, stamp int64, limit int) (model.Likes, error) {
	items := make(model.Likes, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "created_at": bson.M{"$lt": stamp}}
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, d.Pool.Doit(ctx, "BoardLike", "like", handle)
}

func (d *db) ListUserLikes(ctx context.Context, uid string) (model.Likes, error) {
	items := make(model.Likes, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"uid": uid}
		return c.Find(query).Sort("-created_at").Limit(500).All(&items)
	}
	return items, d.Pool.Doit(ctx, "BoardLike", "like", handle)
}

func (d *db) GetLike(ctx context.Context, id string) (*model.Like, error) {
	pitem := &model.Like{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, d.Pool.Doit(ctx, "BoardLike", "like", handle)
}

func (d *db) NewLike(ctx context.Context, pitem *model.Like) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return d.Pool.Doit(ctx, "BoardLike", "like", handle)
}

func (d *db) DelLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return d.Pool.Doit(ctx, "BoardLike", "like", handle)
}
