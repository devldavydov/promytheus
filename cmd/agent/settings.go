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
	_defaultConfigAddress        = "127.0.0.1:8080"
	_defaultConfigReportInterval = 10 * time.Second
	_defaultConfigPollInterval   = 2 * time.Second
	_defaultConfigLogLevel       = "DEBUG"
	_defaultConfigLogFile        = "agent.log"
	_defaultHmacKey              = ""
)

type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	HmacKey        string
	LogLevel       string
	LogFile        string
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	config := &Config{}

	flagSet.StringVar(&config.Address, "a", _defaultConfigAddress, "server address")
	flagSet.DurationVar(&config.ReportInterval, "r", _defaultConfigReportInterval, "report interval")
	flagSet.DurationVar(&config.PollInterval, "p", _defaultConfigPollInterval, "poll interval")
	flagSet.StringVar(&config.HmacKey, "k", _defaultHmacKey, "sign key")
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

	config.HmacKey, err = env.GetVariable("KEY", env.CastString, config.HmacKey)
	if err != nil {
		return nil, err
	}

	config.LogLevel, err = env.GetVariable("LOG_LEVEL", env.CastString, _defaultConfigLogLevel)
	if err != nil {
		return nil, err
	}

	config.LogFile, err = env.GetVariable("LOG_FILE", env.CastString, _defaultConfigLogFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func AgentSettingsAdapt(config *Config) (agent.ServiceSettings, error) {
	agentSettings, err := agent.NewServiceSettings(
		"http://"+config.Address,
		config.PollInterval,
		config.ReportInterval,
		config.HmacKey)
	if err != nil {
		return agent.ServiceSettings{}, err
	}
	return agentSettings, nil
}
