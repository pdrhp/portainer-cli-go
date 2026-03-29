package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pdrhp/portainer-go-cli/internal/client"
	"github.com/pdrhp/portainer-go-cli/internal/config"
	"github.com/pdrhp/portainer-go-cli/internal/envvars"
	"github.com/pdrhp/portainer-go-cli/internal/wizard"
	"github.com/pdrhp/portainer-go-cli/pkg/types"
	"github.com/spf13/cobra"
)

var (
	createSwarmGitName                     string
	createSwarmGitRepositoryURL            string
	createSwarmGitSwarmID                  string
	createSwarmGitEndpointID               int
	createSwarmGitComposeFile              string
	createSwarmGitRepositoryReferenceName  string
	createSwarmGitTLSSkipVerify            bool
	createSwarmGitRepositoryUsername       string
	createSwarmGitRepositoryPassword       string
	createSwarmGitEnv                      []string
	createSwarmGitAdditionalFiles          []string
	createSwarmGitAutoUpdateInterval       string
	createSwarmGitAutoUpdateWebhook        string
	createSwarmGitAutoUpdateForcePullImage bool
	createSwarmGitAutoUpdateForceUpdate    bool
)

var stacksCreateSwarmGitCmd = &cobra.Command{
	Use:   "create-swarm-git",
	Short: "Create a new swarm stack from a git repository",
	Long: `Create a new Docker Swarm stack by pulling the compose file from a Git repository.

Examples:
  # Create stack with all flags
  portainer stacks create-swarm-git --name myStack --repository-url https://github.com/user/repo --swarm-id jpofkc0i9uo9wtx1zesuk649w --endpoint-id 1 --compose-file docker-compose.yml --repository-reference-name refs/heads/main --env KEY1=value1 --env KEY2=value2

  # Create stack with Git authentication
  portainer stacks create-swarm-git --name myStack --repository-url https://github.com/user/repo --swarm-id jpofkc0i9uo9wtx1zesuk649w --endpoint-id 1 --repository-username user --repository-password pass

  # Create stack with auto-update (GitOps)
  portainer stacks create-swarm-git --name myStack --repository-url https://github.com/user/repo --swarm-id jpofkc0i9uo9wtx1zesuk649w --endpoint-id 1 --auto-update-interval 1h --auto-update-webhook abc123

  # Interactive creation
  portainer stacks create-swarm-git`,
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

		var payload types.StackCreateSwarmGitPayload
		var endpointID int

		useWizard := createSwarmGitName == "" && createSwarmGitRepositoryURL == "" && createSwarmGitSwarmID == "" && createSwarmGitEndpointID == 0

		if useWizard {
			wizardPayload, endpointIDFromWizard, err := wizard.RunCreateSwarmGitWizard()
			if err != nil {
				return fmt.Errorf("wizard failed: %w", err)
			}
			payload = *wizardPayload
			endpointID = endpointIDFromWizard
		} else {
			if createSwarmGitName == "" {
				return fmt.Errorf("flag --name is required when not using wizard")
			}
			if createSwarmGitRepositoryURL == "" {
				return fmt.Errorf("flag --repository-url is required when not using wizard")
			}
			if createSwarmGitSwarmID == "" {
				return fmt.Errorf("flag --swarm-id is required when not using wizard")
			}
			if createSwarmGitEndpointID == 0 {
				return fmt.Errorf("flag --endpoint-id is required when not using wizard")
			}

			payload, err = buildPayloadFromFlags()
			if err != nil {
				return fmt.Errorf("invalid input flags: %w", err)
			}
			endpointID = createSwarmGitEndpointID
		}

		cl := client.New(serverURL)
		cl.SetToken(cfg.Token)

		fmt.Printf("Creating swarm stack '%s' from git repository...\n", payload.Name)

		stack, err := cl.CreateSwarmStackFromGit(cmd.Context(), endpointID, payload)
		if err != nil {
			var httpErr *client.HTTPError
			if errors.As(err, &httpErr) {
				switch httpErr.StatusCode {
				case 400:
					return fmt.Errorf("Invalid request: %s", httpErr.Message)
				case 409:
					return fmt.Errorf("Stack name or webhook ID already exists")
				case 401, 403:
					return fmt.Errorf("Authentication failed. Please run 'portainer auth' again")
				default:
					return fmt.Errorf("Failed to create stack (HTTP %d): %s", httpErr.StatusCode, httpErr.Message)
				}
			}
			return fmt.Errorf("Failed to create stack: %w", err)
		}

		fmt.Printf("Stack '%s' created successfully with ID: %d\n", stack.Name, stack.ID)
		return nil
	},
}

func buildPayloadFromFlags() (types.StackCreateSwarmGitPayload, error) {
	payload := types.StackCreateSwarmGitPayload{
		Name:                    createSwarmGitName,
		RepositoryURL:           createSwarmGitRepositoryURL,
		SwarmID:                 createSwarmGitSwarmID,
		ComposeFile:             createSwarmGitComposeFile,
		RepositoryReferenceName: createSwarmGitRepositoryReferenceName,
		TLSSkipVerify:           createSwarmGitTLSSkipVerify,
		AdditionalFiles:         createSwarmGitAdditionalFiles,
		FromAppTemplate:         false,
	}

	if payload.ComposeFile == "" {
		payload.ComposeFile = "docker-compose.yml"
	}
	if payload.RepositoryReferenceName == "" {
		payload.RepositoryReferenceName = "refs/heads/master"
	}

	hasRepoUser := createSwarmGitRepositoryUsername != ""
	hasRepoPass := createSwarmGitRepositoryPassword != ""
	if hasRepoUser != hasRepoPass {
		return types.StackCreateSwarmGitPayload{}, fmt.Errorf("--repository-username and --repository-password must be provided together")
	}

	if hasRepoUser {
		payload.RepositoryAuthentication = true
		payload.RepositoryUsername = createSwarmGitRepositoryUsername
		payload.RepositoryPassword = createSwarmGitRepositoryPassword
	}

	if len(createSwarmGitEnv) > 0 {
		env, err := envvars.Parse(strings.Join(createSwarmGitEnv, "\n"))
		if err != nil {
			return types.StackCreateSwarmGitPayload{}, err
		}
		payload.Env = env
	}

	if createSwarmGitAutoUpdateInterval != "" || createSwarmGitAutoUpdateWebhook != "" {
		payload.AutoUpdate = &types.AutoUpdateSettings{
			Interval:       createSwarmGitAutoUpdateInterval,
			Webhook:        createSwarmGitAutoUpdateWebhook,
			ForcePullImage: createSwarmGitAutoUpdateForcePullImage,
			ForceUpdate:    createSwarmGitAutoUpdateForceUpdate,
		}
	}

	return payload, nil
}

func init() {
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitName, "name", "", "Name of the stack (required)")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitRepositoryURL, "repository-url", "", "URL of the Git repository (required)")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitSwarmID, "swarm-id", "", "Swarm cluster identifier (required)")
	stacksCreateSwarmGitCmd.Flags().IntVar(&createSwarmGitEndpointID, "endpoint-id", 0, "Identifier of the environment (required)")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitComposeFile, "compose-file", "docker-compose.yml", "Path to the compose file in the repository")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitRepositoryReferenceName, "repository-reference-name", "refs/heads/master", "Git reference (branch/tag)")
	stacksCreateSwarmGitCmd.Flags().BoolVar(&createSwarmGitTLSSkipVerify, "tlsskip-verify", false, "Skip TLS verification for Git repository")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitRepositoryUsername, "repository-username", "", "Username for Git repository authentication")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitRepositoryPassword, "repository-password", "", "Password for Git repository authentication")
	stacksCreateSwarmGitCmd.Flags().StringSliceVar(&createSwarmGitEnv, "env", []string{}, "Environment variables (format: KEY=value)")
	stacksCreateSwarmGitCmd.Flags().StringSliceVar(&createSwarmGitAdditionalFiles, "additional-files", []string{}, "Additional compose files")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitAutoUpdateInterval, "auto-update-interval", "", "Auto-update interval (e.g., 1h, 30m)")
	stacksCreateSwarmGitCmd.Flags().StringVar(&createSwarmGitAutoUpdateWebhook, "auto-update-webhook", "", "Webhook ID for auto-update")
	stacksCreateSwarmGitCmd.Flags().BoolVar(&createSwarmGitAutoUpdateForcePullImage, "auto-update-force-pull-image", false, "Force pull latest image on auto-update")
	stacksCreateSwarmGitCmd.Flags().BoolVar(&createSwarmGitAutoUpdateForceUpdate, "auto-update-force-update", false, "Force update even without repository changes")
}
