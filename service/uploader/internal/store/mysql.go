package store

import (
	"context"

	"github.com/jinzhu/gorm"

	"a.com/go-server/common/mysql"

	"a.com/go-server/service/uploader/internal/model"
)

type db struct {
	Pool *mysql.Pool
}

func NewMysqlRepo(pool *mysql.Pool) Store {
	return &db{Pool: pool}
}

func (d *db) GetFileInfo(ctx context.Context, fid string) (*model.FileInfo, error) {
	pitem := &model.FileInfo{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("fid=?", fid).Find(pitem).Error
	}
	return pitem, d.Pool.Slave(d.table()).Doit(ctx, handle)
}

func (d *db) AddFileInfo(ctx context.Context, item *model.FileInfo) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(item).Error
	}
	return d.Pool.Master(d.table()).Doit(ctx, handle)
}

func (d *db) table() string {
	return "uploader"
}
