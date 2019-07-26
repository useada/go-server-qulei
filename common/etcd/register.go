package etcd

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second
)

// NewRegister create a new consul register
func NewRegister(addr string) *Register {
	return &Register{
		Addrs:            []string{addr},
		KeepaliveTime:    keepaliveTime,
		KeepaliveTimeout: keepaliveTimeout,
	}
}

// Register  service register
type Register struct {
	Addrs            []string
	KeepaliveTime    time.Duration
	KeepaliveTimeout time.Duration
}

// Registe service
func (r *Register) Registe(name, ip string, port int) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:        r.Addrs,
		KeepaliveTime:    r.KeepaliveTime,
		KeepaliveTimeout: r.KeepaliveTimeout,
	})
	if err != nil {
		return err
	}

	lease := clientv3.NewLease(client)

	resp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		return err
	}

	if _, err := client.KV.Put(context.TODO(),
		fmt.Sprintf("%v/%v", name, ip), // key
		fmt.Sprintf("%v:%v", ip, port), // val
		clientv3.WithLease(resp.ID)); err != nil {
		return err
	}

	_, err = lease.KeepAlive(context, resp.ID)
	return err
}
