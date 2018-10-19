package main

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

var User *UserHandler

func (u *UserHandler) SearchByName(ctx *gin.Context) *JsonResponse {
	return SuccessResponse("")
}

func (u *UserHandler) SearchByNear(ctx *gin.Context) *JsonResponse {
	return SuccessResponse("")
}
