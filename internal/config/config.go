package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigDir  = ".portainer-cli"
	ConfigFile = "config.yaml"
)

type Config struct {
	ServerURL string `mapstructure:"server_url" yaml:"server_url"`
	Username  string `mapstructure:"username" yaml:"username"`
	Password  string `mapstructure:"password" yaml:"password"`
	Token     string `mapstructure:"token" yaml:"token"`
	APIKey    string `mapstructure:"api_key" yaml:"api_key"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("$HOME/" + ConfigDir)
	viper.AddConfigPath("/etc/portainer-cli")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}

	configDir := filepath.Join(homeDir, ConfigDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	configPath := filepath.Join(configDir, ConfigFile)

	viper.Set("server_url", cfg.ServerURL)
	viper.Set("username", cfg.Username)
	viper.Set("password", cfg.Password)
	viper.Set("token", cfg.Token)
	viper.Set("api_key", cfg.APIKey)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}

	return filepath.Join(homeDir, ConfigDir, ConfigFile), nil
}
