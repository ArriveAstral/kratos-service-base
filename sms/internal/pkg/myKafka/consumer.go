package myKafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ZQCard/kratos-service-base/sms/internal/conf"
	"github.com/ZQCard/kratos-service-base/sms/internal/data/model"
	aliSms "github.com/ZQCard/kratos-service-base/sms/internal/pkg/sms"
	"github.com/ZQCard/kratos-service-base/sms/internal/pkg/util/random"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var hosts []string
var mysqlSource string
var mysqlConn *gorm.DB

// 读取短信配置
func setSmsConfig() error {
	c := config.New(
		config.WithSource(
			file.NewSource("../../configs"),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return err
	}
	hosts = bc.Data.Kafka.Addrs
	mysqlSource = bc.Data.Database.Source
	return nil
}

func ConsumerRegisterSmsSend() error {
	topic := "user_register"
	if len(hosts) == 0 {
		if err := setSmsConfig(); err != nil {
			return err
		}
	}
	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_1
	config.Consumer.Offsets.AutoCommit.Enable = true
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration, error: %v", err)
	}

	consumer, err := sarama.NewConsumer(hosts, config)
	if err != nil {
		return err
	}
	defer consumer.Close()
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for _, partition := range partitions {
		// 读取当前最新offset
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		defer partitionConsumer.Close()

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
					msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
				if err := sendUserRegisterSmsCode(string(msg.Value), topic); err != nil {
					fmt.Printf("err :%s\n", err.Error())
				}

			case err := <-partitionConsumer.Errors():
				fmt.Printf("err :%s\n", err.Error())
			}
		}
	}
	return nil
}

func getMysqlConn() {
	if mysqlConn != nil {
		return
	}
	// mysql数据库连接
	db, err := gorm.Open(mysql.Open(mysqlSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mysqlConn = db
}

func sendUserRegisterSmsCode(mobile, scene string) error {
	getMysqlConn()
	// 验证码时长5分钟
	fiveMinuteAfter := time.Now().Add(time.Duration(5) * time.Minute)
	// 生成验证码
	code := random.GenerateNumber(6)
	err := mysqlConn.Model(&model.Sms{}).Create(&model.Sms{
		Id:         0,
		Mobile:     mobile,
		Content:    code,
		Type:       model.SmsTypeVerifyCode,
		Scene:      scene,
		IsExpire:   model.SmsIsExpireNO,
		ExpireTime: &fiveMinuteAfter,
	}).Error
	if err != nil {
		return err
	}
	// 发送短信验证码
	if err := aliSms.SendAliSmsVerifyCode(mobile, code); err != nil {
		return err
	}
	return nil
}
