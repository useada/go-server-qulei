package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"

	"a.com/go-server/common/configor"
	"a.com/go-server/common/tracing"
)

func Doit(c context.Context, db, collect string, h func(*mgo.Collection) error) error {
	span := tracing.StartDBSpan(c, "mongo", "do")
	defer span.Finish()

	sess, ok := gMgo[db]
	if !ok {
		return errors.New("mongo session is nil")
	}
	conn := sess.Copy()
	defer conn.Close()
	return h(conn.DB(db).C(collect))
}

var gMgo map[string]*mgo.Session

func Init(conf configor.MongoConfigor) error {
	gMgo = make(map[string]*mgo.Session)
	for _, db := range conf.Database {
		addr := "mongodb://" + conf.Auth + conf.Host + "/" + db
		fmt.Println("初始化Mongo数据库:", addr, " 连接池")

		sess, err := mgo.DialWithTimeout(addr, 10*time.Second)
		if err != nil {
			fmt.Println("Mongo数据库: ", db, "连接异常! ", err)
			return err
		}
		sess.SetMode(mgo.Monotonic, true)
		gMgo[db] = sess
	}
	return nil
}
