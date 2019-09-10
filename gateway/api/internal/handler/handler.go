package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"a.com/go-server/common/tracing"
	"a.com/go-server/gclient"
)

type Handler struct {
	Grpc *gclient.Client
	Log  *zap.SugaredLogger
}

func NewApiHandler(gclient *gclient.Client, log *zap.SugaredLogger) *Handler {
	return &Handler{
		Grpc: gclient,
		Log:  log,
	}
}

func (h *Handler) TraceId(ctx *gin.Context) string {
	return tracing.GetID(ctx.Request.Context())
}
