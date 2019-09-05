package store

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/service/board/model"
)

func (d *db) GetSummaries(ctx context.Context, ids []string) (model.Summaries, error) {
	items := make(model.Summaries, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) GetSummary(ctx context.Context, id string) (*model.Summary, error) {
	pitem := &model.Summary{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) IncrSummaryComm(ctx context.Context, cid, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": 1}}
		if len(cid) == 0 {
			data["$inc"] = bson.M{"comms_count": 1, "comms_first_count": 1}
		}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) DecrSummaryComm(ctx context.Context, cid, id string, count int) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": -count}}
		if len(cid) == 0 {
			data["$inc"] = bson.M{"comms_count": -count, "comms_first_count": -1}
		}
		return c.Update(bson.M{"_id": id}, data)
	}
	return d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) IncrSummaryLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) DecrSummaryLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) IncrSummaryRepost(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}

func (d *db) DecrSummaryRepost(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return d.Pool.Doit(ctx, "BoardSummary", "summary", handle)
}
