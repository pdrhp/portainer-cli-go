package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "portainer-go",
	Short: "Portainer Go CLI for CI/CD automation",
	Long:  `A command-line interface for managing Portainer stacks and resources`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("server-url", "", "Portainer server URL")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table|json|yaml)")

	viper.BindPFlag("server_url", rootCmd.PersistentFlags().Lookup("server-url"))

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PORTAINER")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(stacksCmd)
}
