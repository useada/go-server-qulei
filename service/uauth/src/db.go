package main

import (
	"context"

	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"
)

type DbHandle struct {
}

var DB *DbHandle

func (db *DbHandle) GetAccount(ctx context.Context, id string) (*Account, error) {
	pitem := &Account{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("id=?", id).Find(pitem).Error
	}
	return pitem, mysql.Slave(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) NewAccount(ctx context.Context, pitem *Account) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(pitem).Error
	}
	return mysql.Master(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) ModAccount(ctx context.Context, id string, data map[string]interface{}) error {
	handle := func(orm *gorm.DB) error {
		return orm.Where("id=?", id).Updates(data).Error
	}
	return mysql.Master(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) DataBase() string {
	return "uauth"
}
