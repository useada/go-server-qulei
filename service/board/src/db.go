package main

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/common/mongo"
)

type DbHandle struct {
}

var DB = &DbHandle{}

// -- Comment

func (db *DbHandle) ListComments(ctx context.Context, oid, cid string, stamp int64, limit int) (CommentModels, error) {
	items := make(CommentModels, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "cid": cid, "created_at": bson.M{"$lt": stamp}}
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) MutiGetComments(ctx context.Context, ids []string) (CommentModels, error) {
	items := make(CommentModels, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) GetComment(ctx context.Context, id string) (*CommentModel, error) {
	pitem := &CommentModel{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) NewComment(ctx context.Context, pitem *CommentModel) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) DelComment(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) IncrCommReply(ctx context.Context, cid string, pitem *CommentModel) error {
	handle := func(c *mgo.Collection) error {
		replys := bson.M{"$each": CommentModels{*pitem},
			"$sort": bson.M{"created_at": -1}, "$slice": 2}
		return c.Update(bson.M{"_id": cid},
			bson.M{"$push": bson.M{"replys": replys},
				"$inc": bson.M{"reply_count": 1}})
	}
	return mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) DecrCommReply(ctx context.Context, cid, rid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$pull": bson.M{"replys": bson.M{"_id": rid}},
				"$inc": bson.M{"reply_count": -1}})
	}
	return mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) IncrCommLike(ctx context.Context, cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": 1}})
	}
	return mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) DecrCommLike(ctx context.Context, cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": -1}})
	}
	return mongo.Doit(ctx, "BoardComment", "comment", handle)
}

func (db *DbHandle) ListUserCommLikes(ctx context.Context, uid string) (CommentLikeModels, error) {
	items := make(CommentLikeModels, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).Sort("-created_at").Limit(500).All(&items)
	}
	return items, mongo.Doit(ctx, "BoardComment", "like", handle)
}

func (db *DbHandle) NewCommLike(ctx context.Context, pitem *CommentLikeModel) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return mongo.Doit(ctx, "BoardComment", "like", handle)
}

func (db *DbHandle) DelCommLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit(ctx, "BoardComment", "like", handle)
}

// -- Like

func (db *DbHandle) ListLikes(ctx context.Context, oid string, stamp int64, limit int) (LikeModels, error) {
	items := make(LikeModels, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "created_at": bson.M{"$lt": stamp}}
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit(ctx, "BoardLike", "like", handle)
}

func (db *DbHandle) ListUserLikes(ctx context.Context, uid string) (LikeModels, error) {
	items := make(LikeModels, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"uid": uid}
		return c.Find(query).Sort("-created_at").Limit(500).All(&items)
	}
	return items, mongo.Doit(ctx, "BoardLike", "like", handle)
}

func (db *DbHandle) GetLike(ctx context.Context, id string) (*LikeModel, error) {
	pitem := &LikeModel{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, mongo.Doit(ctx, "BoardLike", "like", handle)
}

func (db *DbHandle) NewLike(ctx context.Context, pitem *LikeModel) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return mongo.Doit(ctx, "BoardLike", "like", handle)
}

func (db *DbHandle) DelLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit(ctx, "BoardLike", "like", handle)
}

// -- Summary

func (db *DbHandle) MutiGetSummary(ctx context.Context, ids []string) (SummaryModels, error) {
	items := make(SummaryModels, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) GetSummary(ctx context.Context, id string) (*SummaryModel, error) {
	pitem := &SummaryModel{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) IncrSummaryComm(ctx context.Context, cid, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": 1}}
		if len(cid) == 0 {
			data["$inc"] = bson.M{"comms_count": 1, "comms_first_count": 1}
		}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) DecrSummaryComm(ctx context.Context, cid, id string, count int) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": -count}}
		if len(cid) == 0 {
			data["$inc"] = bson.M{"comms_count": -count, "comms_first_count": -1}
		}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) IncrSummaryLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) DecrSummaryLike(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) IncrSummaryRepost(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit(ctx, "BoardSummary", "summary", handle)
}

func (db *DbHandle) DecrSummaryRepost(ctx context.Context, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit(ctx, "BoardSummary", "summary", handle)
}
