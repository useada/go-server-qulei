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

func (s *SvrHandler) Upload(ctx context.Context,
	in *pb.UploadRequest) (*pb.FileInfo, error) {

	res := pb.FileInfo{
		Typef: in.Typef,
		Ex:    in.Ex,
		Id:    in.Id,
		Size:  int64(len(in.Data)),
	}
	// 如果是图片 解析图片 获取宽/高
	if in.Typef != pb.FileType_FILE {
		img, _, err := image.DecodeConfig(bytes.NewReader(in.Data))
		if err != nil {
			return nil, err
		}
		res.Width = int32(img.Width)
		res.Height = int32(img.Height)
	}

	// 文件上传到S3
	if err := S3.Upload(in.Id, in.Ex, in.Typef, in.Data); err != nil {
		return nil, err
	}

	// 上传成功 生成一条记录
	if err := DB.AddFileInfo(&FileInfoDB{
		Id:     res.Id,
		Ex:     res.Ex,
		Typef:  int(res.Typef),
		Width:  int(res.Width),
		Height: int(res.Height),
		Size:   res.Size,
	}); err != nil {
		fmt.Println("add file info error:", err)
	}
	return &res, nil
}

func (s *SvrHandler) GetFileInfo(ctx context.Context,
	in *pb.FileRequest) (*pb.FileInfo, error) {

	row, err := DB.GetFileInfo(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.FileInfo{
		Id:     row.Id,
		Ex:     row.Ex,
		Typef:  pb.FileType(row.Typef),
		Width:  int32(row.Width),
		Height: int32(row.Height),
		Size:   row.Size,
	}, nil
}
