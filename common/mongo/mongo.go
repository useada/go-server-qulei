package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"

	"a.com/go-server/common/tracing"
)

type Config struct {
	Host     string
	Auth     string
	Database []string
}

type Pool struct {
	Sessions map[string]*mgo.Session
}

func (p *Pool) Doit(c context.Context, db, collect string, h func(*mgo.Collection) error) error {
	span := tracing.StartDBSpan(c, "mongo", "do")
	defer span.Finish()

	session, ok := p.Sessions[db]
	if !ok {
		return errors.New("mongo session is nil")
	}
	conn := session.Copy()

	defer conn.Close()
	return h(conn.DB(db).C(collect))
}

func NewPool(conf Config) *Pool {
	pool := &Pool{
		Sessions: make(map[string]*mgo.Session),
	}

	for _, db := range conf.Database {
		addr := "mongodb://" + conf.Auth + conf.Host + "/" + db

		session, err := mgo.DialWithTimeout(addr, 10*time.Second)
		if err != nil {
			fmt.Println("Mongo数据库: ", db, "连接异常! ", err)
			panic(err)
		}

		session.SetMode(mgo.Monotonic, true)
		pool.Sessions[db] = session
	}

	fmt.Println("初始化Mongo连接池 FINISH")
	return pool
}
