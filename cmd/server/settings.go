package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
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
	_defaultCryptoPrivKeyPath   = ""
	_defaultConfigFilePath      = ""
	_defaultTrustedSubnet       = ""
)

type Config struct {
	Address           string
	StoreFile         string
	HmacKey           string
	DatabaseDsn       string
	LogLevel          string
	LogFile           string
	CryptoPrivKeyPath string
	TrustedSubnet     string
	StoreInterval     time.Duration
	Restore           bool
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	var configFilePath string
	config := &Config{}

	// Check flags
	flagSet.StringVar(&config.Address, "a", _defaultConfigAddress, "server address")
	flagSet.DurationVar(&config.StoreInterval, "i", _defaultconfigStoreInterval, "store interval")
	flagSet.StringVar(&config.StoreFile, "f", _defaultConfigStoreFile, "store file")
	flagSet.BoolVar(&config.Restore, "r", _defaultConfigRestore, "restore")
	flagSet.StringVar(&config.HmacKey, "k", _defaultHmacKey, "sign key")
	flagSet.StringVar(&config.DatabaseDsn, "d", _defaultDatabaseDsn, "database dsn")
	flagSet.StringVar(&config.CryptoPrivKeyPath, "crypto-key", _defaultCryptoPrivKeyPath, "crypto private key path")
	flagSet.StringVar(&config.TrustedSubnet, "t", _defaultTrustedSubnet, "trusted subnet")
	//
	flagSet.StringVar(&configFilePath, "c", _defaultConfigFilePath, "config file path")
	flagSet.StringVar(&configFilePath, "config", _defaultConfigFilePath, "config file path")
	//
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

	config.CryptoPrivKeyPath, err = env.GetVariable("CRYPTO_KEY", env.CastString, config.CryptoPrivKeyPath)
	if err != nil {
		return nil, err
	}

	config.TrustedSubnet, err = env.GetVariable("TRUSTED_SUBNET", env.CastString, config.TrustedSubnet)
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

	//
	configFilePath, err = env.GetVariable("CONFIG", env.CastString, configFilePath)
	if err != nil {
		return nil, err
	}

	if err = applyConfigFile(config, configFilePath); err != nil {
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

	var trustedSubnet *net.IPNet
	if config.TrustedSubnet != "" {
		_, trustedSubnet, err = net.ParseCIDR(config.TrustedSubnet)
		if err != nil {
			return server.ServiceSettings{}, err
		}
	}

	persistSettings := storage.NewPersistSettings(config.StoreInterval, config.StoreFile, config.Restore)
	return server.NewServiceSettings(
		address,
		port,
		config.HmacKey,
		config.DatabaseDsn,
		persistSettings,
		config.CryptoPrivKeyPath,
		trustedSubnet), nil
}

type configFile struct {
	Address           *string        `json:"address"`
	Restore           *bool          `json:"restore"`
	StoreInterval     *time.Duration `json:"store_interval"`
	StoreFile         *string        `json:"store_file"`
	DatabaseDsn       *string        `json:"database_dsn"`
	HmacKey           *string        `json:"hmac_key"`
	CryptoPrivKeyPath *string        `json:"crypto_key"`
	TrustedSubnet     *string        `json:"trusted_subnet"`
}

func applyConfigFile(config *Config, configFilePath string) error {
	if configFilePath == "" {
		return nil
	}

	f, err := os.OpenFile(configFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	configFromFile := configFile{}
	if err = json.NewDecoder(f).Decode(&configFromFile); err != nil {
		return err
	}

	if configFromFile.Address != nil && config.Address == _defaultConfigAddress {
		config.Address = *configFromFile.Address
	}
	if configFromFile.Restore != nil && config.Restore {
		config.Restore = *configFromFile.Restore
	}
	if configFromFile.StoreInterval != nil && config.StoreInterval == _defaultconfigStoreInterval {
		config.StoreInterval = *configFromFile.StoreInterval
	}
	if configFromFile.StoreFile != nil && config.StoreFile == _defaultConfigStoreFile {
		config.StoreFile = *configFromFile.StoreFile
	}
	if configFromFile.DatabaseDsn != nil && config.DatabaseDsn == _defaultDatabaseDsn {
		config.DatabaseDsn = *configFromFile.DatabaseDsn
	}
	if configFromFile.HmacKey != nil && config.HmacKey == _defaultHmacKey {
		config.HmacKey = *configFromFile.HmacKey
	}
	if configFromFile.CryptoPrivKeyPath != nil && config.CryptoPrivKeyPath == _defaultCryptoPrivKeyPath {
		config.CryptoPrivKeyPath = *configFromFile.CryptoPrivKeyPath
	}
	if configFromFile.TrustedSubnet != nil && config.TrustedSubnet == _defaultTrustedSubnet {
		config.TrustedSubnet = *configFromFile.TrustedSubnet
	}

	return nil
}
