package main

import (
	"context"
	"os"
	"os/signal"
	"pet-kafka/internal/config"
	handle "pet-kafka/internal/handler"
	"pet-kafka/internal/infrastructure/kafka"
	"pet-kafka/internal/infrastructure/logger"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
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

	g, ctx := errgroup.WithContext(context.Background())

	h := handle.NewHandler()
	c1, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 1)
	if err != nil {
		l.Errorf("First consumer start error: %v", err)
	}

	c2, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 2)
	if err != nil {
		l.Errorf("Second consumer start error: %v", err)
	}

	c3, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 3)
	if err != nil {
		l.Errorf("Third consumer start error: %v", err)
	}

	g.Go(func() error {
		return c1.Start(ctx)
	})
	ctxTemp, _ := context.WithTimeout(ctx, time.Millisecond*5)
	g.Go(func() error {
		return c2.Start(ctxTemp)
	})
	g.Go(func() error {
		return c3.Start(ctx)
	})

	if err := g.Wait(); err != nil {
		l.Errorf("Wait errGroup error: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	l.Fatal(c1.Stop(), c2.Stop(), c3.Stop())
}
