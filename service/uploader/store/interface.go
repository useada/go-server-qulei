package store

import (
	"context"

	"a.com/go-server/service/uploader/model"
)

type Store interface {
	GetFileInfo(ctx context.Context, fid string) (*model.FileInfo, error)
	AddFileInfo(ctx context.Context, item *model.FileInfo) error
}
