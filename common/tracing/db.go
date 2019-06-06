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

	ext.DBType.Set(span, dbtype)
	return span
}
