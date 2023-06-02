package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerSettingsAdaptDefault(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1", serverSettings.HttpSettings.Host)
	assert.Equal(t, 8080, serverSettings.HttpSettings.Port)
	assert.Nil(t, serverSettings.HmacKey)
	assert.Equal(t, 300*time.Second, serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/devops-metrics-db.json", serverSettings.PersistSettings.StoreFile)
	assert.Equal(t, "", serverSettings.DatabaseDsn)
	assert.True(t, serverSettings.PersistSettings.Restore)
	assert.Nil(t, serverSettings.CryptoPrivKeyPath)
	assert.Nil(t, serverSettings.TrustedSubnet)
	assert.Nil(t, serverSettings.GrpcSettings)
}

func TestServerSettingsAdaptCustomEnv(t *testing.T) {
	testStoreFile := []string{"/foo/bar", ""}
	for _, storeFile := range testStoreFile {
		t.Setenv("ADDRESS", "1.1.1.1:9999")
		t.Setenv("STORE_INTERVAL", "0")
		t.Setenv("STORE_FILE", storeFile)
		t.Setenv("RESTORE", "false")
		t.Setenv("KEY", "123")
		t.Setenv("DATABASE_DSN", "postgre:5444")
		t.Setenv("CRYPTO_KEY", "/home/.ssh/id_rsa")
		t.Setenv("TRUSTED_SUBNET", "192.168.0.0/16")
		t.Setenv("GRPC_ADDRESS", "10.0.0.0:5555")

		testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
		config, err := LoadConfig(*testFlagSet, []string{})
		assert.NoError(t, err)

		serverSettings, err := ServerSettingsAdapt(config)
		assert.NoError(t, err)

		assert.Equal(t, "1.1.1.1", serverSettings.HttpSettings.Host)
		assert.Equal(t, 9999, serverSettings.HttpSettings.Port)
		assert.Equal(t, "123", *serverSettings.HmacKey)
		assert.Equal(t, time.Duration(0), serverSettings.PersistSettings.StoreInterval)
		assert.Equal(t, storeFile, serverSettings.PersistSettings.StoreFile)
		assert.False(t, serverSettings.PersistSettings.Restore)
		assert.Equal(t, "postgre:5444", serverSettings.DatabaseDsn)
		assert.Equal(t, "/home/.ssh/id_rsa", *serverSettings.CryptoPrivKeyPath)
		assert.Equal(t, getIPNet("192.168.0.0/16"), serverSettings.TrustedSubnet)
		assert.Equal(t, "10.0.0.0", serverSettings.GrpcSettings.Host)
		assert.Equal(t, 5555, serverSettings.GrpcSettings.Port)
	}
}

func TestServerSettingsAdaptCustomFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(
		*testFlagSet,
		[]string{
			"-a", "1.1.1.1:9999",
			"-i", "0s",
			"-f", "/tmp/ttt",
			"-r=false",
			"-k", "123",
			"-d", "postgre:5444",
			"-crypto-key", "/home/.ssh/id_rsa",
			"-t", "192.168.0.0/16",
			"-g", "10.0.0.0:5555"})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "1.1.1.1", serverSettings.HttpSettings.Host)
	assert.Equal(t, 9999, serverSettings.HttpSettings.Port)
	assert.Equal(t, "123", *serverSettings.HmacKey)
	assert.Equal(t, time.Duration(0), serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/ttt", serverSettings.PersistSettings.StoreFile)
	assert.False(t, serverSettings.PersistSettings.Restore)
	assert.Equal(t, "postgre:5444", serverSettings.DatabaseDsn)
	assert.Equal(t, "/home/.ssh/id_rsa", *serverSettings.CryptoPrivKeyPath)
	assert.Equal(t, getIPNet("192.168.0.0/16"), serverSettings.TrustedSubnet)
	assert.Equal(t, "10.0.0.0", serverSettings.GrpcSettings.Host)
	assert.Equal(t, 5555, serverSettings.GrpcSettings.Port)
}

func TestServerSettingsAdaptCustomEnvAndFlag(t *testing.T) {
	t.Setenv("ADDRESS", "1.1.1.1:9999")
	t.Setenv("STORE_INTERVAL", "0")
	t.Setenv("STORE_FILE", "/tmp/ttt")
	t.Setenv("RESTORE", "false")
	t.Setenv("KEY", "123")
	t.Setenv("DATABASE_DSN", "postgre:5444")
	t.Setenv("CRYPTO_KEY", "/home/.ssh/id_rsa")
	t.Setenv("TRUSTED_SUBNET", "10.0.0.0/16")
	t.Setenv("GRPC_ADDRESS", "10.0.0.0:5555")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(
		*testFlagSet,
		[]string{
			"-a", "7.7.7.7:7777",
			"-i", "10s",
			"-f", "/tmp/aaa",
			"-r=true",
			"-k", "456",
			"-d", "postgre:5444",
			"-crypto-key", "./id_rsa",
			"-t", "192.168.0.0/16",
			"-g", "10.10.10.10:5555"})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "1.1.1.1", serverSettings.HttpSettings.Host)
	assert.Equal(t, 9999, serverSettings.HttpSettings.Port)
	assert.Equal(t, "123", *serverSettings.HmacKey)
	assert.Equal(t, time.Duration(0), serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/ttt", serverSettings.PersistSettings.StoreFile)
	assert.False(t, serverSettings.PersistSettings.Restore)
	assert.Equal(t, "postgre:5444", serverSettings.DatabaseDsn)
	assert.Equal(t, "/home/.ssh/id_rsa", *serverSettings.CryptoPrivKeyPath)
	assert.Equal(t, getIPNet("10.0.0.0/16"), serverSettings.TrustedSubnet)
	assert.Equal(t, "10.0.0.0", serverSettings.GrpcSettings.Host)
	assert.Equal(t, 5555, serverSettings.GrpcSettings.Port)
}

func TestServerSettingsAdaptCustomEnvAndFlagMix(t *testing.T) {
	t.Setenv("STORE_INTERVAL", "1s")
	t.Setenv("RESTORE", "true")
	t.Setenv("CRYPTO_KEY", "/home/.ssh/id_rsa")
	t.Setenv("TRUSTED_SUBNET", "10.0.0.0/16")
	t.Setenv("GRPC_ADDRESS", "10.0.0.0:5555")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-i", "5m", "-r=false", "-k", "123"})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1", serverSettings.HttpSettings.Host)
	assert.Equal(t, 8080, serverSettings.HttpSettings.Port)
	assert.Equal(t, "123", *serverSettings.HmacKey)
	assert.Equal(t, 1*time.Second, serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/devops-metrics-db.json", serverSettings.PersistSettings.StoreFile)
	assert.True(t, serverSettings.PersistSettings.Restore)
	assert.Equal(t, "", serverSettings.DatabaseDsn)
	assert.Equal(t, "/home/.ssh/id_rsa", *serverSettings.CryptoPrivKeyPath)
	assert.Equal(t, getIPNet("10.0.0.0/16"), serverSettings.TrustedSubnet)
	assert.Equal(t, "10.0.0.0", serverSettings.GrpcSettings.Host)
	assert.Equal(t, 5555, serverSettings.GrpcSettings.Port)
}

func TestServerSettingsAdaptCustomError(t *testing.T) {
	for i, tt := range []struct {
		varName string
		varVal  string
	}{
		{varName: "ADDRESS", varVal: "1.1.1.1"},
		{varName: "ADDRESS", varVal: "1.1.1.1:foobar"},
		{varName: "GRPC_ADDRESS", varVal: "1.1.1.1"},
		{varName: "GRPC_ADDRESS", varVal: "1.1.1.1:foobar"},
		{varName: "TRUSTED_SUBNET", varVal: "abcdef"},
		{varName: "TRUSTED_SUBNET", varVal: "10.0.0.0"},
		{varName: "TRUSTED_SUBNET", varVal: "10.0.0.0/500"},
	} {
		tt := tt
		i := i
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			t.Setenv(tt.varName, tt.varVal)

			testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
			config, err := LoadConfig(*testFlagSet, []string{})
			assert.NoError(t, err)

			_, err = ServerSettingsAdapt(config)
			assert.Error(t, err)
		})
	}
}

func TestServerSettingsCastEnvError(t *testing.T) {
	for i, tt := range []struct {
		envVarName string
		envVarVal  string
	}{
		{envVarName: "STORE_INTERVAL", envVarVal: "foobar"},
		{envVarName: "RESTORE", envVarVal: "foobar"},
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

func TestServerSettingsWithConfigFile(t *testing.T) {
	// Create temp config file
	fCfg, err := os.CreateTemp("", "cfg")
	require.NoError(t, err)

	defer func() {
		fCfg.Close()
		os.Remove(fCfg.Name())
	}()

	cfgAddr := "172.100.1.1:9090"
	cfgStoreInt := 100 * time.Minute
	databaseDsn := "foobar"
	trustedSubnet := "10.0.0.0/16"

	tempCfg := configFile{
		Address:       &cfgAddr,
		StoreInterval: &cfgStoreInt,
		DatabaseDsn:   &databaseDsn,
		TrustedSubnet: &trustedSubnet,
	}
	assert.NoError(t, json.NewEncoder(fCfg).Encode(&tempCfg))

	// Check
	t.Setenv("DATABASE_DSN", "postgre:5444")
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-i", "10s", "-config", fCfg.Name()})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "172.100.1.1", serverSettings.HttpSettings.Host)
	assert.Equal(t, 9090, serverSettings.HttpSettings.Port)
	assert.Nil(t, serverSettings.HmacKey)
	assert.Equal(t, 10*time.Second, serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/devops-metrics-db.json", serverSettings.PersistSettings.StoreFile)
	assert.Equal(t, "postgre:5444", serverSettings.DatabaseDsn)
	assert.True(t, serverSettings.PersistSettings.Restore)
	assert.Nil(t, serverSettings.CryptoPrivKeyPath)
	assert.Equal(t, getIPNet("10.0.0.0/16"), serverSettings.TrustedSubnet)
	assert.Nil(t, serverSettings.GrpcSettings)
}

func TestServerSettingsAllFromConfigFile(t *testing.T) {
	// Create temp config file
	fCfg, err := os.CreateTemp("", "cfg")
	require.NoError(t, err)

	defer func() {
		fCfg.Close()
		os.Remove(fCfg.Name())
	}()

	cfgAddr := "172.100.1.1:9090"
	cfgRestore := true
	cfgStoreInt := 100 * time.Minute
	cfgStoreFile := "/tmp/store"
	cfgDatabaseDsn := "foobar"
	cfgHmacKey := "hmac_key"
	cfgPrivKey := "/tmp/id_rsa"
	cfgTrustedSubnet := "10.0.0.0/16"
	cfgGrpcAddress := "10.0.0.0:5555"

	tempCfg := configFile{
		Address:           &cfgAddr,
		Restore:           &cfgRestore,
		StoreInterval:     &cfgStoreInt,
		StoreFile:         &cfgStoreFile,
		DatabaseDsn:       &cfgDatabaseDsn,
		HmacKey:           &cfgHmacKey,
		CryptoPrivKeyPath: &cfgPrivKey,
		TrustedSubnet:     &cfgTrustedSubnet,
		GrpcAddress:       &cfgGrpcAddress,
	}
	assert.NoError(t, json.NewEncoder(fCfg).Encode(&tempCfg))

	// Check
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{"-config", fCfg.Name()})
	assert.NoError(t, err)

	serverSettings, err := ServerSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "172.100.1.1", serverSettings.HttpSettings.Host)
	assert.Equal(t, 9090, serverSettings.HttpSettings.Port)
	assert.Equal(t, "hmac_key", *serverSettings.HmacKey)
	assert.Equal(t, 100*time.Minute, serverSettings.PersistSettings.StoreInterval)
	assert.Equal(t, "/tmp/store", serverSettings.PersistSettings.StoreFile)
	assert.Equal(t, "foobar", serverSettings.DatabaseDsn)
	assert.True(t, serverSettings.PersistSettings.Restore)
	assert.Equal(t, "/tmp/id_rsa", *serverSettings.CryptoPrivKeyPath)
	assert.Equal(t, getIPNet("10.0.0.0/16"), serverSettings.TrustedSubnet)
	assert.Equal(t, "10.0.0.0", serverSettings.GrpcSettings.Host)
	assert.Equal(t, 5555, serverSettings.GrpcSettings.Port)
}

func getIPNet(cidr string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidr)
	return ipNet
}
