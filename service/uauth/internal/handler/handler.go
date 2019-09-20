package handler

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"a.com/go-server/proto/pb"
	"a.com/go-server/service/uauth/internal/cache"
	"a.com/go-server/service/uauth/internal/store"
)

type SvrHandler struct {
	Cache cache.Cache
	Store store.Store
	Log   *zap.SugaredLogger
}

func RegisterHandler(svr *grpc.Server, cache cache.Cache, store store.Store, log *zap.SugaredLogger) {
	pb.RegisterUauthServer(svr, &SvrHandler{
		Cache: cache,
		Store: store,
		Log:   log,
	})
}

func (s *SvrHandler) Login(ctx context.Context, in *pb.AuthLoginArgs) (*pb.AuthTokenInfo, error) {
	uid, err := s.checkLogin(ctx, in)
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
	var (
		uid string
		err error
	)
	switch in.Method {
	case pb.AuthMethod_PASSWD:
		uid, err = s.byPasswd(ctx, in.Openid, in.Code)
	case pb.AuthMethod_SMS:
		uid, err = s.bySmsCode(ctx, in.Openid, in.Code)
	case pb.AuthMethod_WECHAT:
		uid, err = s.byWechat(ctx, in.Openid, in.Code)
	case pb.AuthMethod_QICQ:
		uid, err = s.byQicq(ctx, in.Openid, in.Code)
	default:
		uid, err = "", errors.New("method not support")
	}
	return uid, err
}

func (s *SvrHandler) byPasswd(ctx context.Context, openid, code string) (string, error) {
	return "", nil
}

func (s *SvrHandler) bySmsCode(ctx context.Context, openid, code string) (string, error) {
	return "", nil
}

func (s *SvrHandler) byWechat(ctx context.Context, openid, code string) (string, error) {
	return "", nil
}

func (s *SvrHandler) byQicq(ctx context.Context, openid, code string) (string, error) {
	return "", nil
}
