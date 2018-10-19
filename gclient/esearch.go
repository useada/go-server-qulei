package gclient

import (
	"context"
	"time"

	proto "a.com/go-server/proto/pb/esearch"
)

type EsClient struct {
}

var Esearch *EsClient

func (e *EsClient) SearchByName(args *proto.NameRequest) (*proto.UserInfos, error) {
	conn, err := getConn(e.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout())
	defer cancel()
	return proto.NewEsearchClient(conn).SearchByName(ctx, args)
}

func (e *EsClient) SearchByNear(args *proto.NearRequest) (*proto.UserInfos, error) {
	conn, err := getConn(e.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout())
	defer cancel()
	return proto.NewEsearchClient(conn).SearchByNear(ctx, args)
}

func (e *EsClient) timeout() time.Duration {
	return 100 * time.Millisecond
}

func (e *EsClient) service() string {
	return "esearch"
}
