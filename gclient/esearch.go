package gclient

import (
	"context"
	"time"

	"a.com/go-server/proto/pb"
)

type EsClient struct {
}

var Esearch *EsClient

func (e *EsClient) UsersByName(args *pb.UsersByNameArgs) (*pb.UserInfos, error) {
	conn, err := GetConn(e.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout())
	defer cancel()
	return pb.NewEsearchClient(conn).UsersByName(ctx, args)
}

func (e *EsClient) UsersByNear(args *pb.UsersByNearArgs) (*pb.UserInfos, error) {
	conn, err := GetConn(e.service())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout())
	defer cancel()
	return pb.NewEsearchClient(conn).UsersByNear(ctx, args)
}

func (e *EsClient) timeout() time.Duration {
	return 100 * time.Millisecond
}

func (e *EsClient) service() string {
	return "esearch"
}
