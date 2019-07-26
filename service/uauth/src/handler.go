package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"
)

func RegisterHandler(svr *grpc.Server) {
	pb.RegisterUauthServer(svr, &SvrHandler{})
}

type SvrHandler struct{}

func (s *SvrHandler) Login(ctx context.Context, in *pb.AuthLoginArgs) (*pb.AuthTokenInfo, error) {
	uid, err := checkLogin(ctx, in)
	fmt.Println(uid, err)

	return nil, nil
}

func (s *SvrHandler) Passwd(ctx context.Context, in *pb.AuthPasswdArgs) (*pb.AuthTokenInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Refresh(ctx context.Context, in *pb.AuthRefreshArgs) (*pb.AuthTokenInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Bind(ctx context.Context, in *pb.AuthBindArgs) (*pb.AuthUserInfo, error) {
	return nil, nil
}

func (s *SvrHandler) UnBind(ctx context.Context, in *pb.AuthUnBindArgs) (*pb.AuthUserInfo, error) {
	return nil, nil
}

func (s *SvrHandler) Detail(ctx context.Context, in *pb.AuthDetailArgs) (*pb.AuthUserInfo, error) {
	return nil, nil
}

func (s *SvrHandler) checkLogin(ctx context.Context, in *pb.AuthLoginArgs) (string, error) {
	switch in.Method {
	case pb.AuthMethod_PASSWD:
		return checkPasswd(ctx, in.Openid, in.Code)
	case pb.AuthMethod_SMS:
		return checkSmsCode(ctx, in.Openid, in.Code)
	case pb.AuthMethod_WECHAT:
		return checkWechat(in.Openid, in.Code)
	case pb.AuthMethod_QICQ:
		return checkQicq(in.Openid, in.Code)
	}
	return "", errors.New("login method don't support")
}

func (s *SvrHandler) checkPasswd(ctx context.Context, openid, code string) (string, error) {
	return "", nil
}

func (s *SvrHandler) checkSmsCode(ctx context.Context, openid, code string) (string, error) {
	return "", nil
}

func (s *SvrHandler) checkWechat(openid, code string) (string, error) {
	return "", nil
}

func (s *SvrHandler) checkQicq(openid, code string) (string, error) {
	return "", nil
}
