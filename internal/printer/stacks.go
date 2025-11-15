package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v3"

	"github.com/pdrhp/portainer-go-cli/pkg/types"
)

func PrintStacks(stacks []types.Stack, format string) error {
	switch format {
	case "json":
		return printStacksJSON(stacks)
	case "yaml":
		return printStacksYAML(stacks)
	default:
		return printStacksTable(stacks)
	}
}

func printStacksTable(stacks []types.Stack) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tNAME\tTYPE\tSTATUS\tENDPOINT\tSWARM ID")
	fmt.Fprintln(w, "--\t----\t----\t------\t--------\t--------")

	for _, stack := range stacks {
		swarmID := stack.SwarmID
		if len(swarmID) > 12 {
			swarmID = swarmID[:12] + "..."
		}
		if swarmID == "" {
			swarmID = "-"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%s\n",
			stack.ID,
			stack.Name,
			stack.Type.String(),
			stack.StatusString(),
			stack.EndpointID,
			swarmID,
		)
	}

	return w.Flush()
}

func printStacksJSON(stacks []types.Stack) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stacks)
}

func printStacksYAML(stacks []types.Stack) error {
	return yaml.NewEncoder(os.Stdout).Encode(stacks)
}
