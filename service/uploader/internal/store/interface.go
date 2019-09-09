package store

//go:generate mockgen -destination=../mock/store_mock.go -package=mock a.com/go-server/service/upload/internal/store Store

import (
	"context"

	"a.com/go-server/service/uploader/internal/model"
)

type Store interface {
	GetFileInfo(ctx context.Context, fid string) (*model.FileInfo, error)
	AddFileInfo(ctx context.Context, item *model.FileInfo) error
}
