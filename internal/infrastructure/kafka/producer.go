package kafka

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const flushTimeout = 5000

type Producer struct {
	producer *kafka.Producer
	l        *zap.SugaredLogger
}

func NewProducer(address []string, l *zap.SugaredLogger) (*Producer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	l.Info("Create new kafka producer")
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		l.Errorf("Create kafka producer error, NewProducer method: %v", err)
		return nil, fmt.Errorf("create kafka producer error, NewProducer method: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic, key string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:     []byte(message),
		Key:       []byte(key),
		Timestamp: time.Now(),
	}

	kafkaChan := make(chan kafka.Event)

	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		p.l.Errorf("Send message error, Produce method: %v", err)
		return fmt.Errorf("send message error, Produce method: %w", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		p.l.Errorf("Event control error, Produce method: %v", ev)
		return fmt.Errorf("event control error, Produce method: %w", ev)
	default:
		p.l.Errorf("Event control error, Produce method: %v", errors.New("unknown event type"))
		return fmt.Errorf("event control error, Produce method: %w", errors.New("unknown event type"))
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
