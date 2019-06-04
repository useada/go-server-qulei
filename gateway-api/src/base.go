package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONResponse response info struct
type JSONResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// APIHandler api function handler
type APIHandler func(ctx *gin.Context) *JSONResponse

// ResponseWrapper response wrapper
func ResponseWrapper(handle APIHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, handle(ctx))
	}
}

// ErrorResponse error response
func ErrorResponse(code, msg string) *JSONResponse {
	return &JSONResponse{Code: code, Msg: msg}
}

// SuccessResponse success response
func SuccessResponse(data interface{}) *JSONResponse {
	return &JSONResponse{Code: "2000", Msg: "OK", Data: data}
}
