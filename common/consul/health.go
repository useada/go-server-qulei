package consul

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthImpl 健康检查实现
type HealthImpl struct{}

// Check 实现健康检查接口
func (h *HealthImpl) Check(ctx context.Context, args *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthImpl) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func RegisterGrpcHealth(svr *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(svr, &HealthImpl{})
}
