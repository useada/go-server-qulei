package gclient

import (
	"context"

	"a.com/go-server/proto/constant"
	"a.com/go-server/proto/pb"
)

func (cli *Client) UsersByName(c context.Context, args *pb.UsersByNameArgs) (*pb.UserInfos, error) {
	conn, err := cli.Get(constant.SearchService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewEsearchClient(conn).UsersByName(ctx, args)
}

func (cli *Client) UsersByNear(c context.Context, args *pb.UsersByNearArgs) (*pb.UserInfos, error) {
	conn, err := cli.Get(constant.SearchService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewEsearchClient(conn).UsersByNear(ctx, args)
}
