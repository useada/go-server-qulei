package gclient

import (
	"context"
	"time"

	pb "a.com/go-server/proto/pb"
)

type UploaderClient struct {
}

var Uploader *UploaderClient

func (u *UploaderClient) Upload(args *pb.FileUploadArgs) (*pb.FileInfo, error) {
	conn, err := getConn(u.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout())
	defer cancel()
	return pb.NewUploaderClient(conn).Upload(ctx, args)
}

func (u *UploaderClient) Query(args *pb.FileQueryArgs) (*pb.FileInfo, error) {
	conn, err := getConn(u.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout())
	defer cancel()
	return pb.NewUploaderClient(conn).Query(ctx, args)
}

func (u *UploaderClient) timeout() time.Duration {
	return 500 * time.Millisecond
}

func (u *UploaderClient) service() string {
	return "uploader"
}
