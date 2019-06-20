package main

import (
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
}

var Search *SearchHandler

func (s *SearchHandler) UsersByName(ctx *gin.Context) *JSONResponse {
	return SuccessResponse("")
}

func (s *SearchHandler) UsersByNear(ctx *gin.Context) *JSONResponse {
	return SuccessResponse("")
}
