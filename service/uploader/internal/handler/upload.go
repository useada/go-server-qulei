package handler

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"
	"a.com/go-server/service/uploader/internal/cloud"
	"a.com/go-server/service/uploader/internal/model"
	"a.com/go-server/service/uploader/internal/store"
)

func RegisterHandler(svr *grpc.Server, store store.Store, cloud cloud.Cloud, log *zap.SugaredLogger) {
	pb.RegisterUploaderServer(svr, &SvrHandler{
		Store: store,
		Cloud: cloud,
		Log:   log,
	})
}

type SvrHandler struct {
	Store store.Store
	Cloud cloud.Cloud
	Log   *zap.SugaredLogger
}

func (s *SvrHandler) Upload(ctx context.Context, in *pb.FileUploadArgs) (*pb.FileInfo, error) {
	if in.Type != pb.TYPE_FILE { // 如果是图片 解析图片 获取宽/高
		img, _, _ := image.DecodeConfig(bytes.NewReader(in.Data))
		in.Width, in.Height = int32(img.Width), int32(img.Height)
	}

	if err := s.Cloud.Upload(in.Id, in.Ex, in.Type, in.Data); err != nil {
		return nil, err
	}

	pitem := &model.FileInfo{}
	if err := s.Store.AddFileInfo(ctx, pitem.DestructPb(in)); err != nil {
		s.Log.Error("add file info err:%v", err)
	}
	return pitem.ConstructPb(), nil
}

func (s *SvrHandler) Query(ctx context.Context, in *pb.FileQueryArgs) (*pb.FileInfo, error) {
	pitem, err := s.Store.GetFileInfo(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return pitem.ConstructPb(), nil
}
