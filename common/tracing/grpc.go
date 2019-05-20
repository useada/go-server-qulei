package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	gRPCComponentTag = opentracing.Tag{string(ext.Component), "gRPC"}
)

func GrpcClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}

		clientSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			gRPCComponentTag,
		)
		defer clientSpan.Finish()

		ctx = injectSpanContext(ctx, tracer, clientSpan)
		err := invoker(ctx, method, req, resp, cc, opts...)
		if err != nil {
			SetSpanTags(clientSpan, err, true)
			clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		return err
	}
}

func GrpcServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		spanContext, _ := extractSpanContext(ctx, tracer)

		serverSpan := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			gRPCComponentTag,
		)
		defer serverSpan.Finish()

		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		resp, err = handler(ctx, req)
		if err != nil {
			SetSpanTags(serverSpan, err, false)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		return resp, err
	}
}

func injectSpanContext(ctx context.Context, tracer opentracing.Tracer, cSpan opentracing.Span) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mdWriter := MetadataReaderWriter{md}

	// 注入信息到meta
	if err := tracer.Inject(cSpan.Context(), opentracing.HTTPHeaders, mdWriter); err != nil {
		cSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
	}
	// meta装入到context中
	return metadata.NewOutgoingContext(ctx, md)
}

func extractSpanContext(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return tracer.Extract(opentracing.HTTPHeaders, MetadataReaderWriter{md})
}
