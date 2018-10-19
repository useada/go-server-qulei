package gclient

import (
	"context"
	"time"

	proto "a.com/go-server/proto/pb/cabinets"
)

type CabClient struct {
}

var Cabinets *CabClient

func (c *CabClient) Upload(args *proto.UploadRequest) (*proto.FileInfo, error) {
	conn, err := getConn(c.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout())
	defer cancel()
	return proto.NewCabinetsClient(conn).Upload(ctx, args)
}

func (c *CabClient) GetFileInfo(args *proto.FileRequest) (*proto.FileInfo, error) {
	conn, err := getConn(c.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout())
	defer cancel()
	return proto.NewCabinetsClient(conn).GetFileInfo(ctx, args)
}

func (c *CabClient) timeout() time.Duration {
	return 500 * time.Millisecond
}

func (c *CabClient) service() string {
	return "cabinets"
}
