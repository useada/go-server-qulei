package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterUploaderServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) Upload(ctx context.Context, in *pb.FileUploadArgs) (*pb.FileInfo, error) {
	if in.Type != pb.TYPE_FILE { // 如果是图片 解析图片 获取宽/高
		img, _, _ := image.DecodeConfig(bytes.NewReader(in.Data))
		in.Width, in.Height = int32(img.Width), int32(img.Height)
	}

	if err := S3.Upload(in.Id, in.Ex, in.Type, in.Data); err != nil {
		return nil, err
	}

	pitem := &FileInfoModel{}
	if err := DB.AddFileInfo(ctx, pitem.DestructPb(in)); err != nil {
		fmt.Println("add file info error:", err)
	}
	return pitem.ConstructPb(), nil
}

func (s *SvrHandler) Query(ctx context.Context, in *pb.FileQueryArgs) (*pb.FileInfo, error) {
	pitem, err := DB.GetFileInfo(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return pitem.ConstructPb(), nil
}
