package main

import (
	"os"
	"pet-kafka/internal/infrastructure/kafka"
	"pet-kafka/internal/infrastructure/logger"
)

func main() {
	logger.Init()
	l := logger.GetLogger()
	defer l.Sync()

	pidPath := os.Args[1]
	l.Info(pidPath)
	err := kafka.StopConsumer(pidPath, l)
	if err != nil {
		l.Fatal("Stop consumer error")
	}
}
