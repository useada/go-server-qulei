package mysql

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"a.com/go-server/common/tracing"
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
	Slaves []ConfigNode
}

type Client struct {
	*gorm.DB
}

func (c *Client) Doit(ctx context.Context, h func(*gorm.DB) error) error {
	if c == nil {
		return errors.New("mysql instance is nil")
	}

	span := tracing.StartDBSpan(ctx, "mysql", "do")
	defer span.Finish()

	return h(c.DB)
}

type dbInstance struct {
	Name   string
	master *Client
	slaves []*Client
	next   uint64
	total  uint64
}

type Pool struct {
	Instances map[string]dbInstance
}

func (p *Pool) Master(dbname string) *Client {
	instance, ok := p.Instances[dbname]
	if !ok {
		return nil
	}
	return instance.master
}

func (p *Pool) Slave(dbname string) *Client {
	instance, ok := p.Instances[dbname]
	if !ok {
		return nil
	}

	if instance.total == 0 {
		return instance.master
	}

	next := atomic.AddUint64(&instance.next, 1)
	return instance.slaves[next%instance.total]
}

func NewPool(confs []Config) *Pool {
	pool := &Pool{
		Instances: make(map[string]dbInstance),
	}
	for _, conf := range confs {
		orm, err := connect(conf.Name, conf.Option, conf.Master)
		if err != nil {
			panic(err)
		}
		instance := dbInstance{
			master: &Client{DB: orm},
		}

		for _, slave := range conf.Slaves {
			orm, err := connect(conf.Name, conf.Option, slave)
			if err != nil {
				continue
			}
			instance.slaves = append(instance.slaves, &Client{DB: orm})
			instance.total++
		}
		pool.Instances[conf.Name] = instance
	}

	fmt.Println("初始化Mysql连接池 FINISH")
	return pool
}

func connect(dbname string, dboption string, node ConfigNode) (*gorm.DB, error) {
	dst := fmt.Sprintf("%s@tcp(%s)/%s", node.Auth, node.Host, dbname)
	if len(dboption) > 0 {
		dst = dst + "?" + dboption
	}

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
