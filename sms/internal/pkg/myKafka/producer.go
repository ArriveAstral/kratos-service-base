package myKafka

import (
	"github.com/Shopify/sarama"
	"github.com/ZQCard/kratos-service-base/sms/internal/conf"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gopkg.in/yaml.v2"
	"time"
)

// 创建生产者
func CreateProducer(addrs []string, idempotent bool) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	// request.timeout.ms
	config.Producer.Timeout = time.Second * 5
	// request.required.acks
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = idempotent
	config.Version = sarama.V0_11_0_1

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return sarama.NewSyncProducer(addrs, config)
}

func SendSyncMessage(configPath, topic, key, value, headers, metadata string, offset int64, partition int32, idempotent bool) (partitionResult int32, offsetResult int64, err error) {
	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := c.Load(); err != nil {
		return 0, 0, err
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return 0, 0, err
	}

	producer, err := CreateProducer(bc.Data.Kafka.Addrs, idempotent)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			return
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(value),
		Headers:   nil,
		Metadata:  metadata,
		Offset:    offset,
		Partition: partition,
		Timestamp: time.Time{},
	}

	partitionResult, offsetResult, err = producer.SendMessage(msg)
	return
}

func CreateAsyncProducer(addrs []string, idempotent bool) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	// request.timeout.ms
	config.Producer.Timeout = time.Second * 5
	// request.required.acks
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = idempotent
	config.Version = sarama.V0_11_0_1

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return sarama.NewAsyncProducer(addrs, config)
}
