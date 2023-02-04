package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerSettingsAdaptDefault(t *testing.T) {
	envConfig, err := LoadEnvConfig()
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(envConfig)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1", serverSettings.GetServerAddress())
	assert.Equal(t, 8080, serverSettings.GetServerPort())
	assert.Equal(t, 300*time.Second, serverSettings.GetPersistenSettings().GetStoreInterval())
	assert.Equal(t, "/tmp/devops-metrics-db.json", serverSettings.GetPersistenSettings().GetStoreFile())
	assert.True(t, serverSettings.GetPersistenSettings().GetRestore())
}

func TestServerSettingsAdaptCustom(t *testing.T) {
	testStoreFile := []string{"/foo/bar", ""}
	for _, storeFile := range testStoreFile {
		t.Setenv("ADDRESS", "1.1.1.1:9999")
		t.Setenv("STORE_INTERVAL", "0")
		t.Setenv("STORE_FILE", storeFile)
		t.Setenv("RESTORE", "false")

		envConfig, err := LoadEnvConfig()
		assert.NoError(t, err)

		serverSettings, err := ServerSettingsAdapt(envConfig)
		assert.NoError(t, err)

		assert.Equal(t, "1.1.1.1", serverSettings.GetServerAddress())
		assert.Equal(t, 9999, serverSettings.GetServerPort())
		assert.Equal(t, time.Duration(0), serverSettings.GetPersistenSettings().GetStoreInterval())
		assert.Equal(t, storeFile, serverSettings.GetPersistenSettings().GetStoreFile())
		assert.False(t, serverSettings.GetPersistenSettings().GetRestore())
	}
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
