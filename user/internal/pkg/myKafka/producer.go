package myKafka

import (
	"github.com/Shopify/sarama"
	"time"
)

var Address []string

// 创建生产者
func CreateProducer(idempotent bool) (sarama.SyncProducer, error) {
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
	return sarama.NewSyncProducer([]string{"121.41.201.112:9092"}, config)
}

func SendSyncMessage(topic, key, value, headers, metadata string, offset int64, partition int32, idempotent bool) (partitionResult int32, offsetResult int64, err error) {

	producer, err := CreateProducer(idempotent)

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

func CreateAsyncProducer(idempotent bool) (sarama.AsyncProducer, error) {
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
	return sarama.NewAsyncProducer(Address, config)
}
