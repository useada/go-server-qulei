package main

import (
	"context"

	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterUauthServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) Login(ctx context.Context,
	in *pb.AuthLoginArgs) (*pb.AuthTokenInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Passwd(ctx context.Context,
	in *pb.AuthPasswdArgs) (*pb.AuthTokenInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Refresh(ctx context.Context,
	in *pb.AuthRefreshArgs) (*pb.AuthTokenInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Bind(ctx context.Context,
	in *pb.AuthBindArgs) (*pb.AuthUserInfo, error) {
	return nil, nil
}

func (s *SvrHandler) UnBind(ctx context.Context,
	in *pb.AuthUnBindArgs) (*pb.AuthUserInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Get(ctx context.Context,
	in *pb.AuthGetArgs) (*pb.AuthUserInfo, error) {
	return nil, nil
}
