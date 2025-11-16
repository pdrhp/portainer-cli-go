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

func TestStackStatus_String(t *testing.T) {
	tests := []struct {
		stack    Stack
		expected string
	}{
		{Stack{Status: 1}, "running"},
		{Stack{Status: 2}, "stopped"},
		{Stack{Status: 3}, "failed"},
		{Stack{Status: 999}, "unknown"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.stack.StatusString())
	}
}

func TestPair_Creation(t *testing.T) {
	pair := Pair{
		Name:  "TEST_VAR",
		Value: "test_value",
	}

	assert.Equal(t, "TEST_VAR", pair.Name)
	assert.Equal(t, "test_value", pair.Value)
}

func TestAutoUpdateSettings_Creation(t *testing.T) {
	settings := AutoUpdateSettings{
		Interval:       "1h",
		Webhook:        "webhook123",
		ForcePullImage: true,
		ForceUpdate:    false,
	}

	assert.Equal(t, "1h", settings.Interval)
	assert.Equal(t, "webhook123", settings.Webhook)
	assert.True(t, settings.ForcePullImage)
	assert.False(t, settings.ForceUpdate)
}

func TestStackCreateSwarmGitPayload_Creation(t *testing.T) {
	payload := StackCreateSwarmGitPayload{
		Name:                     "test-stack",
		RepositoryURL:            "https://github.com/user/repo",
		SwarmID:                  "jpofkc0i9uo9wtx1zesuk649w",
		ComposeFile:              "docker-compose.yml",
		RepositoryReferenceName:  "refs/heads/master",
		RepositoryAuthentication: true,
		RepositoryUsername:       "user",
		RepositoryPassword:       "pass",
		TLSSkipVerify:            false,
		FromAppTemplate:          false,
	}

	assert.Equal(t, "test-stack", payload.Name)
	assert.Equal(t, "https://github.com/user/repo", payload.RepositoryURL)
	assert.Equal(t, "jpofkc0i9uo9wtx1zesuk649w", payload.SwarmID)
	assert.True(t, payload.RepositoryAuthentication)
	assert.Equal(t, "user", payload.RepositoryUsername)
	assert.Equal(t, "pass", payload.RepositoryPassword)
}

func TestStackCreateSwarmGitPayload_WithEnvAndFiles(t *testing.T) {
	payload := StackCreateSwarmGitPayload{
		Name:          "test-stack",
		RepositoryURL: "https://github.com/user/repo",
		SwarmID:       "jpofkc0i9uo9wtx1zesuk649w",
		Env: []Pair{
			{Name: "KEY1", Value: "value1"},
			{Name: "KEY2", Value: "value2"},
		},
		AdditionalFiles: []string{"file1.yml", "file2.yml"},
	}

	assert.Len(t, payload.Env, 2)
	assert.Equal(t, "KEY1", payload.Env[0].Name)
	assert.Equal(t, "value1", payload.Env[0].Value)

	assert.Len(t, payload.AdditionalFiles, 2)
	assert.Equal(t, "file1.yml", payload.AdditionalFiles[0])
	assert.Equal(t, "file2.yml", payload.AdditionalFiles[1])
}

func TestStackCreateSwarmGitPayload_WithAutoUpdate(t *testing.T) {
	payload := StackCreateSwarmGitPayload{
		Name:          "test-stack",
		RepositoryURL: "https://github.com/user/repo",
		SwarmID:       "jpofkc0i9uo9wtx1zesuk649w",
		AutoUpdate: &AutoUpdateSettings{
			Interval:       "2h",
			Webhook:        "webhook456",
			ForcePullImage: false,
			ForceUpdate:    true,
		},
	}

	assert.NotNil(t, payload.AutoUpdate)
	assert.Equal(t, "2h", payload.AutoUpdate.Interval)
	assert.Equal(t, "webhook456", payload.AutoUpdate.Webhook)
	assert.False(t, payload.AutoUpdate.ForcePullImage)
	assert.True(t, payload.AutoUpdate.ForceUpdate)
}

func TestStackGitRedeployPayload_Creation(t *testing.T) {
	payload := StackGitRedeployPayload{
		RepositoryReferenceName:  "refs/heads/main",
		RepositoryAuthentication: true,
		RepositoryUsername:       "user",
		RepositoryPassword:       "pass",
		Env: []Pair{
			{Name: "ENV_VAR", Value: "value"},
		},
		Prune:     true,
		PullImage: false,
		StackName: "test-stack",
	}

	assert.Equal(t, "refs/heads/main", payload.RepositoryReferenceName)
	assert.True(t, payload.RepositoryAuthentication)
	assert.Equal(t, "user", payload.RepositoryUsername)
	assert.Equal(t, "pass", payload.RepositoryPassword)
	assert.Len(t, payload.Env, 1)
	assert.Equal(t, "ENV_VAR", payload.Env[0].Name)
	assert.Equal(t, "value", payload.Env[0].Value)
	assert.True(t, payload.Prune)
	assert.False(t, payload.PullImage)
	assert.Equal(t, "test-stack", payload.StackName)
}
