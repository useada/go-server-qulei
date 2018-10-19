package consul

import (
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

func NewConsulResolver(address string, service string) naming.Resolver {
	return &consulResolver{address: address, service: service}
}

type consulResolver struct {
	address string
	service string
}

func (r *consulResolver) Resolve(target string) (naming.Watcher, error) {
	client, err := api.NewClient(&api.Config{Address: r.address})
	if err != nil {
		return nil, err
	}
	return &consulWatcher{
		client:  client,
		service: r.service,
		addrs:   map[string]struct{}{},
	}, nil
}

type consulWatcher struct {
	client    *api.Client
	service   string
	addrs     map[string]struct{}
	lastIndex uint64
}

func (w *consulWatcher) Next() ([]*naming.Update, error) {
	for {
		services, metainfo, err := w.client.Health().Service(w.service, "", true, &api.QueryOptions{
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

func (w *consulWatcher) Close() {
	// nothing to do
}
