package main

import (
	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"
)

type DbHandle struct {
}

var DB *DbHandle

func (db *DbHandle) GetFileInfo(fid string) (FileInfoDB, error) {
	row := FileInfoDB{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("fid=?", fid).Find(&row).Error
	}
	return row, mysql.Doit(db.DataBase(), handle)
}

func (db *DbHandle) AddFileInfo(info *FileInfoDB) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(info).Error
	}
	return mysql.Doit(db.DataBase(), handle)
}

func (db *DbHandle) DataBase() string {
	return "uploader"
}
