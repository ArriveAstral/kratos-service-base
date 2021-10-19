package data

import (
	"context"
	smsv1 "github.com/ZQCard/kratos-service-base/api/sms/v1"
	"github.com/ZQCard/kratos-service-base/user/internal/conf"
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
	NewDiscovery,
	NewUserRepo,
	NewAddressRepo,
	NewSmsServiceClient,
)

// Data .
type Data struct {
	db        *gorm.DB
	smsClient smsv1.SmsClient
}

// NewData .
func NewData(
	conf *conf.Data,
	logger log.Logger,
	smsClient smsv1.SmsClient,
) (*Data, func(), error) {
	log := log.NewHelper(logger)
	// mysql数据库连接
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		db:        db.Debug(),
		smsClient: smsClient,
	}

	return d, func() {
		log.Info("message", "closing the data resources")
	}, nil
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

func NewSmsServiceClient(r registry.Discovery) smsv1.SmsClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///kratos-service-base.sms.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := smsv1.NewSmsClient(conn)
	return c
}
