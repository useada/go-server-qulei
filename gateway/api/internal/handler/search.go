package handler

import (
	"github.com/gin-gonic/gin"

	"a.com/go-server/gateway/api/internal/base"
)

func (h *Handler) UsersByName(ctx *gin.Context) *base.JSONResponse {
	return base.SuccessResponse("")
}

func (h *Handler) UsersByNear(ctx *gin.Context) *base.JSONResponse {
	return base.SuccessResponse("")
}
