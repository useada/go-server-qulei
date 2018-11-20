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

func (db *DBHandle) ListFirstComms(oid, direct string,
	stamp int64, limit int) (CommRecordList, error) {
	items := make(CommRecordList, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "status": 0, "level": 1}
		if len(direct) != 0 {
			query["created_at"] = bson.M{"$" + direct: stamp} //$gt, $lt
		}

		sorter := "-created_at"
		if direct == "gt" || direct == "gte" {
			sorter = "created_at"
		}

		err := c.Find(query).Sort(sorter).Limit(limit).All(&items)
		if len(items) > 0 && sorter == "-created_at" {
			sort.Sort(items)
		}
		return err
	}
	return items, mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) ListChildComms(oid, cid, direct string,
	stamp int64, limit int) (CommRecordList, error) {
	items := make(CommRecordList, 0)
	handle := func(c *mgo.Collection) error {
		query := bson.M{"oid": oid, "status": 0, "level": 2, "cid": cid}
		if len(direct) != 0 {
			query["created_at"] = bson.M{"$" + direct: stamp} //$gt, $lt
		}

		sorter := "-created_at"
		if direct == "gt" || direct == "gte" {
			sorter = "created_at"
		}

		err := c.Find(query).Sort(sorter).Limit(limit).All(&items)
		if len(items) > 0 && sorter == "-created_at" {
			sort.Sort(items)
		}
		return err
	}
	return items, mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) MGetComments(ids []string) (CommRecordList, error) {
	items := make(CommRecordList, 0)
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&items)
	}
	return items, mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) GetComment(id string) (CommRecord, error) {
	item := CommRecord{}
	handle := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": id}).One(&item)
	}
	return item, mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) NewComment(item CommRecord) error {
	handle := func(c *mgo.Collection) error {
		return c.Insert(item)
	}
	return mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) DelComment(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": id},
			bson.M{"$set": bson.M{"status": DELETE_STATUS}})
	}
	return mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) IncrCommReply(id string) error {
	return nil
}

func (db *DBHandle) DecrCommReply(id string) error {
	return nil
}

func (db *DBHandle) IncrCommLike(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": id},
			bson.M{"$inc": bson.M{"likes_count": 1}})
	}
	return mongo.Doit("Comment", "comment", handle)
}

func (db *DBHandle) DecrCommLike(id string) error {
	handle := func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": id},
			bson.M{"$inc": bson.M{"likes_count": -1}})
	}
	return mongo.Doit("Comment", "comment", handle)
}
