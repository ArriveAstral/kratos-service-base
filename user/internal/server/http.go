package server

import (
	"context"
	"fmt"
	v1 "github.com/ZQCard/kratos-service-base/api/user/v1"
	"github.com/ZQCard/kratos-service-base/user/internal/conf"
	"github.com/ZQCard/kratos-service-base/user/internal/pkg/middleware/jwt"
	"github.com/ZQCard/kratos-service-base/user/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func getOperation(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		if tr, ok := transport.FromServerContext(ctx); ok {
			fmt.Println(tr.Operation())
		}
		return handler(ctx, req)
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, userService *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			// 对于需要登录的路由进行jwt中间件验证
			selector.Server(jwt.AuthMiddleware()).
				//Prefix("/").
				Path(
					// 账户
					"/api.user.v1.User/GetUser",
					"/api.user.v1.User/CheckUserOK",

					// 地址
					"/api.user.v1.User/ListAddress",
					"/api.user.v1.User/CreateAddress",
					"/api.user.v1.User/GetAddress",
					"/api.user.v1.User/UpdateAddress",
					"/api.user.v1.User/DeleteAddress",
				).
				Build(),
			getOperation,
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterUserHTTPServer(srv, userService)
	return srv
}
