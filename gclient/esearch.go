package gclient

import (
	"context"
	"time"

	"a.com/go-server/proto/pb"
)

type EsClient struct {
}

var Esearch *EsClient

func (e *EsClient) SearchByName(args *pb.SearchNameArgs) (*pb.SearchUserInfos, error) {
	conn, err := getConn(e.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout())
	defer cancel()
	return pb.NewEsearchClient(conn).SearchByName(ctx, args)
}

func (e *EsClient) SearchByNear(args *pb.SearchNearArgs) (*pb.SearchUserInfos, error) {
	conn, err := getConn(e.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout())
	defer cancel()
	return pb.NewEsearchClient(conn).SearchByNear(ctx, args)
}

func (e *EsClient) timeout() time.Duration {
	return 100 * time.Millisecond
}

func (e *EsClient) service() string {
	return "esearch"
}
