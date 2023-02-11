package main

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerSettingsAdaptDefault(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1", serverSettings.ServerAddress)
	assert.Equal(t, 8080, serverSettings.ServerPort)
	assert.Equal(t, 300*time.Second, serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/devops-metrics-db.json", serverSettings.PersistSettings.StoreFile)
	assert.True(t, serverSettings.PersistSettings.Restore)
}

func TestServerSettingsAdaptCustomEnv(t *testing.T) {
	testStoreFile := []string{"/foo/bar", ""}
	for _, storeFile := range testStoreFile {
		t.Setenv("ADDRESS", "1.1.1.1:9999")
		t.Setenv("STORE_INTERVAL", "0")
		t.Setenv("STORE_FILE", storeFile)
		t.Setenv("RESTORE", "false")

		testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
		config, err := LoadConfig(*testFlagSet, []string{})
		assert.NoError(t, err)

		serverSettings, err := ServerSettingsAdapt(config)
		assert.NoError(t, err)

		assert.Equal(t, "1.1.1.1", serverSettings.ServerAddress)
		assert.Equal(t, 9999, serverSettings.ServerPort)
		assert.Equal(t, time.Duration(0), serverSettings.PersistSettings.StoreInterval)
		assert.Equal(t, storeFile, serverSettings.PersistSettings.StoreFile)
		assert.False(t, serverSettings.PersistSettings.Restore)
	}
}

func TestServerSettingsAdaptCustomFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "1.1.1.1:9999", "-i", "0s", "-f", "/tmp/ttt", "-r=false"})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "1.1.1.1", serverSettings.ServerAddress)
	assert.Equal(t, 9999, serverSettings.ServerPort)
	assert.Equal(t, time.Duration(0), serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/ttt", serverSettings.PersistSettings.StoreFile)
	assert.False(t, serverSettings.PersistSettings.Restore)
}

func TestServerSettingsAdaptCustomEnvAndFlag(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("STORE_INTERVAL", "0")
	t.Setenv("STORE_FILE", "/tmp/ttt")
	t.Setenv("RESTORE", "false")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-a", "7.7.7.7:7777", "-i", "10s", "-f", "/tmp/aaa", "-r=true"})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "1.1.1.1", serverSettings.ServerAddress)
	assert.Equal(t, 9999, serverSettings.ServerPort)
	assert.Equal(t, time.Duration(0), serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/ttt", serverSettings.PersistSettings.StoreFile)
	assert.False(t, serverSettings.PersistSettings.Restore)
}

func TestServerSettingsAdaptCustomEnvAndFlagMix(t *testing.T) {
	t.Setenv("STORE_INTERVAL", "1s")
	t.Setenv("RESTORE", "true")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-i", "5m", "-r=false"})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1", serverSettings.ServerAddress)
	assert.Equal(t, 8080, serverSettings.ServerPort)
	assert.Equal(t, 1*time.Second, serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/devops-metrics-db.json", serverSettings.PersistSettings.StoreFile)
	assert.True(t, serverSettings.PersistSettings.Restore)
}

func TestServerSettingsAdaptCustomError(t *testing.T) {
	testAddress := []string{"1.1.1.1", "1.1.1.1:foobar"}

	for _, addr := range testAddress {
		t.Setenv("ADDRESS", addr)

		testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
		config, err := LoadConfig(*testFlagSet, []string{})
		assert.NoError(t, err)

		_, err = ServerSettingsAdapt(config)
		assert.Error(t, err)
	}
}
