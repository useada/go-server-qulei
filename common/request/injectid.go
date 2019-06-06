package request

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// InjectID gin middleware
func InjectID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.Request.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = xid.New().String()
		}
		ctx.Set("RequestID", requestID)
		ctx.Next()
	}
}

func GetRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value("RequestID").(string); ok {
		return id
	}
	return ""
}
