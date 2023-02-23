package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/devldavydov/promytheus/internal/common/env"
	"github.com/devldavydov/promytheus/internal/server"
	"github.com/devldavydov/promytheus/internal/server/storage"
)

const (
	_defaultConfigAddress       = "127.0.0.1:8080"
	_defaultConfigLogLevel      = "DEBUG"
	_defaultConfigLogFile       = "server.log"
	_defaultconfigStoreInterval = 300 * time.Second
	_defaultConfigStoreFile     = "/tmp/devops-metrics-db.json"
	_defaultConfigRestore       = true
	_defaultHmacKey             = ""
	_defaultDatabaseDsn         = ""
)

type Config struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	HmacKey       string
	DatabaseDsn   string
	LogLevel      string
	LogFile       string
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	config := &Config{}

	// Check flags
	flagSet.StringVar(&config.Address, "a", _defaultConfigAddress, "server address")
	flagSet.DurationVar(&config.StoreInterval, "i", _defaultconfigStoreInterval, "store interval")
	flagSet.StringVar(&config.StoreFile, "f", _defaultConfigStoreFile, "store file")
	flagSet.BoolVar(&config.Restore, "r", _defaultConfigRestore, "restore")
	flagSet.StringVar(&config.HmacKey, "k", _defaultHmacKey, "sign key")
	flagSet.StringVar(&config.DatabaseDsn, "d", _defaultDatabaseDsn, "database dsn")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	err = flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	// Check env
	config.Address, err = env.GetVariable("ADDRESS", env.CastString, config.Address)
	if err != nil {
		return nil, err
	}

	config.StoreInterval, err = env.GetVariable("STORE_INTERVAL", env.CastDuration, config.StoreInterval)
	if err != nil {
		return nil, err
	}

	config.StoreFile, err = env.GetVariable("STORE_FILE", env.CastString, config.StoreFile)
	if err != nil {
		return nil, err
	}

	config.Restore, err = env.GetVariable("RESTORE", env.CastBool, config.Restore)
	if err != nil {
		return nil, err
	}

	config.HmacKey, err = env.GetVariable("KEY", env.CastString, config.HmacKey)
	if err != nil {
		return nil, err
	}

	config.DatabaseDsn, err = env.GetVariable("DATABASE_DSN", env.CastString, config.DatabaseDsn)
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

func ServerSettingsAdapt(config *Config) (server.ServiceSettings, error) {
	parts := strings.Split(config.Address, ":")
	if len(parts) != 2 {
		return server.ServiceSettings{}, fmt.Errorf("wrong address format")
	}

	address := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return server.ServiceSettings{}, fmt.Errorf("wrong address format")
	}

	persistSettings := storage.NewPersistSettings(config.StoreInterval, config.StoreFile, config.Restore)
	return server.NewServiceSettings(address, port, config.HmacKey, config.DatabaseDsn, persistSettings), nil
}
