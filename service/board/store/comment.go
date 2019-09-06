package store

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/service/board/model"
)

func (d *db) ListComments(ctx context.Context, oid, cid string, stamp int64, limit int) (model.Comments, error) {
	items := make(model.Comments, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "cid": cid, "created_at": bson.M{"$lt": stamp}}
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) GetComments(ctx context.Context, ids []string) (model.Comments, error) {
	items := make(model.Comments, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) GetComment(ctx context.Context, id string) (*model.Comment, error) {
	pitem := &model.Comment{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) NewComment(ctx context.Context, pitem *model.Comment) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) DelComment(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) IncrCommReply(ctx context.Context, cid string, pitem *model.Comment) error {
	handle := func(c *mgo.Collection) error {
		replys := bson.M{"$each": model.Comments{*pitem}, "$sort": bson.M{"created_at": -1}, "$slice": 2}
		return c.Update(bson.M{"_id": cid},
			bson.M{"$push": bson.M{"replys": replys}, "$inc": bson.M{"reply_count": 1}})
	}
	return d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) DecrCommReply(ctx context.Context, cid, rid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$pull": bson.M{"replys": bson.M{"_id": rid}}, "$inc": bson.M{"reply_count": -1}})
	}
	return d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) IncrCommLike(ctx context.Context, cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": 1}})
	}
	return d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}

func (d *db) DecrCommLike(ctx context.Context, cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": -1}})
	}
	return d.Pool.Doit(ctx, "BoardComment", "comment", handle)
}
