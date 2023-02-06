package main

import (
	"flag"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentSettingsAdaptDefault(t *testing.T) {
	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagCfg, err := LoadFlagConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg, flagCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://127.0.0.1:8080")
	assert.Equal(t, 10*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 2*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustomEnv(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "1s")
	t.Setenv("POLL_INTERVAL", "2s")

	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagCfg, err := LoadFlagConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg, flagCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 1*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 2*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustomEnvAndFlag(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "2s")
	t.Setenv("POLL_INTERVAL", "4s")

	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagCfg, err := LoadFlagConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-r", "11s", "-p", "3s"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg, flagCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 2*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 4*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustomEnvAndFlagMix(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")

	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagCfg, err := LoadFlagConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-p", "3s"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg, flagCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 10*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 3*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustomFlag(t *testing.T) {
	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagCfg, err := LoadFlagConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-r", "11s", "-p", "3s"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg, flagCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://8.8.8.8:8888")
	assert.Equal(t, 11*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 3*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustomError(t *testing.T) {
	t.Setenv("ADDRESS", "a.%^7b.c.d.e.f")

	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flagCfg, err := LoadFlagConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	_, err = AgentSettingsAdapt(envCfg, flagCfg)
	assert.Error(t, err)

}
