package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func StartDBSpan(ctx context.Context, dbtype, cmd string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		dbtype+" "+cmd,
	)

	ext.SpanKindRPCClient.Set(span)
	ext.DBType.Set(span, dbtype)
	ext.Component.Set(span, "golang/"+dbtype)
	return span
}
