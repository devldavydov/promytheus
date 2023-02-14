package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/env"
)

const (
	defaultConfigAddress        = "127.0.0.1:8080"
	defaultConfigReportInterval = 10 * time.Second
	defaultConfigPollInterval   = 2 * time.Second
	defaultConfigLogLevel       = "DEBUG"
)

type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	LogLevel       string
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	config := &Config{}

	flagSet.StringVar(&config.Address, "a", defaultConfigAddress, "server address")
	flagSet.DurationVar(&config.ReportInterval, "r", defaultConfigReportInterval, "report interval")
	flagSet.DurationVar(&config.PollInterval, "p", defaultConfigPollInterval, "poll interval")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	err = flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	config.Address, err = env.GetVariable("ADDRESS", env.CastString, config.Address)
	if err != nil {
		return nil, err
	}

	config.ReportInterval, err = env.GetVariable("REPORT_INTERVAL", env.CastDuration, config.ReportInterval)
	if err != nil {
		return nil, err
	}

	config.PollInterval, err = env.GetVariable("POLL_INTERVAL", env.CastDuration, config.PollInterval)
	if err != nil {
		return nil, err
	}

	config.LogLevel, err = env.GetVariable("LOG_LEVEL", env.CastString, defaultConfigLogLevel)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func AgentSettingsAdapt(config *Config) (agent.ServiceSettings, error) {
	agentSettings, err := agent.NewServiceSettings(
		"http://"+config.Address,
		config.PollInterval,
		config.ReportInterval)
	if err != nil {
		return agent.ServiceSettings{}, err
	}
	return agentSettings, nil
}
