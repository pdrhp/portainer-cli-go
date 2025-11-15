package cmd

import (
	"fmt"
	"strings"

	"github.com/pdrhp/portainer-go-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Set and get configuration values for the Portainer CLI`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Available keys:
- server-url: Portainer server URL
- username: Default username for authentication
- password: Default password for authentication
- api-key: API key for authentication`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		switch key {
		case "server-url", "server_url":
			cfg.ServerURL = value
		case "username":
			cfg.Username = value
		case "password":
			cfg.Password = value
		case "api-key", "api_key":
			cfg.APIKey = value
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Configuration updated: %s = %s\n", key, maskSensitiveValue(key, value))
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value. If no key is specified, shows all values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(args) == 0 {
			fmt.Println("Current configuration:")
			fmt.Printf("Server URL: %s\n", cfg.ServerURL)
			fmt.Printf("Username: %s\n", cfg.Username)
			fmt.Printf("Password: %s\n", maskValue(cfg.Password))
			fmt.Printf("API Key: %s\n", maskValue(cfg.APIKey))
			fmt.Printf("Token: %s\n", maskValue(cfg.Token))
		} else {
			key := args[0]
			switch key {
			case "server-url", "server_url":
				fmt.Println(cfg.ServerURL)
			case "username":
				fmt.Println(cfg.Username)
			case "password":
				fmt.Println(maskValue(cfg.Password))
			case "api-key", "api_key":
				fmt.Println(maskValue(cfg.APIKey))
			case "token":
				fmt.Println(maskValue(cfg.Token))
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
}

func maskSensitiveValue(key, value string) string {
	sensitiveKeys := []string{"password", "token", "api-key", "api_key"}
	for _, sensitiveKey := range sensitiveKeys {
		if key == sensitiveKey {
			return maskValue(value)
		}
	}
	return value
}

func maskValue(value string) string {
	if value == "" {
		return ""
	}
	return strings.Repeat("*", len(value))
}
