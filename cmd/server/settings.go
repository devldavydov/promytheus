package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/promytheus/internal/server"
)

type EnvConfig struct {
	Address  string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

func LoadEnvConfig() (EnvConfig, error) {
	envCfg := EnvConfig{}
	if err := env.Parse(&envCfg); err != nil {
		return EnvConfig{}, err
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

	return server.NewServiceSettings(address, port), nil
}
