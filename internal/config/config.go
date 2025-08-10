package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	ServerConfig
	KafkaConfig
}

type ServerConfig struct {
	Port int `env:"APP_SERVER_PORT"`
}

type KafkaConfig struct {
	Addresses []string `env:"KAFKA_BOOTSTRAP_SERVERS, required"`
}

func NewConfig(l *zap.SugaredLogger) (*Config, error) {
	l.Info("Make new config")
	err := godotenv.Load()
	if err != nil {
		l.Errorf("Load env errof, NewConfig method: %w", err)
		return nil, fmt.Errorf("load env errof, NewConfig method: %w", err)
	}

	var config Config
	err = cleanenv.ReadEnv(&config)
	if err != nil {
		l.Errorf("Init config error, NewConfig method: %w", err)
		return nil, fmt.Errorf("init config error, NewConfig method: %w", err)
	}
	return &config, nil
}
