package store

import (
	"context"

	"github.com/jinzhu/gorm"

	"a.com/go-server/service/uauth/internal/model"
)

func (d *db) GetAccount(ctx context.Context, id string) (*model.Account, error) {
	pitem := &model.Account{}
	handle := func(orm *gorm.DB) error {
		return orm.Where("id=?", id).Find(pitem).Error
	}
	return pitem, d.Pool.Slave(d.table()).Doit(ctx, handle)
}

func (d *db) NewAccount(ctx context.Context, pitem *model.Account) error {
	handle := func(orm *gorm.DB) error {
		return orm.Create(pitem).Error
	}
	return d.Pool.Master(d.table()).Doit(ctx, handle)
}

func (d *db) ModAccount(ctx context.Context, id string, data map[string]interface{}) error {
	handle := func(orm *gorm.DB) error {
		return orm.Where("id=?", id).Updates(data).Error
	}
	return d.Pool.Master(d.table()).Doit(ctx, handle)
}

func (d *db) table() string {
	return "uauth"
}
