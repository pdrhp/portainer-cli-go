package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pdrhp/portainer-go-cli/internal/client"
	"github.com/pdrhp/portainer-go-cli/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigsCommand() *cobra.Command {
	group := &cobra.Command{Use: "configs", Short: "Manage Docker Swarm configs"}
	var endpointID int
	var file string
	create := &cobra.Command{
		Use: "create NAME", Short: "Create an immutable Docker Swarm config from a file or stdin",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("name must not be empty")
			}
			if endpointID <= 0 {
				return fmt.Errorf("endpoint-id must be positive")
			}
			output := cmd.Flag("output").Value.String()
			if output != "table" && output != "json" && output != "yaml" {
				return fmt.Errorf("unsupported output format: %s", output)
			}
			var reader io.Reader = cmd.InOrStdin()
			if file != "-" {
				f, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("read config file: %w", err)
				}
				defer f.Close()
				reader = f
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("read config file: %w", err)
			}
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cfg.Token == "" {
				return fmt.Errorf("not authenticated. Please run 'portainer-go auth' first")
			}
			serverURL := cmd.Flag("server-url").Value.String()
			if serverURL == "" {
				serverURL = cfg.ServerURL
			}
			if serverURL == "" {
				return fmt.Errorf("server URL not configured. Use --server-url or set it in config")
			}
			cl := client.New(strings.TrimRight(serverURL, "/"))
			cl.SetToken(cfg.Token)
			result, err := cl.CreateConfig(cmd.Context(), endpointID, args[0], data)
			if err != nil {
				return err
			}
			switch output {
			case "json":
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			case "yaml":
				return yaml.NewEncoder(cmd.OutOrStdout()).Encode(result)
			default:
				_, err = fmt.Fprintln(cmd.OutOrStdout(), result.ID)
				return err
			}
		},
	}
	create.Flags().IntVar(&endpointID, "endpoint-id", 0, "Portainer Swarm environment ID")
	create.Flags().StringVar(&file, "file", "", "Config file path, or - for stdin")
	_ = create.MarkFlagRequired("endpoint-id")
	_ = create.MarkFlagRequired("file")
	group.AddCommand(create)
	return group
}
