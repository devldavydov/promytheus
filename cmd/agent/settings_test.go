package main

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentSettingsAdaptDefault(t *testing.T) {
	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://127.0.0.1:8080")
	assert.Equal(t, 10*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 2*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustom(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("REPORT_INTERVAL", "1s")
	t.Setenv("POLL_INTERVAL", "2s")

	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	agentSettings, err := AgentSettingsAdapt(envCfg)
	assert.NoError(t, err)

	expURL, _ := url.Parse("http://1.1.1.1:9999")
	assert.Equal(t, 1*time.Second, agentSettings.GetReportInterval())
	assert.Equal(t, 2*time.Second, agentSettings.GetPollInterval())
	assert.Equal(t, expURL, agentSettings.GetServerAddress())
}

func TestAgentSettingsAdaptCustomError(t *testing.T) {
	t.Setenv("ADDRESS", "a.%^7b.c.d.e.f")

	envCfg, err := LoadEnvConfig()
	assert.NoError(t, err)

	_, err = AgentSettingsAdapt(envCfg)
	assert.Error(t, err)

}
