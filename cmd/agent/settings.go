package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/env"
	"github.com/devldavydov/promytheus/internal/common/settings"
)

const (
	defaultConfigAddress        = "127.0.0.1:8080"
	defaultConfigReportInterval = 10 * time.Second
	defaultConfigPollInterval   = 2 * time.Second
	defaultConfigLogLevel       = "DEBUG"
)

type EnvConfig struct {
	Address        *env.EnvPair[string]
	ReportInterval *env.EnvPair[time.Duration]
	PollInterval   *env.EnvPair[time.Duration]
	LogLevel       *env.EnvPair[string]
}

type FlagConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func LoadEnvConfig() (*EnvConfig, error) {
	var err error
	envCfg := &EnvConfig{}

	envCfg.Address, err = env.GetVariable("ADDRESS", env.CastString, defaultConfigAddress)
	if err != nil {
		return nil, err
	}

	envCfg.ReportInterval, err = env.GetVariable("REPORT_INTERVAL", env.CastDuration, defaultConfigReportInterval)
	if err != nil {
		return nil, err
	}

	envCfg.PollInterval, err = env.GetVariable("POLL_INTERVAL", env.CastDuration, defaultConfigPollInterval)
	if err != nil {
		return nil, err
	}

	envCfg.LogLevel, err = env.GetVariable("LOG_LEVEL", env.CastString, defaultConfigLogLevel)
	if err != nil {
		return nil, err
	}

	return envCfg, nil
}

func LoadFlagConfig(flagSet flag.FlagSet, flags []string) (*FlagConfig, error) {
	flagConfig := &FlagConfig{}
	flagSet.StringVar(&flagConfig.Address, "a", defaultConfigAddress, "server address")
	flagSet.DurationVar(&flagConfig.ReportInterval, "r", defaultConfigReportInterval, "report interval")
	flagSet.DurationVar(&flagConfig.PollInterval, "p", defaultConfigPollInterval, "poll interval")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	err := flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	return flagConfig, nil
}

func AgentSettingsAdapt(envConfig *EnvConfig, flagConfig *FlagConfig) (agent.ServiceSettings, error) {
	agentSettings, err := agent.NewServiceSettings(
		"http://"+settings.GetPriorityParam(envConfig.Address, flagConfig.Address),
		settings.GetPriorityParam(envConfig.PollInterval, flagConfig.PollInterval),
		settings.GetPriorityParam(envConfig.ReportInterval, flagConfig.ReportInterval))
	if err != nil {
		return agent.ServiceSettings{}, err
	}
	return agentSettings, nil
}
