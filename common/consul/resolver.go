package consul

import (
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

func NewResolver(addr string) (naming.Resolver, error) {
	return &resolver{address: addr}, nil
}

type resolver struct {
	address string
}

func (r *resolver) Resolve(target string) (naming.Watcher, error) {
	client, err := api.NewClient(&api.Config{Address: r.address})
	if err != nil {
		return nil, err
	}

	return &watcher{
		client: client,
		target: target,
		addrs:  map[string]struct{}{},
	}, nil
}

type watcher struct {
	client    *api.Client
	target    string
	addrs     map[string]struct{}
	lastIndex uint64
}

func (w *watcher) Next() ([]*naming.Update, error) {
	for {
		services, metainfo, err := w.client.Health().Service(w.target, "", true, &api.QueryOptions{
			AllowStale: true,
			WaitIndex:  w.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
		})
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}

		w.lastIndex = metainfo.LastIndex
		addrs := map[string]struct{}{}
		for _, service := range services {
			addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
		}

		var updates []*naming.Update
		for addr := range w.addrs {
			if _, ok := addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr})
			}
		}
		for addr := range addrs {
			if _, ok := w.addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr})
			}
		}

		if len(updates) != 0 {
			w.addrs = addrs
			return updates, nil
		}
	}
}

func (w *watcher) Close() {
}
