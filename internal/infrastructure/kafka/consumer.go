package kafka

import (
	"context"
	"fmt"
	handle "pet-kafka/internal/handler"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	sessionTimeout = 7000
	noTimeout      = -1
)

type Consumer struct {
	consumer       *kafka.Consumer
	handler        handle.Handler
	l              *zap.SugaredLogger
	consumerNumber int
}

func NewConsumer(handler handle.Handler, l *zap.SugaredLogger, address []string, topic, consumerGroup string, consumerNumber int) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		l.Errorf("New consumer error, NewConsumer method: %v", err)
		return nil, fmt.Errorf("new consumer error, NewConsumer method: %w", err)
	}

	if err = c.Subscribe(topic, nil); err != nil {
		l.Errorf("Subscribe error, NewConsumer method: %w", err)
		return nil, fmt.Errorf("subscribe error, NewConsumer method: %w", err)
	}
	return &Consumer{
		consumer:       c,
		handler:        handler,
		l:              l,
		consumerNumber: consumerNumber,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
			if err != nil {
				c.l.Errorf("Read message error, Start method: %v", err)
				return fmt.Errorf("read message error, Start method: %w", err)
			}
			if kafkaMsg == nil {
				continue
			}
			if err = c.handler.HandlerImplMessage(kafkaMsg.Value, kafkaMsg.TopicPartition, c.consumerNumber); err != nil {
				c.l.Errorf("HandlerImpl message error, Start method: %v", err)
				return fmt.Errorf("handlerImpl message error, Start method: %w", err)
			}
			if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
				c.l.Errorf("Store message error, Start method: %v")
				return fmt.Errorf("store message error, Start method: %w", err)
			}
		}
	}
}

func (c *Consumer) Stop() error {
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	c.l.Infof("Commited offset")
	return c.consumer.Close()
}
