package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

// NewRegister create a new consul register
func NewRegister(addr string) *Register {
	return &Register{
		address:  addr,
		Timeout:  time.Duration(1) * time.Minute,
		Interval: time.Duration(5) * time.Second,
	}
}

// Register  service register
type Register struct {
	address  string
	Timeout  time.Duration
	Interval time.Duration
}

// Registe service
func (r *Register) Registe(name, ip string, port int) error {
	client, err := api.NewClient(&api.Config{Address: r.address})
	if err != nil {
		return err
	}

	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:   fmt.Sprintf("%v-%v-%v", name, ip, port), // 服务节点的名称
		Name: fmt.Sprintf("%v", name),                 // 服务名称
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%v:%v/%v", ip, port, name), // grpc支持,健康检查的地址
			Interval:                       r.Interval.String(),                     // 健康检查间隔
			DeregisterCriticalServiceAfter: r.Timeout.String(),                      // 注销时间，相当于过期时间
		},
		Address: ip,   // 服务IP
		Port:    port, // 服务端口
	})
}
