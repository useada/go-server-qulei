package main

import (
	"context"

	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"
	"a.com/go-server/common/utime"
)

type DbHandle struct {
}

var DB *DbHandle

func (db *DbHandle) GetAuthInfo(ctx context.Context, id string) (*AuthInfoModel, error) {
	pitem := &AuthInfoModel{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("id=?", id).Find(pitem).Error
	}
	return pitem, mysql.Slave(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) NewAuthInfo(ctx context.Context, pitem *AuthInfoModel) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(pitem).Error
	}
	return mysql.Master(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) ModAuthInfo(ctx context.Context, id, field string, value interface{}) error {
	handle := func(orm *gorm.DB) error {
		return orm.Where("id=?", id).Updates(map[string]interface{}{
			field:        value,
			"updated_at": utime.Millisec(),
		}).Error
	}
	return mysql.Master(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) DataBase() string {
	return "uauth"
}
