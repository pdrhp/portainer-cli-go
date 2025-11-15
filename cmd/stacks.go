package cmd

import (
	"github.com/spf13/cobra"
)

var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "Manage stacks",
	Long:  `Manage Docker stacks in Portainer environment`,
}

func init() {
	stacksCmd.AddCommand(stacksListCmd)
}
