package wizard

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pdrhp/portainer-go-cli/pkg/types"
)

type CreateSwarmGitData struct {
	Name                     string
	RepositoryURL            string
	SwarmID                  string
	EndpointID               string
	ComposeFile              string
	RepositoryReferenceName  string
	TLSSkipVerify            bool
	UseAuth                  bool
	RepositoryUsername       string
	RepositoryPassword       string
	EnvVars                  string
	AdditionalFiles          string
	UseAutoUpdate            bool
	AutoUpdateInterval       string
	AutoUpdateWebhook        string
	AutoUpdateForcePullImage bool
	AutoUpdateForceUpdate    bool
}

func RunCreateSwarmGitWizard() (*types.StackCreateSwarmGitPayload, int, error) {
	var data CreateSwarmGitData

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Stack Name").
				Value(&data.Name).
				Validate(func(s string) error {
					if len(s) < 1 {
						return errors.New("stack name cannot be empty")
					}
					if len(s) < 3 {
						return errors.New("stack name must be at least 3 characters")
					}
					return nil
				}),

			huh.NewInput().
				Title("Git Repository URL").
				Value(&data.RepositoryURL).
				Validate(func(s string) error {
					if len(s) < 1 {
						return errors.New("repository URL cannot be empty")
					}
					if !strings.HasPrefix(s, "http") {
						return errors.New("repository URL must start with http:// or https://")
					}
					return nil
				}),

			huh.NewInput().
				Title("Swarm Cluster ID").
				Value(&data.SwarmID).
				Validate(func(s string) error {
					if len(s) < 1 {
						return errors.New("swarm ID cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Endpoint ID").
				Value(&data.EndpointID).
				Validate(func(s string) error {
					if len(s) < 1 {
						return errors.New("endpoint ID cannot be empty")
					}
					if _, err := strconv.Atoi(s); err != nil {
						return errors.New("endpoint ID must be a number")
					}
					return nil
				}),

			huh.NewInput().
				Title("Compose File Path").
				Value(&data.ComposeFile).
				Placeholder("docker-compose.yml"),

			huh.NewInput().
				Title("Git Reference").
				Value(&data.RepositoryReferenceName).
				Placeholder("refs/heads/master"),

			huh.NewConfirm().
				Title("Skip TLS Verification?").
				Value(&data.TLSSkipVerify),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Use Git Authentication?").
				Value(&data.UseAuth),
		).WithHideFunc(func() bool { return false }),

		huh.NewGroup(
			huh.NewInput().
				Title("Git Username").
				Value(&data.RepositoryUsername),

			huh.NewInput().
				Title("Git Password").
				EchoMode(huh.EchoModePassword).
				Value(&data.RepositoryPassword),
		).WithHideFunc(func() bool { return !data.UseAuth }),

		huh.NewGroup(
			huh.NewText().
				Title("Environment Variables").
				Description("Enter environment variables in KEY=value format, one per line").
				Value(&data.EnvVars),

			huh.NewText().
				Title("Additional Files").
				Description("Enter additional compose files, one per line").
				Value(&data.AdditionalFiles),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Auto-Update (GitOps)?").
				Value(&data.UseAutoUpdate),
		).WithHideFunc(func() bool { return false }),

		huh.NewGroup(
			huh.NewInput().
				Title("Auto-Update Interval").
				Description("e.g., 1h, 30m, 1h30m").
				Value(&data.AutoUpdateInterval),

			huh.NewInput().
				Title("Webhook ID").
				Description("UUID for webhook trigger").
				Value(&data.AutoUpdateWebhook),

			huh.NewConfirm().
				Title("Force Pull Image on Update?").
				Value(&data.AutoUpdateForcePullImage),

			huh.NewConfirm().
				Title("Force Update (ignore repo changes)?").
				Value(&data.AutoUpdateForceUpdate),
		).WithHideFunc(func() bool { return !data.UseAutoUpdate }),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return nil, 0, fmt.Errorf("wizard cancelled: %w", err)
	}

	endpointID, _ := strconv.Atoi(data.EndpointID)

	payload := &types.StackCreateSwarmGitPayload{
		Name:                    data.Name,
		RepositoryURL:           data.RepositoryURL,
		SwarmID:                 data.SwarmID,
		ComposeFile:             data.ComposeFile,
		RepositoryReferenceName: data.RepositoryReferenceName,
		TLSSkipVerify:           data.TLSSkipVerify,
		AdditionalFiles:         []string{},
		FromAppTemplate:         false,
	}

	if payload.ComposeFile == "" {
		payload.ComposeFile = "docker-compose.yml"
	}
	if payload.RepositoryReferenceName == "" {
		payload.RepositoryReferenceName = "refs/heads/master"
	}

	if data.UseAuth {
		payload.RepositoryAuthentication = true
		payload.RepositoryUsername = data.RepositoryUsername
		payload.RepositoryPassword = data.RepositoryPassword
	}

	if data.EnvVars != "" {
		lines := strings.Split(strings.TrimSpace(data.EnvVars), "\n")
		payload.Env = make([]types.Pair, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					payload.Env = append(payload.Env, types.Pair{
						Name:  strings.TrimSpace(parts[0]),
						Value: strings.TrimSpace(parts[1]),
					})
				}
			}
		}
	}

	if data.AdditionalFiles != "" {
		lines := strings.Split(strings.TrimSpace(data.AdditionalFiles), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				payload.AdditionalFiles = append(payload.AdditionalFiles, line)
			}
		}
	}

	if data.UseAutoUpdate {
		payload.AutoUpdate = &types.AutoUpdateSettings{
			Interval:       data.AutoUpdateInterval,
			Webhook:        data.AutoUpdateWebhook,
			ForcePullImage: data.AutoUpdateForcePullImage,
			ForceUpdate:    data.AutoUpdateForceUpdate,
		}
	}

	return payload, endpointID, nil
}
