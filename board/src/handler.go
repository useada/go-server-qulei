package main

import (
	"context"

	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterBoardServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) ListFirstComments(ctx context.Context,
	in *pb.CommListArgs) (*pb.CommentInfos, error) {
	return nil, nil
}

func (s *SvrHandler) ListChildComments(ctx context.Context,
	in *pb.CommListArgs) (*pb.CommentInfos, error) {
	return nil, nil
}

func (s *SvrHandler) GetComment(ctx context.Context,
	in *pb.CommGetArgs) (*pb.CommentInfo, error) {
	return nil, nil
}

func (s *SvrHandler) NewComment(ctx context.Context,
	in *pb.CommNewArgs) (*pb.CommentInfo, error) {
	return nil, nil
}

func (s *SvrHandler) DelComment(ctx context.Context,
	in *pb.CommDelArgs) (*pb.ReplyBaseInfo, error) {
	return nil, nil
}

func (s *SvrHandler) LikeComment(ctx context.Context,
	in *pb.CommLikeArgs) (*pb.ReplyBaseInfo, error) {
	return nil, nil
}

func (s *SvrHandler) UnLikeComment(ctx context.Context,
	in *pb.CommUnLikeArgs) (*pb.ReplyBaseInfo, error) {
	return nil, nil
}

func (s *SvrHandler) ListLikes(ctx context.Context,
	in *pb.LikeListArgs) (*pb.LikeInfos, error) {
	return nil, nil
}

func (s *SvrHandler) GetLike(ctx context.Context,
	in *pb.LikeGetArgs) (*pb.LikeInfo, error) {
	return nil, nil
}

func (s *SvrHandler) NewLike(ctx context.Context,
	in *pb.LikeNewArgs) (*pb.LikeInfo, error) {
	return nil, nil
}

func (s *SvrHandler) DelLike(ctx context.Context,
	in *pb.LikeDelArgs) (*pb.ReplyBaseInfo, error) {
	return nil, nil
}

func (s *SvrHandler) GetsSummary(ctx context.Context,
	in *pb.BoardSummaryArgs) (*pb.BoardSummaryInfos, error) {
	return nil, nil
}
