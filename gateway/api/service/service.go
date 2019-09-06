package service

import (
	"go.uber.org/zap"

	"a.com/go-server/gclient"
)

func NewApiService(gclient *gclient.Client, log *zap.SugaredLogger) *Service {
	return &Service{
		Grpc: gclient,
		Log:  log,
	}
}

type Service struct {
	Grpc *gclient.Client
	Log  *zap.SugaredLogger
}
