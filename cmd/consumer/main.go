package main

import (
	"os"
	"os/signal"
	"pet-kafka/internal/config"
	handle "pet-kafka/internal/handler"
	"pet-kafka/internal/infrastructure/kafka"
	"pet-kafka/internal/infrastructure/logger"
	"syscall"
)

const (
	topic         = "my-topic"
	consumerGroup = "my-consumer-group"
)

func main() {
	logger.Init()
	l := logger.GetLogger()
	defer l.Sync()
	cfg, _ := config.NewConfig(l)

	h := handle.NewHandler()
	c1, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 1)
	if err != nil {
		l.Fatal(err)
	}

	c2, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 2)
	if err != nil {
		l.Fatal(err)
	}

	c3, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 3)
	if err != nil {
		l.Fatal(err)
	}

	go func() {
		c1.Start()
	}()
	go func() {
		c2.Start()
	}()
	go func() {
		c3.Start()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	l.Fatal(c1.Stop(), c2.Stop(), c3.Stop())
}
