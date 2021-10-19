package data

import (
	"context"
	userv1 "github.com/ZQCard/kratos-service-base/api/user/v1"
	"github.com/ZQCard/kratos-service-base/coupon/internal/conf"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewCouponRepo,
	NewDiscovery,
	NewUserServiceClient,
)

// Data .
type Data struct {
	db         *gorm.DB
	userClient userv1.UserClient
}

// NewData .
func NewData(
	conf *conf.Data,
	logger log.Logger,
	userClient userv1.UserClient,
) (*Data, func(), error) {

	log := log.NewHelper(logger)
	// mysql数据库连接
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		db:         db,
		userClient: userClient,
	}

	return d, func() {
		log.Info("message", "closing the data resources")
	}, nil
}

func NewUserServiceClient(r registry.Discovery) userv1.UserClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///kratos-service-base.user.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := userv1.NewUserClient(conn)
	return c
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
