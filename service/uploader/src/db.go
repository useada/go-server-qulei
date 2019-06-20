package main

import (
	"context"

	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"
)

type DbHandle struct {
}

var DB *DbHandle

func (db *DbHandle) GetFileInfo(ctx context.Context, fid string) (*FileInfoModel, error) {
	pitem := &FileInfoModel{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("fid=?", fid).Find(pitem).Error
	}
	return pitem, mysql.Slave(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) AddFileInfo(ctx context.Context, item *FileInfoModel) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(item).Error
	}
	return mysql.Master(db.DataBase()).Doit(ctx, handle)
}

func (db *DbHandle) DataBase() string {
	return "uploader"
}
