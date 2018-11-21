package main

import (
	"sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"a.com/go-server/common/mongo"
)

type DBHandle struct {
}

var DB = &DBHandle{}

// -- Comment

func (db *DBHandle) ListFirstComms(oid, direct string,
	stamp int64, limit int) (CommRecordList, error) {
	items := make(CommRecordList, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "level": COMMENT_LEVEL_FIRST}
		if len(direct) != 0 {
			query["created_at"] = bson.M{"$" + direct: stamp} //$gt, $lt
		}

		if direct == "gt" || direct == "gte" {
			return c.Find(query).Sort("created_at").Limit(limit).All(&items)
		}

		defer sort.Sort(items) // created_at升序返回
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) ListChildComms(oid, cid, direct string,
	stamp int64, limit int) (CommRecordList, error) {
	items := make(CommRecordList, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "level": COMMENT_LEVEL_CHILD, "cid": cid}
		if len(direct) != 0 {
			query["created_at"] = bson.M{"$" + direct: stamp} //$gt, $lt
		}

		if direct == "gt" || direct == "gte" {
			return c.Find(query).Sort("created_at").Limit(limit).All(&items)
		}

		defer sort.Sort(items) // created_at升序返回
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) MGetComments(ids []string) (CommRecordList, error) {
	items := make(CommRecordList, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) GetComment(id string) (CommRecord, error) {
	item := CommRecord{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(&item)
	}
	return item, mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) NewComment(item CommRecord) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(item)
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) DelComment(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) IncrCommReply(cid string, item CommRecord) error {
	handle := func(c *mgo.Collection) error {
		replys := bson.M{"$each": CommRecordList{item},
			"$sort": "-created_at", "$slice": 3}
		return c.Update(bson.M{"_id": cid},
			bson.M{"$push": bson.M{"replys": replys},
				"$inc": bson.M{"reply_count": 1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) DecrCommReply(cid, rid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$pull": bson.M{"replys": bson.M{"_id": rid}},
				"$inc": bson.M{"reply_count": -1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) IncrCommLike(cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": 1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) DecrCommLike(cid string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": cid},
			bson.M{"$inc": bson.M{"likes_count": -1}})
	}
	return mongo.Doit("BoardComment", "comment", handle)
}

func (db *DBHandle) GetsUserCommLikes(uid string) (CommLikeRecordList, error) {
	items := make(CommLikeRecordList, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).Sort("-created_at").Limit(500).All(&items)
	}
	return items, mongo.Doit("BoardComment", "like", handle)
}

func (db *DBHandle) NewCommLike(item CommLikeRecord) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(item)
	}
	return mongo.Doit("BoardComment", "like", handle)
}

func (db *DBHandle) DelCommLike(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit("BoardComment", "like", handle)
}

// -- Like

func (db *DBHandle) ListLikes(oid string,
	stamp int64, limit int) (LikeRecordList, error) {
	items := make(LikeRecordList, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid}
		if stamp != 0 {
			query["created_at"] = bson.M{"$lt": stamp}
		}
		return c.Find(query).Sort("-created_at").Limit(limit).All(&items)
	}
	return items, mongo.Doit("BoardLike", "like", handle)
}

func (db *DBHandle) GetsUserLikes(uid string) (LikeRecordList, error) {
	items := make(LikeRecordList, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"uid": uid}
		return c.Find(query).Sort("-created_at").Limit(500).All(&items)
	}
	return items, mongo.Doit("BoardLike", "like", handle)
}

func (db *DBHandle) GetLike(id string) (LikeRecord, error) {
	item := LikeRecord{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(&item)
	}
	return item, mongo.Doit("BoardLike", "like", handle)
}

func (db *DBHandle) NewLike(item LikeRecord) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(item)
	}
	return mongo.Doit("BoardLike", "like", handle)
}

func (db *DBHandle) DelLike(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Remove(bson.M{"_id": id})
	}
	return mongo.Doit("BoardLike", "like", handle)
}

// -- Summary

func (db *DBHandle) MGetsSummary(ids []string) (SumRecordList, error) {
	items := make(SumRecordList, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) GetSummary(id string) (SumRecord, error) {
	item := SumRecord{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(&item)
	}
	return item, mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) IncrSummaryComm(id string, level int) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": 1}}
		if level == COMMENT_LEVEL_FIRST {
			data["$inc"] = bson.M{"comms_count": 1, "comms_first_count": 1}
		}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) DecrSummaryComm(id string, level, n int) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"comms_count": -n}}
		if level == COMMENT_LEVEL_FIRST {
			data["$inc"] = bson.M{"comms_count": -n, "comms_first_count": -1}
		}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) IncrSummaryLike(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) DecrSummaryLike(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"likes_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) IncrSummaryRepost(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": 1}}
		_, err := c.Upsert(bson.M{"_id": id}, data)
		return err
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}

func (db *DBHandle) DecrSummaryRepost(id string) error {
	handle := func(c *mgo.Collection) error {
		data := bson.M{"$inc": bson.M{"repost_count": -1}}
		return c.Update(bson.M{"_id": id}, data)
	}
	return mongo.Doit("BoardSummary", "summary", handle)
}
