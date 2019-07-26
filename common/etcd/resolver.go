package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"
	grpc "google.golang.org/grpc/naming"
)

func NewResolver(addr string) (grpc.Resolver, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:        []string{addr},
		KeepaliveTime:    keepaliveTime,
		KeepaliveTimeout: keepaliveTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &naming.GRPCResolver{Client: client}, nil
}
