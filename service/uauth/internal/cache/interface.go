package cache

//go:generate mockgen -destination=../mock/cache_mock.go -package=mock a.com/go-server/service/uauth/internal/cache Cache

import (
	"context"

	"a.com/go-server/service/uauth/internal/model"
)

type Cache interface {
	GetRecord(ctx context.Context, uid, device string) (model.Record, error)
	NewRecord(ctx context.Context, r model.Record) error
}
