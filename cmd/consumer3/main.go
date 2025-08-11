package main

import (
	"context"
	"os"
	"os/signal"
	"pet-kafka/internal/config"
	handle "pet-kafka/internal/handler"
	"pet-kafka/internal/infrastructure/kafka"
	"pet-kafka/internal/infrastructure/logger"
	"strconv"
	"syscall"

	"golang.org/x/sync/errgroup"
)

const (
	topic         = "my-topic"
	consumerGroup = "my-consumer-group3"
)

func main() {
	logger.Init()
	l := logger.GetLogger()
	defer l.Sync()

	cfg, _ := config.NewConfig(l)

	pidFile := "consumer3.pid"
	err := os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		l.Errorf("Create pidFile Error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	h := handle.NewHandler()
	c3, err := kafka.NewConsumer(*h, l, cfg.Addresses, topic, consumerGroup, 1)
	if err != nil {
		l.Errorf("Third consumer start error: %v", err)
	}

	l.Infof("Start third consumer in group %s", consumerGroup)
	g.Go(func() error {
		return c3.Start(ctx)
	})

	if err := g.Wait(); err != nil {
		l.Errorf("Wait errGroup error: %v", err)
	}

	if err = os.Remove(pidFile); err != nil {
		l.Errorf("Remove pid file error")
	}

	l.Infof("Consumer 3 stop gracefully")
}
