package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/devldavydov/promytheus/internal/agent"
	"github.com/devldavydov/promytheus/internal/common/env"
)

const (
	_defaultConfigAddress          = "127.0.0.1:8080"
	_defaultConfigReportInterval   = 10 * time.Second
	_defaultConfigPollInterval     = 2 * time.Second
	_defaultConfigLogLevel         = "DEBUG"
	_defaultConfigLogFile          = "agent.log"
	_defaultConfigHmacKey          = ""
	_defaultConfigRateLimit        = 2
	_defaultConfigCryptoPubKeyPath = ""
	_defaultConfigFilePath         = ""
	_defaultConfigUseGRPC          = false
	_defaultConfigGRPCCACertPath   = ""
)

type Config struct {
	Address          string
	HmacKey          string
	LogLevel         string
	LogFile          string
	CryptoPubKeyPath string
	ReportInterval   time.Duration
	PollInterval     time.Duration
	RateLimit        int
	UseGRPC          bool
	GRPCCACertPath   string
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	var configFilePath string
	config := &Config{}

	flagSet.StringVar(&config.Address, "a", _defaultConfigAddress, "server address")
	flagSet.DurationVar(&config.ReportInterval, "r", _defaultConfigReportInterval, "report interval")
	flagSet.DurationVar(&config.PollInterval, "p", _defaultConfigPollInterval, "poll interval")
	flagSet.StringVar(&config.HmacKey, "k", _defaultConfigHmacKey, "sign key")
	flagSet.IntVar(&config.RateLimit, "l", _defaultConfigRateLimit, "rate limit")
	flagSet.StringVar(&config.CryptoPubKeyPath, "crypto-key", _defaultConfigCryptoPubKeyPath, "crypto public key path")
	flagSet.BoolVar(&config.UseGRPC, "g", _defaultConfigUseGRPC, "use gRPC insted of HTTP")
	flagSet.StringVar(&config.GRPCCACertPath, "gca", _defaultConfigGRPCCACertPath, "gRPC TLS CA certificate path")
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

	config.RateLimit, err = env.GetVariable("RATE_LIMIT", env.CastInt, config.RateLimit)
	if err != nil {
		return nil, err
	}

	config.CryptoPubKeyPath, err = env.GetVariable("CRYPTO_KEY", env.CastString, config.CryptoPubKeyPath)
	if err != nil {
		return nil, err
	}

	config.UseGRPC, err = env.GetVariable("USE_GRPC", env.CastBool, config.UseGRPC)
	if err != nil {
		return nil, err
	}

	config.GRPCCACertPath, err = env.GetVariable("GRPC_CA_CERT", env.CastString, config.GRPCCACertPath)
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

func AgentSettingsAdapt(config *Config) (agent.ServiceSettings, error) {
	agentSettings, err := agent.NewServiceSettings(
		config.Address,
		config.PollInterval,
		config.ReportInterval,
		config.HmacKey,
		config.RateLimit,
		config.CryptoPubKeyPath,
		config.UseGRPC,
		config.GRPCCACertPath)
	if err != nil {
		return agent.ServiceSettings{}, err
	}
	return agentSettings, nil
}

type configFile struct {
	Address          *string        `json:"address"`
	ReportInterval   *time.Duration `json:"report_interval"`
	PollInterval     *time.Duration `json:"poll_interval"`
	HmacKey          *string        `json:"hmac_key"`
	RateLimit        *int           `json:"rate_limit"`
	CryptoPubKeyPath *string        `json:"crypto_key"`
	UseGRPC          *bool          `json:"use_grpc"`
	GRPCCACertPath   *string        `json:"grpc_ca_cert"`
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
	if configFromFile.ReportInterval != nil && config.ReportInterval == _defaultConfigReportInterval {
		config.ReportInterval = *configFromFile.ReportInterval
	}
	if configFromFile.PollInterval != nil && config.PollInterval == _defaultConfigPollInterval {
		config.PollInterval = *configFromFile.PollInterval
	}
	if configFromFile.HmacKey != nil && config.HmacKey == _defaultConfigHmacKey {
		config.HmacKey = *configFromFile.HmacKey
	}
	if configFromFile.RateLimit != nil && config.RateLimit == _defaultConfigRateLimit {
		config.RateLimit = *configFromFile.RateLimit
	}
	if configFromFile.CryptoPubKeyPath != nil && config.CryptoPubKeyPath == _defaultConfigCryptoPubKeyPath {
		config.CryptoPubKeyPath = *configFromFile.CryptoPubKeyPath
	}
	if configFromFile.UseGRPC != nil && !config.UseGRPC {
		config.UseGRPC = *configFromFile.UseGRPC
	}
	if configFromFile.GRPCCACertPath != nil && config.GRPCCACertPath == _defaultConfigGRPCCACertPath {
		config.GRPCCACertPath = *configFromFile.GRPCCACertPath
	}

	return nil
}
