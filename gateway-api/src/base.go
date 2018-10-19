package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ApiHandler func(ctx *gin.Context) *JsonResponse

func ResponseWrapper(handle ApiHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, handle(ctx))
	}
}

func ErrorResponse(code, msg string) *JsonResponse {
	return &JsonResponse{Code: code, Msg: msg}
}

func SuccessResponse(data interface{}) *JsonResponse {
	return &JsonResponse{Code: "2000", Msg: "OK", Data: data}
}
