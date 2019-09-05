package cloud

import (
	"a.com/go-server/proto/pb"
)

type Cloud interface {
	Upload(id, ex string, xtype pb.TYPE, data []byte) error
	Delete(remoteDir string) error
	GetBucket() string
	RemoteDir(id, ex string, xtype pb.TYPE) string
}
