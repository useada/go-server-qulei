package gclient

import (
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"a.com/go-server/common/consul"
	"a.com/go-server/common/tracing"
)

type Client struct {
	Conns         map[string]*grpc.ClientConn
	RequstTimeout int64
}

func (cli *Client) Get(service string) (*grpc.ClientConn, error) {
	conn, ok := cli.Conns[service]
	if !ok {
		return nil, errors.New(service + " conn not exist")
	}
	return conn, nil
}

func (cli *Client) timeout() time.Duration {
	return time.Duration(cli.RequstTimeout) * time.Millisecond
}

type Config struct {
	Discover      string
	Services      []string
	RequstTimeout int64 `toml:"request_timeout"`
}

func NewGrpcClient(conf Config) *Client {
	if len(conf.Discover) == 0 || len(conf.Services) == 0 {
		panic("discovery or services empty")
	}
	if conf.RequstTimeout == 0 {
		conf.RequstTimeout = 500
	}

	client := &Client{
		Conns:         make(map[string]*grpc.ClientConn, len(conf.Services)),
		RequstTimeout: conf.RequstTimeout,
	}

	for _, service := range conf.Services {
		conn, err := connect(conf.Discover, service)
		if err != nil {
			panic(err)
		}
		client.Conns[service] = conn
	}
	return client
}

func connect(discover, service string) (*grpc.ClientConn, error) {
	r, err := consul.NewResolver(discover)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(service, grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			PermitWithoutStream: true,
			Time:                500 * time.Millisecond,
		}),
		grpc.WithUnaryInterceptor(tracing.GrpcClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithBalancer(grpc.RoundRobin(r)))
}
