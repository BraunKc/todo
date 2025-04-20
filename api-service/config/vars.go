package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	Logger *zap.Logger
	Config config
)

type config struct {
	DBService struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"db-service"`
}

func InitYaml() {
	file, err := os.ReadFile("api-service/config/config.yml")
	if err != nil {
		Logger.Fatal("yaml init error", zap.Error(err))
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		Logger.Fatal("init yaml error", zap.Error(err))
	}

	Logger.Info("yaml inited", zap.Any("config", Config))
}
