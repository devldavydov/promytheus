package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSettingsAdaptDefault(t *testing.T) {
	envConfig, err := LoadEnvConfig()
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(envConfig)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1", serverSettings.GetServerAddress())
	assert.Equal(t, 8080, serverSettings.GetServerPort())
}

func TestServerSettingsAdaptCustom(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")

	envConfig, err := LoadEnvConfig()
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(envConfig)
	assert.NoError(t, err)

	assert.Equal(t, "1.1.1.1", serverSettings.GetServerAddress())
	assert.Equal(t, 9999, serverSettings.GetServerPort())
}

func TestServerSettingsAdaptCustomError(t *testing.T) {
	testAddress := []string{"1.1.1.1", "1.1.1.1:foobar"}

	for _, addr := range testAddress {
		t.Setenv("ADDRESS", addr)

		envConfig, err := LoadEnvConfig()
		assert.NoError(t, err)

		_, err = ServerSettingsAdapt(envConfig)
		assert.Error(t, err)
	}
}
