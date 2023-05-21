package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentSettingsAdaptDefault(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://127.0.0.1:8080")
	assert.Equal(t, 10*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 2*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
	assert.Nil(t, agentSettings.HmacKey)
	assert.Nil(t, agentSettings.CryptoPubKeyPath)
	assert.Equal(t, 2, agentSettings.RateLimit)
}

func TestAgentSettingsAdaptCustomEnv(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "1s")
	t.Setenv("POLL_INTERVAL", "2s")
	t.Setenv("KEY", "123")
	t.Setenv("RATE_LIMIT", "10")
	t.Setenv("CRYPTO_KEY", "/home/.ssh/id_rsa.pub")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 1*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 2*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
	assert.Equal(t, "123", *agentSettings.HmacKey)
	assert.Equal(t, "/home/.ssh/id_rsa.pub", *agentSettings.CryptoPubKeyPath)
	assert.Equal(t, 10, agentSettings.RateLimit)
}

func TestAgentSettingsAdaptCustomFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(
		*testFlagSet,
		[]string{"-a", "8.8.8.8:8888", "-r", "11s", "-p", "3s", "-k", "123", "-l", "5", "-crypto-key", "./key.pub"},
	)
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://8.8.8.8:8888")
	assert.Equal(t, 11*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 3*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
	assert.Equal(t, "123", *agentSettings.HmacKey)
	assert.Equal(t, "./key.pub", *agentSettings.CryptoPubKeyPath)
	assert.Equal(t, 5, agentSettings.RateLimit)
}

func TestAgentSettingsAdaptCustomEnvAndFlag(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "2s")
	t.Setenv("POLL_INTERVAL", "4s")
	t.Setenv("KEY", "123")
	t.Setenv("RATE_LIMIT", "15")
	t.Setenv("CRYPTO_KEY", "/home/.ssh/id_rsa.pub")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(
		*testFlagSet,
		[]string{"-a", "8.8.8.8:8888", "-r", "11s", "-p", "3s", "-k", "456", "-l", "1", "-crypto-key", "./key.pub"},
	)
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 2*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 4*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
	assert.Equal(t, "123", *agentSettings.HmacKey)
	assert.Equal(t, "/home/.ssh/id_rsa.pub", *agentSettings.CryptoPubKeyPath)
	assert.Equal(t, 15, agentSettings.RateLimit)
}

func TestAgentSettingsAdaptCustomEnvAndFlagMix(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("CRYPTO_KEY", "/home/.ssh/id_rsa.pub")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-p", "3s", "-l", "11"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 10*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 3*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
	assert.Nil(t, agentSettings.HmacKey)
	assert.Equal(t, "/home/.ssh/id_rsa.pub", *agentSettings.CryptoPubKeyPath)
	assert.Equal(t, 11, agentSettings.RateLimit)
}

func TestAgentSettingsAdaptCustomError(t *testing.T) {
	t.Setenv("ADDRESS", "a.%^7b.c.d.e.f")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	_, err = AgentSettingsAdapt(config)
	assert.Error(t, err)
}

func TestAgentSettingsCastEnvError(t *testing.T) {
	for i, tt := range []struct {
		envVarName string
		envVarVal  string
	}{
		{envVarName: "REPORT_INTERVAL", envVarVal: "foobar"},
		{envVarName: "POLL_INTERVAL", envVarVal: "foobar"},
		{envVarName: "RATE_LIMIT", envVarVal: "foobar"},
	} {
		tt := tt
		i := i
		t.Run(fmt.Sprintf("check%d", i), func(t *testing.T) {
			t.Setenv(tt.envVarName, tt.envVarVal)
			testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
			_, err := LoadConfig(*testFlagSet, []string{})
			assert.Error(t, err)
		})
	}
}
func TestAgentSettingsWithConfigFile(t *testing.T) {
	// Create temp config file
	fCfg, err := os.CreateTemp("", "cfg")
	require.NoError(t, err)

	defer func() {
		fCfg.Close()
		os.Remove(fCfg.Name())
	}()

	cfgAddr := "172.100.1.1:9090"
	cfgRepInt := 100 * time.Minute
	cfgPollInt := 200 * time.Minute

	tempCfg := configFile{
		Address:        &cfgAddr,
		ReportInterval: &cfgRepInt,
		PollInterval:   &cfgPollInt,
	}
	assert.NoError(t, json.NewEncoder(fCfg).Encode(&tempCfg))

	// Check
	t.Setenv("REPORT_INTERVAL", "1s")
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-p", "3s", "-c", fCfg.Name()})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://172.100.1.1:9090")
	assert.Equal(t, 1*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 3*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
	assert.Nil(t, agentSettings.HmacKey)
	assert.Nil(t, agentSettings.CryptoPubKeyPath)
	assert.Equal(t, 2, agentSettings.RateLimit)
}
