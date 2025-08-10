package main

import (
	"fmt"
	"pet-kafka/internal/config"
	"pet-kafka/internal/infrastructure/kafka"
	"pet-kafka/internal/infrastructure/logger"

	"github.com/google/uuid"
)

const (
	topic = "my-topic"
	count = 10
)

func main() {
	logger.Init()
	l := logger.GetLogger()
	defer l.Sync()

	cfg, err := config.NewConfig(l)
	if err != nil {
		l.Fatalf("Create config error, RunProducer method: %v", err)
	}

	p, err := kafka.NewProducer(cfg.Addresses, l)
	if err != nil {
		l.Fatalf("Create new kafka producer error, RunProducer method: %v", err)
	}

	keys := generateUUIDString()
	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Message %d", i)
		key := keys[i%count]
		if err = p.Produce(msg, topic, key); err != nil {
			l.Errorf("Sending message error, RunProducer method: %v", err)
		}
	}
}

func generateUUIDString() [count]string {
	var uuids [count]string
	for i := 0; i < count; i++ {
		uuids[i] = uuid.NewString()
	}
	return uuids
}
