package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/devldavydov/promytheus/internal/common/env"
	"github.com/devldavydov/promytheus/internal/common/settings"
	"github.com/devldavydov/promytheus/internal/server"
	"github.com/devldavydov/promytheus/internal/server/storage"
)

const (
	defaultConfigAddress       = "127.0.0.1:8080"
	defaultConfigLogLevel      = "DEBUG"
	defaultconfigStoreInterval = 300 * time.Second
	defaultConfigStoreFile     = "/tmp/devops-metrics-db.json"
	defaultConfigRestore       = true
)

type EnvConfig struct {
	Address       *env.EnvPair[string]
	LogLevel      *env.EnvPair[string]
	StoreInterval *env.EnvPair[time.Duration]
	StoreFile     *env.EnvPair[string]
	Restore       *env.EnvPair[bool]
}

type FlagConfig struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

func LoadEnvConfig() (*EnvConfig, error) {
	var err error
	envCfg := &EnvConfig{}

	envCfg.Address, err = env.GetVariable("ADDRESS", env.CastString, defaultConfigAddress)
	if err != nil {
		return nil, err
	}

	envCfg.LogLevel, err = env.GetVariable("LOG_LEVEL", env.CastString, defaultConfigLogLevel)
	if err != nil {
		return nil, err
	}

	envCfg.StoreInterval, err = env.GetVariable("STORE_INTERVAL", env.CastDuration, defaultconfigStoreInterval)
	if err != nil {
		return nil, err
	}

	envCfg.StoreFile, err = env.GetVariable("STORE_FILE", env.CastString, defaultConfigStoreFile)
	if err != nil {
		return nil, err
	}

	envCfg.Restore, err = env.GetVariable("RESTORE", env.CastBool, defaultConfigRestore)
	if err != nil {
		return nil, err
	}

	return envCfg, nil
}

func LoadFlagConfig(flagSet flag.FlagSet, flags []string) (*FlagConfig, error) {
	flagConfig := &FlagConfig{}
	flagSet.StringVar(&flagConfig.Address, "a", defaultConfigAddress, "server address")
	flagSet.DurationVar(&flagConfig.StoreInterval, "i", defaultconfigStoreInterval, "store interval")
	flagSet.StringVar(&flagConfig.StoreFile, "f", defaultConfigStoreFile, "store file")
	flagSet.BoolVar(&flagConfig.Restore, "r", defaultConfigRestore, "restore")
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

func ServerSettingsAdapt(envConfig *EnvConfig, flagConfig *FlagConfig) (server.ServiceSettings, error) {
	parts := strings.Split(settings.GetPriorityParam(envConfig.Address, flagConfig.Address), ":")
	if len(parts) != 2 {
		return server.ServiceSettings{}, fmt.Errorf("wrong address format")
	}

	address := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return server.ServiceSettings{}, fmt.Errorf("wrong address format")
	}

	persistSettings := storage.NewPersistSettings(
		settings.GetPriorityParam(envConfig.StoreInterval, flagConfig.StoreInterval),
		settings.GetPriorityParam(envConfig.StoreFile, flagConfig.StoreFile),
		settings.GetPriorityParam(envConfig.Restore, flagConfig.Restore))
	return server.NewServiceSettings(address, port, persistSettings), nil
}
