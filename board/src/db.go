package main

import (
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/common/mongo"
)

type DbHandle struct {
}

var DB = &DbHandle{}

// -- Comment

func (db *DbHandle) ListComments(oid, cid, direction string, stamp int64,
	limit int) (CommentModels, error) {
	items := make(CommentModels, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "cid": cid}
		if len(direction) != 0 {
			query["created_at"] = bson.M{"$" + direction: stamp} //$gt, $lt
		}

		if direction == "gt" || direction == "gte" {
			return c.Find(query).Sort("created_at").Limit(limit).All(&items)
		}

		defer sort.Sort(items) // created_at升序返回
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) MutiGetComments(ids []string) (CommentModels, error) {
	items := make(CommentModels, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) GetComment(id string) (*CommentModel, error) {
	pitem := &CommentModel{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) NewComment(pitem *CommentModel) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) DelComment(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) IncrCommReply(cid string, pitem *CommentModel) error {
	if len(cid) == 0 {
		return nil
	}
	handle := func(c *mgo.Collection) error {
		replys := bson.M{"$each": CommentModels{*pitem},
			"$sort": "-created_at", "$slice": 3}
		return c.Update(bson.M{"_id": cid},
			bson.M{"$push": bson.M{"replys": replys},
				"$inc": bson.M{"reply_count": 1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) DecrCommReply(cid, rid string) error {
	if len(cid) == 0 {
		return nil
	}
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$pull": bson.M{"replys": bson.M{"_id": rid}},
				"$inc": bson.M{"reply_count": -1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) IncrCommLike(cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": 1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) DecrCommLike(cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": -1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DbHandle) ListUserCommLikes(uid string) (CommentLikeModels, error) {
	items := make(CommentLikeModels, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).Sort("-created_at").Limit(500).All(&items)
	}
	return items, mongo.Doit("BoardComment", "like", handle)
}

func (db *DbHandle) NewCommLike(pitem *CommentLikeModel) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return mongo.Doit("BoardComment", "like", handle)
}

func (db *DbHandle) DelCommLike(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit("BoardComment", "like", handle)
}

// -- Like

func (db *DbHandle) ListLikes(oid string, stamp int64,
	limit int) (LikeModels, error) {
	items := make(LikeModels, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid}
		if stamp != 0 {
			query["created_at"] = bson.M{"$lt": stamp}
		}
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit("BoardLike", "like", handle)
}

func (db *DbHandle) ListUserLikes(uid string) (LikeModels, error) {
	items := make(LikeModels, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"uid": uid}
		return c.Find(query).Sort("-created_at").Limit(500).All(&items)
	}
	return items, mongo.Doit("BoardLike", "like", handle)
}

func (db *DbHandle) GetLike(id string) (*LikeModel, error) {
	pitem := &LikeModel{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, mongo.Doit("BoardLike", "like", handle)
}

func (db *DbHandle) NewLike(pitem *LikeModel) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(pitem)
	}
	return mongo.Doit("BoardLike", "like", handle)
}

func (db *DbHandle) DelLike(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit("BoardLike", "like", handle)
}

// -- Summary

func (db *DbHandle) MutiGetSummary(ids []string) (SummaryModels, error) {
	items := make(SummaryModels, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) GetSummary(id string) (*SummaryModel, error) {
	pitem := &SummaryModel{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(pitem)
	}
	return pitem, mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) IncrSummaryComm(cid, id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": 1}}
		if len(cid) == 0 {
			data["$inc"] = bson.M{"comms_count": 1, "comms_first_count": 1}
		}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) DecrSummaryComm(cid, id string, count int) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": -count}}
		if len(cid) == 0 {
			data["$inc"] = bson.M{"comms_count": -count, "comms_first_count": -1}
		}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) IncrSummaryLike(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) DecrSummaryLike(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) IncrSummaryRepost(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DbHandle) DecrSummaryRepost(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}
