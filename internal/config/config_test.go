package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load_Save(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "portainer-cli-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cfg := &Config{
		ServerURL: "https://portainer.example.com",
		Username:  "testuser",
		Password:  "testpass",
		APIKey:    "testkey",
		Token:     "testtoken",
	}

	err = Save(cfg)
	require.NoError(t, err)

	loadedCfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, cfg.ServerURL, loadedCfg.ServerURL)
	assert.Equal(t, cfg.Username, loadedCfg.Username)
	assert.Equal(t, cfg.Password, loadedCfg.Password)
	assert.Equal(t, cfg.APIKey, loadedCfg.APIKey)
	assert.Equal(t, cfg.Token, loadedCfg.Token)
}

func TestConfig_GetConfigPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "portainer-cli-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	path, err := GetConfigPath()
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, ConfigDir, ConfigFile)
	assert.Equal(t, expectedPath, path)
}
