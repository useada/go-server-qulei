package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"a.com/go-server/common/tracing"
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

func (s *Service) TraceId(ctx *gin.Context) string {
	return tracing.GetID(ctx.Request.Context())
}
