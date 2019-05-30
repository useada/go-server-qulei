package mysql

import (
	"fmt"
	"sync/atomic"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MysqlNode struct {
	Host    string
	Auth    string
	MaxIdle int `toml:"max_idle"`
	MaxOpen int `toml:"max_open"`
}

type MysqlConfigor struct {
	Name   string
	Option string
	Master MysqlNode
	Slave  []MysqlNode
}

type dbInstance struct {
	Name   string
	Master *Client
	Slave  []*Client
	Next   uint64
	Total  uint64
}

func (i *dbInstance) getMaster() *Client {
	return i.Master
}

func (i *dbInstance) getSlave() *Client {
	if i.Total == 0 {
		return i.Master
	}
	next := atomic.AddUint64(&i.Next, 1)
	return i.Slave[next%i.Total]
}

func connect(dbname string, dboption string, node MysqlNode) (*gorm.DB, error) {
	dst := fmt.Sprintf("%s@tcp(%s)/%s", node.Auth, node.Host, dbname)
	if len(dboption) > 0 {
		dst = dst + "?" + dboption
	}
	fmt.Println("连接Mysql:", dst)

	orm, err := gorm.Open("mysql", dst)
	if err != nil {
		fmt.Println("Mysql数据库: ", dbname, "连接异常! ", err)
		return nil, err
	}
	orm.LogMode(true)
	orm.DB().SetMaxIdleConns(node.MaxIdle)
	orm.DB().SetMaxOpenConns(node.MaxOpen)
	return orm, nil
}

var gInstance map[string]dbInstance

func Init(confs []MysqlConfigor) error {
	gInstance = make(map[string]dbInstance)
	for _, conf := range confs {
		instance := dbInstance{}

		// Master
		orm, err := connect(conf.Name, conf.Option, conf.Master)
		if err != nil {
			return err
		}
		instance.Master = &Client{DB: orm}

		// Slave
		for _, slave := range conf.Slave {
			orm, err := connect(conf.Name, conf.Option, slave)
			if err != nil {
				continue
			}
			instance.Slave = append(instance.Slave, &Client{DB: orm})
			instance.Total += 1
		}
		gInstance[conf.Name] = instance
	}
	return nil
}
