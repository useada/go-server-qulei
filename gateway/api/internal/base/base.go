package base

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

// CORS ...
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
