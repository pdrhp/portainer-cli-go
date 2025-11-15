package printer

import (
	"bytes"
	"encoding/json"
	"testing"

	"os"

	"github.com/pdrhp/portainer-go-cli/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintStacks_JSON(t *testing.T) {
	stacks := []types.Stack{
		{
			ID:         1,
			Name:       "test-stack",
			Type:       types.StackTypeDockerCompose,
			EndpointID: 1,
			Status:     1,
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := printStacksJSON(stacks)
	require.NoError(t, err)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result []types.Stack
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test-stack", result[0].Name)
}
