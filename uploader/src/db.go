package main

import (
	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"
)

type DbHandle struct {
}

var DB *DbHandle

func (db *DbHandle) GetFileInfo(fid string) (*FileInfoModel, error) {
	pitem := &FileInfoModel{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("fid=?", fid).Find(pitem).Error
	}
	return pitem, mysql.Doit(db.DataBase(), handle)
}

func (db *DbHandle) AddFileInfo(item *FileInfoModel) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(item).Error
	}
	return mysql.Doit(db.DataBase(), handle)
}

func (db *DbHandle) DataBase() string {
	return "uploader"
}
