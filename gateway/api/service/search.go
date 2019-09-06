package service

import (
	"github.com/gin-gonic/gin"

	"a.com/go-server/gateway/api/base"
)

func (s *Service) UsersByName(ctx *gin.Context) *base.JSONResponse {
	return base.SuccessResponse("")
}

func (s *Service) UsersByNear(ctx *gin.Context) *base.JSONResponse {
	return base.SuccessResponse("")
}
