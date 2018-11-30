package gclient

import (
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"a.com/go-server/common/consul"
)

var GrpcConns map[string]*grpc.ClientConn

func Init(consul string, services []string) error {
	if len(consul) == 0 || len(services) == 0 {
		return errors.New("consul or services empty")
	}
	GrpcConns = make(map[string]*grpc.ClientConn, len(services))
	for _, service := range services {
		conn, err := newConn(consul, service)
		if err != nil {
			return err
		}
		GrpcConns[service] = conn
	}
	return nil
}

func GetConn(service string) (*grpc.ClientConn, error) {
	conn, ok := GrpcConns[service]
	if !ok {
		return nil, errors.New(service + " conn not exist")
	}
	return conn, nil
}

func newConn(host, service string) (*grpc.ClientConn, error) {
	r := consul.NewConsulResolver(host, service)
	return grpc.Dial("", grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			PermitWithoutStream: true,
			Time:                500 * time.Millisecond,
		}),
		//grpc.WithBlock(),
		grpc.WithBalancer(grpc.RoundRobin(r)))
}
