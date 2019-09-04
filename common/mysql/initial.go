package mysql

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ConfigNode struct {
	Host    string
	Auth    string
	MaxIdle int `toml:"max_idle"`
	MaxOpen int `toml:"max_open"`
	MaxLife int `toml:"max_life"`
}

type Config struct {
	Name   string
	Option string
	Master ConfigNode
	Slave  []ConfigNode
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

func connect(dbname string, dboption string, node ConfigNode) (*gorm.DB, error) {
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
	orm.DB().SetConnMaxLifetime(time.Duration(node.MaxLife) * time.Second)
	return orm, nil
}

var gInstance map[string]dbInstance

func Init(confs []Config) error {
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
			instance.Total++
		}
		gInstance[conf.Name] = instance
	}
	return nil
}
