package main

import (
	"flag"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
}

func TestAgentSettingsAdaptCustomEnv(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "1s")
	t.Setenv("POLL_INTERVAL", "2s")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 1*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 2*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
}

func TestAgentSettingsAdaptCustomFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-r", "11s", "-p", "3s"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://8.8.8.8:8888")
	assert.Equal(t, 11*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 3*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
}

func TestAgentSettingsAdaptCustomEnvAndFlag(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "2s")
	t.Setenv("POLL_INTERVAL", "4s")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-r", "11s", "-p", "3s"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 2*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 4*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
}

func TestAgentSettingsAdaptCustomEnvAndFlagMix(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "8.8.8.8:8888", "-p", "3s"})
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(config)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 10*time.Second, agentSettings.ReportInterval)
	assert.Equal(t, 3*time.Second, agentSettings.PollInterval)
	assert.Equal(t, expURL, agentSettings.ServerAddress)
}

func TestAgentSettingsAdaptCustomError(t *testing.T) {
	t.Setenv("ADDRESS", "a.%^7b.c.d.e.f")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	_, err = AgentSettingsAdapt(config)
	assert.Error(t, err)

}
