package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Trace gin middleware
func Trace(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := tracer.StartSpan("HTTP/" + c.Request.Method)
		ext.HTTPMethod.Set(span, c.Request.Method)

		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), span))

		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		span.Finish()
	}
}
