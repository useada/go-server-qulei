package tracing

import (
	"io"
	"time"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitTracing(serviceName string) (io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "127.0.0.1:6831",
		},
	}
	return cfg.InitGlobalTracer(serviceName)
}
