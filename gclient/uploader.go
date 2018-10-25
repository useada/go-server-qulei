package gclient

import (
	"context"
	"time"

	pb "a.com/go-server/proto/pb"
)

type UploaderClient struct {
}

var Uploader *UploaderClient

func (u *UploaderClient) Upload(args *pb.UploadRequest) (*pb.FileInfo, error) {
	conn, err := getConn(u.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout())
	defer cancel()
	return pb.NewUploaderClient(conn).Upload(ctx, args)
}

func (u *UploaderClient) GetFileInfo(args *pb.FileRequest) (*pb.FileInfo, error) {
	conn, err := getConn(u.service())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout())
	defer cancel()
	return pb.NewUploaderClient(conn).GetFileInfo(ctx, args)
}

func (u *UploaderClient) timeout() time.Duration {
	return 500 * time.Millisecond
}

func (u *UploaderClient) service() string {
	return "uploader"
}
