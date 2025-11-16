package cmd

import (
	"errors"
	"fmt"

	"github.com/pdrhp/portainer-go-cli/internal/client"
	"github.com/pdrhp/portainer-go-cli/internal/config"
	"github.com/pdrhp/portainer-go-cli/internal/printer"
	"github.com/pdrhp/portainer-go-cli/pkg/types"
	"github.com/spf13/cobra"
)

var (
	endpointID int
	swarmID    string
)

var stacksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stacks",
	Long: `List all stacks available in the Portainer environment.
	
Examples:
  # List all stacks
  portainer stacks list
  
  # Filter by endpoint
  portainer stacks list --endpoint-id 1
  
  # Filter by Swarm cluster
  portainer stacks list --swarm-id jpofkc0i9uo9wtx1zesuk649w
  
  # Output in JSON format
  portainer stacks list --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Token == "" {
			return fmt.Errorf("not authenticated. Please run 'portainer auth' first")
		}

		serverURL := cmd.Flag("server-url").Value.String()
		if serverURL == "" {
			serverURL = cfg.ServerURL
		}
		if serverURL == "" {
			return fmt.Errorf("server URL not configured. Use --server-url flag or set it in config")
		}

		var filters *types.StackFilters
		if endpointID > 0 || swarmID != "" {
			filters = &types.StackFilters{
				EndpointID: endpointID,
				SwarmID:    swarmID,
			}
		}

		cl := client.New(serverURL)
		cl.SetToken(cfg.Token)

		stacks, err := cl.ListStacks(cmd.Context(), filters)
		if err != nil {
			var httpErr *client.HTTPError
			if errors.As(err, &httpErr) {
				switch httpErr.StatusCode {
				case 401, 403:
					return fmt.Errorf("Authentication failed. Please run 'portainer auth' again")
				case 404:
					return fmt.Errorf("Endpoint not found")
				default:
					return fmt.Errorf("Failed to list stacks (HTTP %d): %s", httpErr.StatusCode, httpErr.Message)
				}
			}
			return fmt.Errorf("Failed to list stacks: %w", err)
		}

		if len(stacks) == 0 {
			fmt.Println("No stacks found.")
			return nil
		}

		outputFormat := cmd.Flag("output").Value.String()

		return printer.PrintStacks(stacks, outputFormat)
	},
}

func init() {
	stacksListCmd.Flags().IntVar(&endpointID, "endpoint-id", 0, "Filter stacks by endpoint ID")
	stacksListCmd.Flags().StringVar(&swarmID, "swarm-id", "", "Filter stacks by Swarm cluster ID")
}
