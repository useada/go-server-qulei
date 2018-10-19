package mysql

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	configor "a.com/go-server/common/configor"
)

func Doit(db string, h func(*gorm.DB) error) error {
	orm, ok := gGorm[db]
	if !ok {
		return errors.New("mysql orm is nil")
	}
	return h(orm)
}

var gGorm map[string]*gorm.DB

func InitMysql(conf configor.MysqlConfigor) error {
	gGorm = make(map[string]*gorm.DB)
	for _, db := range conf.Database {
		dst := fmt.Sprintf("%s@tcp(%s)/%s", conf.Auth, conf.Host, db)
		fmt.Println("连接Mysql:", dst)

		orm, err := gorm.Open("mysql", dst)
		if err != nil {
			fmt.Println("Mysql数据库: ", db, "连接异常! ", err)
			return err
		}

		orm.LogMode(true)
		orm.DB().SetMaxIdleConns(conf.MaxIdle)
		orm.DB().SetMaxOpenConns(conf.MaxOpen)
		gGorm[db] = orm
	}
	return nil
}
