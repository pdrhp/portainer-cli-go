package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/pdrhp/portainer-go-cli/internal/client"
	"github.com/pdrhp/portainer-go-cli/internal/config"
	"github.com/pdrhp/portainer-go-cli/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	authUsername string
	authPassword string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Portainer server",
	Long: `Authenticate with Portainer server using username/password or interactive wizard.
	
Examples:
  # Authenticate with flags
  portainer auth --username admin --password mypass
  
  # Interactive authentication
  portainer auth
  
  # Use config values
  portainer auth --server-url https://portainer.example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		serverURL := cmd.Flag("server-url").Value.String()
		if serverURL == "" {
			serverURL = cfg.ServerURL
		}
		if serverURL == "" {
			return fmt.Errorf("server URL not provided. Use --server-url flag or set it in config")
		}

		var username, password string

		if authUsername != "" && authPassword != "" {
			username = authUsername
			password = authPassword
		} else if cfg.Username != "" && cfg.Password != "" {
			username = cfg.Username
			password = cfg.Password
		} else {
			fmt.Println("Authentication required")
			fmt.Printf("Server: %s\n\n", serverURL)

			creds, err := wizard.RunAuthWizard()
			if err != nil {
				return fmt.Errorf("authentication wizard failed: %w", err)
			}
			username = creds.Username
			password = creds.Password
		}

		cl := client.New(serverURL)

		fmt.Printf("Authenticating with %s...\n", serverURL)
		token, err := cl.Authenticate(ctx, username, password)
		if err != nil {
			var httpErr *client.HTTPError
			if errors.As(err, &httpErr) {
				switch httpErr.StatusCode {
				case 400, 422:
					return fmt.Errorf("Invalid credentials. Please check username and password")
				case 401:
					return fmt.Errorf("Authentication failed. Invalid username or password")
				case 500:
					return fmt.Errorf("Server error. Please try again later")
				default:
					return fmt.Errorf("Authentication failed (HTTP %d): %s", httpErr.StatusCode, httpErr.Message)
				}
			}
			return fmt.Errorf("Authentication failed: %w", err)
		}

		cfg.Token = token
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Authentication successful! Token saved.")
		return nil
	},
}

func init() {
	authCmd.Flags().StringVar(&authUsername, "username", "", "Username for authentication")
	authCmd.Flags().StringVar(&authPassword, "password", "", "Password for authentication")
}
