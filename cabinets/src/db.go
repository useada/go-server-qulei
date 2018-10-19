package main

import (
	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"
)

type DBHandle struct {
}

var DB *DBHandle

func (d *DBHandle) GetFileInfo(fid string) (FileInfoDB, error) {
	row := FileInfoDB{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("fid=?", fid).Find(&row).Error
	}
	return row, mysql.Doit(d.DataBase(), handle)
}

func (d *DBHandle) AddFileInfo(info *FileInfoDB) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(info).Error
	}
	return mysql.Doit(d.DataBase(), handle)
}

func (d *DBHandle) DataBase() string {
	return "cabinets"
}

type FileInfoDB struct {
	Id     string `json:"id"`    // 文件id <primary key>
	Ex     string `json:"ex"`    // 文件扩展名
	Typef  int    `json:"typef"` // 文件类型 头像/图片/普通文件
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int64  `json:"size"`
}

func (f *FileInfoDB) TableName() string {
	return "file_info"
}
