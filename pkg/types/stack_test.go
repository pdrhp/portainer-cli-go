package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackType_String(t *testing.T) {
	tests := []struct {
		stackType StackType
		expected  string
	}{
		{StackTypeDockerCompose, "compose"},
		{StackTypeDockerSwarm, "swarm"},
		{StackType(999), "unknown"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.stackType.String())
	}
}
