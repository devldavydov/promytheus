package main

import (
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/promytheus/internal/agent"
)

type EnvConfig struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	LogLevel       string        `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

func LoadEnvConfig() (EnvConfig, error) {
	envCfg := EnvConfig{}
	if err := env.Parse(&envCfg); err != nil {
		return EnvConfig{}, err
	}

	return envCfg, nil
}

func AgentSettingsAdapt(envConfig EnvConfig) (agent.ServiceSettings, error) {
	agentSettings, err := agent.NewServiceSettings(
		"http://"+envConfig.Address,
		envConfig.PollInterval,
		envConfig.ReportInterval)
	if err != nil {
		return agent.ServiceSettings{}, err
	}
	return agentSettings, nil
}
