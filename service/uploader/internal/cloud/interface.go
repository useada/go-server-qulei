package cloud

//go:generate mockgen -destination=../mock/store_mock.go -package=mock a.com/go-server/service/upload/internal/cloud Cloud

import (
	"a.com/go-server/proto/pb"
)

type Cloud interface {
	Upload(id, ex string, xtype pb.TYPE, data []byte) error
	Delete(remoteDir string) error
	GetBucket() string
	RemoteDir(id, ex string, xtype pb.TYPE) string
}
