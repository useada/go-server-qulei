package gclient

import (
	"context"

	"a.com/go-server/proto/constant"
	"a.com/go-server/proto/pb"
)

func (cli *Client) Upload(c context.Context, args *pb.FileUploadArgs) (*pb.FileInfo, error) {
	conn, err := cli.Get(constant.UploadService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewUploaderClient(conn).Upload(ctx, args)
}

func (cli *Client) Query(c context.Context, args *pb.FileQueryArgs) (*pb.FileInfo, error) {
	conn, err := cli.Get(constant.UploadService)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, cli.timeout())
	defer cancel()
	return pb.NewUploaderClient(conn).Query(ctx, args)
}
