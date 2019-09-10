package store

//go:generate mockgen -destination=../mock/store_mock.go -package=mock a.com/go-server/service/uauth/internal/store Store

import (
	"context"

	"a.com/go-server/service/uauth/internal/model"
)

type Store interface {
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	NewAccount(ctx context.Context, pitem *model.Account) error
	ModAccount(ctx context.Context, id string, data map[string]interface{}) error
}
