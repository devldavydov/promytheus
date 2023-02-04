package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/promytheus/internal/server"
	"github.com/devldavydov/promytheus/internal/server/storage"
)

type EnvConfig struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	LogLevel      string        `env:"LOG_LEVEL" envDefault:"DEBUG"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string
	Restore       bool `env:"RESTORE" envDefault:"true"`
}

func LoadEnvConfig() (EnvConfig, error) {
	envCfg := EnvConfig{}
	if err := env.Parse(&envCfg); err != nil {
		return EnvConfig{}, err
	}

	// Get env var manualy, because caarlos0 set default value, when you set env var to empty value: STORE_FILE= cmd/server/server
	val, exists := os.LookupEnv("STORE_FILE")
	if !exists {
		envCfg.StoreFile = "/tmp/devops-metrics-db.json"
	} else {
		envCfg.StoreFile = val
	}

	return envCfg, nil
}

func ServerSettingsAdapt(envConfig EnvConfig) (server.ServiceSettings, error) {
	parts := strings.Split(envConfig.Address, ":")
	if len(parts) != 2 {
		return server.ServiceSettings{}, fmt.Errorf("wrong address format")
	}

	address := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return server.ServiceSettings{}, fmt.Errorf("wrong address format")
	}

	persistSettings := storage.NewPersistSettings(envConfig.StoreInterval, envConfig.StoreFile, envConfig.Restore)
	return server.NewServiceSettings(address, port, persistSettings), nil
}
