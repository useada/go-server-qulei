```
import (
    "a.com/go-server/common/mongo"
)

// 全局初始化调用一次
func init() {
    if err := mongo.Init(conf); err != nil {
        panic(err)
    }
}

// 使用示例
func DemoListMongoData(sid string, offset, limit int) (CommentList, error) {
    items := make(CommentList, 0)
    h := func(c *mgo.Collection) error {
        return c.Find(bson.M{"source_id": sid, "status": 0}).
            Select(bson.M{"replys": bson.M{"$slice": 2000}}).
            Sort("-created_at").Skip(offset).Limit(limit).All(&items)
    }
    return items, mongo.Doit("Comment", "comment", handle)
}

func DemoCountMongoData(sid string) (count int, err error) {
    h := func(c *mgo.Collection) error {
        count, err = c.Find(bson.M{"source_id": sid, "status": 0}).Count()
        return err
    }
    return count, mongo.Doit("Comment", "comment", h) // 库: Comment; 集合: comment
}
```
