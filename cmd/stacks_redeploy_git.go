package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pdrhp/portainer-go-cli/internal/client"
	"github.com/pdrhp/portainer-go-cli/internal/config"
	"github.com/pdrhp/portainer-go-cli/internal/wizard"
	"github.com/pdrhp/portainer-go-cli/pkg/types"
	"github.com/spf13/cobra"
)

var (
	redeployGitStackID                 int
	redeployGitEndpointID              int
	redeployGitRepositoryReferenceName string
	redeployGitRepositoryUsername      string
	redeployGitRepositoryPassword      string
	redeployGitEnv                     []string
	redeployGitPrune                   bool
	redeployGitPullImage               bool
	redeployGitStackName               string
)

var stacksRedeployGitCmd = &cobra.Command{
	Use:   "redeploy [stack-id]",
	Short: "Redeploy a stack from its Git repository",
	Long: `Redeploy an existing stack by pulling the latest changes from its Git repository.

Examples:
  # Redeploy with flags
  portainer stacks redeploy 123 --endpoint-id 1 --env KEY1=value1 --prune --pull-image

  # Redeploy with Git authentication
  portainer stacks redeploy 123 --endpoint-id 1 --repository-username user --repository-password pass

  # Interactive redeploy
  portainer stacks redeploy`,
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

		var payload types.StackGitRedeployPayload
		var stackID int
		var endpointID int

		if len(args) > 0 {
			if redeployGitStackID > 0 {
				return fmt.Errorf("stack ID cannot be specified both as argument and flag")
			}
			stackID, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid stack ID: %s", args[0])
			}
		} else if redeployGitStackID > 0 {
			stackID = redeployGitStackID
		} else {
			// Use wizard
			wizardPayload, wizardStackID, wizardEndpointID, err := wizard.RunRedeployGitWizard()
			if err != nil {
				return fmt.Errorf("wizard failed: %w", err)
			}
			payload = *wizardPayload
			stackID = wizardStackID
			endpointID = wizardEndpointID
		}

		if endpointID == 0 {
			if redeployGitEndpointID == 0 {
				return fmt.Errorf("--endpoint-id is required")
			}
			endpointID = redeployGitEndpointID

			payload = buildRedeployPayloadFromFlags()
		}

		cl := client.New(serverURL)
		cl.SetToken(cfg.Token)

		fmt.Printf("Redeploying stack %d from Git repository...\n", stackID)

		stack, err := cl.RedeployStackFromGit(cmd.Context(), stackID, endpointID, payload)
		if err != nil {
			var httpErr *client.HTTPError
			if errors.As(err, &httpErr) {
				switch httpErr.StatusCode {
				case 400:
					return fmt.Errorf("Invalid request: %s", httpErr.Message)
				case 403:
					return fmt.Errorf("Permission denied. You don't have access to this stack")
				case 404:
					return fmt.Errorf("Stack not found")
				case 401:
					return fmt.Errorf("Authentication failed. Please run 'portainer auth' again")
				default:
					return fmt.Errorf("Failed to redeploy stack (HTTP %d): %s", httpErr.StatusCode, httpErr.Message)
				}
			}
			return fmt.Errorf("Failed to redeploy stack: %w", err)
		}

		fmt.Printf("Stack '%s' redeployed successfully\n", stack.Name)
		return nil
	},
}

func buildRedeployPayloadFromFlags() types.StackGitRedeployPayload {
	payload := types.StackGitRedeployPayload{
		RepositoryReferenceName: redeployGitRepositoryReferenceName,
		StackName:               redeployGitStackName,
		Prune:                   redeployGitPrune,
		PullImage:               redeployGitPullImage,
	}

	if redeployGitRepositoryUsername != "" && redeployGitRepositoryPassword != "" {
		payload.RepositoryAuthentication = true
		payload.RepositoryUsername = redeployGitRepositoryUsername
		payload.RepositoryPassword = redeployGitRepositoryPassword
	}

	if len(redeployGitEnv) > 0 {
		payload.Env = make([]types.Pair, 0, len(redeployGitEnv))
		for _, env := range redeployGitEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				payload.Env = append(payload.Env, types.Pair{
					Name:  parts[0],
					Value: parts[1],
				})
			}
		}
	}

	return payload
}

func init() {
	stacksRedeployGitCmd.Flags().IntVar(&redeployGitStackID, "stack-id", 0, "Stack ID to redeploy (alternative to positional argument)")
	stacksRedeployGitCmd.Flags().IntVar(&redeployGitEndpointID, "endpoint-id", 0, "Environment identifier (required)")
	stacksRedeployGitCmd.Flags().StringVar(&redeployGitRepositoryReferenceName, "repository-reference-name", "", "Git reference (branch/tag)")
	stacksRedeployGitCmd.Flags().StringVar(&redeployGitRepositoryUsername, "repository-username", "", "Username for Git repository authentication")
	stacksRedeployGitCmd.Flags().StringVar(&redeployGitRepositoryPassword, "repository-password", "", "Password for Git repository authentication")
	stacksRedeployGitCmd.Flags().StringSliceVar(&redeployGitEnv, "env", []string{}, "Environment variables (format: KEY=value)")
	stacksRedeployGitCmd.Flags().BoolVar(&redeployGitPrune, "prune", false, "Remove services that are no longer referenced")
	stacksRedeployGitCmd.Flags().BoolVar(&redeployGitPullImage, "pull-image", false, "Force pull the latest image")
	stacksRedeployGitCmd.Flags().StringVar(&redeployGitStackName, "stack-name", "", "Stack name (Kubernetes only)")

}
