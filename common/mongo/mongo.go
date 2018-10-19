package mongo

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"

	configor "a.com/server/mywork/common/configor"
)

func Doit(db, collect string, h func(*mgo.Collection) error) error {
	sess, ok := gMgo[db]
	if !ok {
		return errors.New("mongo session is nil")
	}
	conn := sess.Copy()
	defer conn.Close()

	return h(conn.DB(db).C(collect))
}

var gMgo map[string]*mgo.Session

func InitMongo(conf configor.MongoConfigor) error {
	gMgo = make(map[string]*mgo.Session)
	for _, db := range conf.Database {
		fmt.Println("初始化Mongo数据库:", db, " 连接池")

		addr := "mongodb://" + conf.Auth + conf.Host + db
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
