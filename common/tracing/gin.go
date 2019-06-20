package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/xid"
)

// Trace gin middleware
func Trace(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := xid.New().String()
		span := tracer.StartSpan("HTTP/" + ctx.Request.Method)
		ext.HTTPMethod.Set(span, ctx.Request.Method)
		span.SetTag("trace.id", traceID)

		newCtx := context.WithValue(ctx.Request.Context(), "TraceID", traceID)
		ctx.Request = ctx.Request.WithContext(opentracing.ContextWithSpan(newCtx, span))

		ctx.Next()

		ext.HTTPStatusCode.Set(span, uint16(ctx.Writer.Status()))
		span.Finish()
	}
}

func GetID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	id, _ := ctx.Value("TraceID").(string)
	return id
}
