package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRedeployPayloadFromFlags_OnlySpecifiedFields(t *testing.T) {
	// Reset all flags first
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = []string{}
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	// Set only some fields
	redeployGitRepositoryReferenceName = "refs/heads/develop"
	redeployGitRepositoryUsername = "testuser"
	redeployGitRepositoryPassword = "testpass"
	redeployGitEnv = []string{"KEY1=value1"}
	redeployGitPrune = true

	payload := buildRedeployPayloadFromFlags()

	assert.Equal(t, "refs/heads/develop", payload.RepositoryReferenceName)
	assert.True(t, payload.RepositoryAuthentication)
	assert.Equal(t, "testuser", payload.RepositoryUsername)
	assert.Equal(t, "testpass", payload.RepositoryPassword)
	assert.Len(t, payload.Env, 1)
	assert.Equal(t, "KEY1", payload.Env[0].Name)
	assert.Equal(t, "value1", payload.Env[0].Value)
	assert.True(t, payload.Prune) // Now includes the boolean values
	assert.Equal(t, "", payload.StackName)
}

func TestBuildRedeployPayloadFromFlags_WithFlagChanges(t *testing.T) {
	// Reset all flags first
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = []string{}
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	payload := buildRedeployPayloadFromFlags()

	// All boolean flags are included with their default values
	assert.False(t, payload.Prune)
	assert.False(t, payload.PullImage)
	assert.Equal(t, "", payload.RepositoryReferenceName)
	assert.False(t, payload.RepositoryAuthentication)
	assert.Equal(t, "", payload.StackName)
}

func TestBuildRedeployPayloadFromFlags_EmptyFields(t *testing.T) {
	// Reset all flags
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = []string{}
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	payload := buildRedeployPayloadFromFlags()

	// All fields should be empty/false when nothing is set
	assert.Equal(t, "", payload.RepositoryReferenceName)
	assert.False(t, payload.RepositoryAuthentication)
	assert.Equal(t, "", payload.RepositoryUsername)
	assert.Equal(t, "", payload.RepositoryPassword)
	assert.Nil(t, payload.Env)
	assert.Equal(t, "", payload.StackName)
	assert.False(t, payload.Prune)
	assert.False(t, payload.PullImage)
}
