package main

import (
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
}

var Search *SearchHandler

func (s *SearchHandler) UsersByName(ctx *gin.Context) *JsonResponse {
	return SuccessResponse("")
}

func (s *SearchHandler) UsersByNear(ctx *gin.Context) *JsonResponse {
	return SuccessResponse("")
}
