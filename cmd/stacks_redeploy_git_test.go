package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	payload, err := buildRedeployPayloadFromFlags()
	require.NoError(t, err)

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

	payload, err := buildRedeployPayloadFromFlags()
	require.NoError(t, err)

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

	payload, err := buildRedeployPayloadFromFlags()
	require.NoError(t, err)

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

func TestBuildRedeployPayloadFromFlags_EnvWithSpaces(t *testing.T) {
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = []string{"GREETING=  hello world  "}
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	payload, err := buildRedeployPayloadFromFlags()
	require.NoError(t, err)
	require.Len(t, payload.Env, 1)

	assert.Equal(t, "GREETING", payload.Env[0].Name)
	assert.Equal(t, "  hello world  ", payload.Env[0].Value)
}

func TestBuildRedeployPayloadFromFlags_InvalidEnv(t *testing.T) {
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = []string{"INVALID"}
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	_, err := buildRedeployPayloadFromFlags()
	require.Error(t, err)
	assert.ErrorContains(t, err, "missing '='")
}

func TestBuildRedeployPayloadFromFlags_RepositoryCredentialsMustBeTogether(t *testing.T) {
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = "user"
	redeployGitRepositoryPassword = ""
	redeployGitEnv = []string{}
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	_, err := buildRedeployPayloadFromFlags()
	require.Error(t, err)
	assert.EqualError(t, err, "--repository-username and --repository-password must be provided together")
}

func TestHasNonInteractiveRedeployInput_FalseWhenNoFlags(t *testing.T) {
	redeployGitEndpointID = 0
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = nil
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	assert.False(t, hasNonInteractiveRedeployInput())
}

func TestHasNonInteractiveRedeployInput_TrueWhenFlagsProvided(t *testing.T) {
	redeployGitEndpointID = 1
	redeployGitRepositoryReferenceName = ""
	redeployGitRepositoryUsername = ""
	redeployGitRepositoryPassword = ""
	redeployGitEnv = nil
	redeployGitPrune = false
	redeployGitPullImage = false
	redeployGitStackName = ""

	assert.True(t, hasNonInteractiveRedeployInput())
}
