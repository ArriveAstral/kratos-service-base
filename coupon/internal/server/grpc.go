package server

import (
	v1 "github.com/ZQCard/kratos-service-base/api/coupon/v1"
	"github.com/ZQCard/kratos-service-base/coupon/internal/conf"
	"github.com/ZQCard/kratos-service-base/coupon/internal/pkg/middleware/jwt"
	"github.com/ZQCard/kratos-service-base/coupon/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, service *service.CouponService, tp *tracesdk.TracerProvider, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			getOperation,
			// 对于需要登录的路由进行jwt中间件验证
			selector.Server(jwt.AuthMiddleware()).Build(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterCouponServer(srv, service)
	return srv
}
