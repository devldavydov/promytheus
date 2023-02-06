package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/settings"
)

const (
	defaultConfigAddress        = "127.0.0.1:8080"
	defaultConfigReportInterval = 10 * time.Second
	defaultConfigPollInterval   = 2 * time.Second
	defaultConfigLogLevel       = "DEBUG"
)

type EnvConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	LogLevel       string        `env:"LOG_LEVEL"`
}

type FlagConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func LoadEnvConfig() (EnvConfig, error) {
	envCfg := EnvConfig{
		Address:        defaultConfigAddress,
		ReportInterval: defaultConfigReportInterval,
		PollInterval:   defaultConfigPollInterval,
		LogLevel:       defaultConfigLogLevel,
	}
	if err := env.Parse(&envCfg); err != nil {
		return EnvConfig{}, err
	}

	return envCfg, nil
}

func LoadFlagConfig(flagSet flag.FlagSet, flags []string) (FlagConfig, error) {
	flagParseConfig := struct {
		address        string
		reportInterval string
		pollInterval   string
	}{}
	flagSet.StringVar(&flagParseConfig.address, "a", defaultConfigAddress, "server address")
	flagSet.StringVar(&flagParseConfig.reportInterval, "r", defaultConfigReportInterval.String(), "report interval")
	flagSet.StringVar(&flagParseConfig.pollInterval, "p", defaultConfigPollInterval.String(), "poll interval")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	flagSet.Parse(flags)

	flagConfig := FlagConfig{Address: flagParseConfig.address}

	durVal, err := time.ParseDuration(flagParseConfig.reportInterval)
	if err != nil {
		return FlagConfig{}, err
	}
	flagConfig.ReportInterval = durVal

	durVal, err = time.ParseDuration(flagParseConfig.pollInterval)
	if err != nil {
		return FlagConfig{}, err
	}
	flagConfig.PollInterval = durVal

	return flagConfig, nil
}

func AgentSettingsAdapt(envConfig EnvConfig, flagConfig FlagConfig) (agent.ServiceSettings, error) {
	agentSettings, err := agent.NewServiceSettings(
		"http://"+settings.GetPriorityParam(envConfig.Address, flagConfig.Address, defaultConfigAddress),
		settings.GetPriorityParam(envConfig.PollInterval, flagConfig.PollInterval, defaultConfigPollInterval),
		settings.GetPriorityParam(envConfig.ReportInterval, flagConfig.ReportInterval, defaultConfigReportInterval))
	if err != nil {
		return agent.ServiceSettings{}, err
	}
	return agentSettings, nil
}
