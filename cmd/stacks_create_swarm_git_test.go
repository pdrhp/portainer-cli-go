package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStacksCreateSwarmGitCmd_Flags(t *testing.T) {
	cmd := stacksCreateSwarmGitCmd

	// Test required flags exist
	nameFlag := cmd.Flag("name")
	assert.NotNil(t, nameFlag)
	assert.Equal(t, "", nameFlag.DefValue)

	repositoryURLFlag := cmd.Flag("repository-url")
	assert.NotNil(t, repositoryURLFlag)
	assert.Equal(t, "", repositoryURLFlag.DefValue)

	swarmIDFlag := cmd.Flag("swarm-id")
	assert.NotNil(t, swarmIDFlag)
	assert.Equal(t, "", swarmIDFlag.DefValue)

	endpointIDFlag := cmd.Flag("endpoint-id")
	assert.NotNil(t, endpointIDFlag)
	assert.Equal(t, "0", endpointIDFlag.DefValue)

	// Test optional flags
	composeFileFlag := cmd.Flag("compose-file")
	assert.NotNil(t, composeFileFlag)
	assert.Equal(t, "docker-compose.yml", composeFileFlag.DefValue)

	repositoryReferenceNameFlag := cmd.Flag("repository-reference-name")
	assert.NotNil(t, repositoryReferenceNameFlag)
	assert.Equal(t, "refs/heads/master", repositoryReferenceNameFlag.DefValue)

	tlsskipVerifyFlag := cmd.Flag("tlsskip-verify")
	assert.NotNil(t, tlsskipVerifyFlag)
	assert.Equal(t, "false", tlsskipVerifyFlag.DefValue)

	repositoryUsernameFlag := cmd.Flag("repository-username")
	assert.NotNil(t, repositoryUsernameFlag)
	assert.Equal(t, "", repositoryUsernameFlag.DefValue)

	repositoryPasswordFlag := cmd.Flag("repository-password")
	assert.NotNil(t, repositoryPasswordFlag)
	assert.Equal(t, "", repositoryPasswordFlag.DefValue)

	envFlag := cmd.Flag("env")
	assert.NotNil(t, envFlag)

	additionalFilesFlag := cmd.Flag("additional-files")
	assert.NotNil(t, additionalFilesFlag)

	autoUpdateIntervalFlag := cmd.Flag("auto-update-interval")
	assert.NotNil(t, autoUpdateIntervalFlag)
	assert.Equal(t, "", autoUpdateIntervalFlag.DefValue)

	autoUpdateWebhookFlag := cmd.Flag("auto-update-webhook")
	assert.NotNil(t, autoUpdateWebhookFlag)
	assert.Equal(t, "", autoUpdateWebhookFlag.DefValue)

	autoUpdateForcePullImageFlag := cmd.Flag("auto-update-force-pull-image")
	assert.NotNil(t, autoUpdateForcePullImageFlag)
	assert.Equal(t, "false", autoUpdateForcePullImageFlag.DefValue)

	autoUpdateForceUpdateFlag := cmd.Flag("auto-update-force-update")
	assert.NotNil(t, autoUpdateForceUpdateFlag)
	assert.Equal(t, "false", autoUpdateForceUpdateFlag.DefValue)
}

func TestBuildPayloadFromFlags_BasicFields(t *testing.T) {
	// Reset all flags
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitComposeFile = ""
	createSwarmGitRepositoryReferenceName = ""
	createSwarmGitTLSSkipVerify = false
	createSwarmGitRepositoryUsername = ""
	createSwarmGitRepositoryPassword = ""
	createSwarmGitEnv = []string{}
	createSwarmGitAdditionalFiles = []string{}
	createSwarmGitAutoUpdateInterval = ""
	createSwarmGitAutoUpdateWebhook = ""
	createSwarmGitAutoUpdateForcePullImage = false
	createSwarmGitAutoUpdateForceUpdate = false

	payload := buildPayloadFromFlags()

	// Test basic fields
	assert.Equal(t, "test-stack", payload.Name)
	assert.Equal(t, "https://github.com/user/repo", payload.RepositoryURL)
	assert.Equal(t, "jpofkc0i9uo9wtx1zesuk649w", payload.SwarmID)

	// Test defaults are applied
	assert.Equal(t, "docker-compose.yml", payload.ComposeFile)
	assert.Equal(t, "refs/heads/master", payload.RepositoryReferenceName)
	assert.False(t, payload.TLSSkipVerify)
	assert.False(t, payload.FromAppTemplate)
}

func TestBuildPayloadFromFlags_WithAuth(t *testing.T) {
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitRepositoryUsername = "testuser"
	createSwarmGitRepositoryPassword = "testpass"
	createSwarmGitEnv = []string{}
	createSwarmGitAdditionalFiles = []string{}
	createSwarmGitAutoUpdateInterval = ""
	createSwarmGitAutoUpdateWebhook = ""

	payload := buildPayloadFromFlags()

	assert.True(t, payload.RepositoryAuthentication)
	assert.Equal(t, "testuser", payload.RepositoryUsername)
	assert.Equal(t, "testpass", payload.RepositoryPassword)
}

func TestBuildPayloadFromFlags_WithEnvVars(t *testing.T) {
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitRepositoryUsername = ""
	createSwarmGitRepositoryPassword = ""
	createSwarmGitEnv = []string{"KEY1=value1", "KEY2=value2"}
	createSwarmGitAdditionalFiles = []string{}
	createSwarmGitAutoUpdateInterval = ""
	createSwarmGitAutoUpdateWebhook = ""

	payload := buildPayloadFromFlags()

	assert.Len(t, payload.Env, 2)
	assert.Equal(t, "KEY1", payload.Env[0].Name)
	assert.Equal(t, "value1", payload.Env[0].Value)
	assert.Equal(t, "KEY2", payload.Env[1].Name)
	assert.Equal(t, "value2", payload.Env[1].Value)
}

func TestBuildPayloadFromFlags_WithAdditionalFiles(t *testing.T) {
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitRepositoryUsername = ""
	createSwarmGitRepositoryPassword = ""
	createSwarmGitEnv = []string{}
	createSwarmGitAdditionalFiles = []string{"file1.yml", "file2.yml"}
	createSwarmGitAutoUpdateInterval = ""
	createSwarmGitAutoUpdateWebhook = ""

	payload := buildPayloadFromFlags()

	assert.Len(t, payload.AdditionalFiles, 2)
	assert.Equal(t, "file1.yml", payload.AdditionalFiles[0])
	assert.Equal(t, "file2.yml", payload.AdditionalFiles[1])
}

func TestBuildPayloadFromFlags_WithAutoUpdate(t *testing.T) {
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitRepositoryUsername = ""
	createSwarmGitRepositoryPassword = ""
	createSwarmGitEnv = []string{}
	createSwarmGitAdditionalFiles = []string{}
	createSwarmGitAutoUpdateInterval = "1h"
	createSwarmGitAutoUpdateWebhook = "webhook123"
	createSwarmGitAutoUpdateForcePullImage = true
	createSwarmGitAutoUpdateForceUpdate = false

	payload := buildPayloadFromFlags()

	assert.NotNil(t, payload.AutoUpdate)
	assert.Equal(t, "1h", payload.AutoUpdate.Interval)
	assert.Equal(t, "webhook123", payload.AutoUpdate.Webhook)
	assert.True(t, payload.AutoUpdate.ForcePullImage)
	assert.False(t, payload.AutoUpdate.ForceUpdate)
}

func TestBuildPayloadFromFlags_WithoutAutoUpdate(t *testing.T) {
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitRepositoryUsername = ""
	createSwarmGitRepositoryPassword = ""
	createSwarmGitEnv = []string{}
	createSwarmGitAdditionalFiles = []string{}
	createSwarmGitAutoUpdateInterval = ""
	createSwarmGitAutoUpdateWebhook = ""

	payload := buildPayloadFromFlags()

	assert.Nil(t, payload.AutoUpdate)
}

func TestBuildPayloadFromFlags_CustomComposeFile(t *testing.T) {
	createSwarmGitName = "test-stack"
	createSwarmGitRepositoryURL = "https://github.com/user/repo"
	createSwarmGitSwarmID = "jpofkc0i9uo9wtx1zesuk649w"
	createSwarmGitComposeFile = "custom-compose.yml"
	createSwarmGitRepositoryReferenceName = "refs/heads/develop"
	createSwarmGitTLSSkipVerify = true
	createSwarmGitRepositoryUsername = ""
	createSwarmGitRepositoryPassword = ""
	createSwarmGitEnv = []string{}
	createSwarmGitAdditionalFiles = []string{}
	createSwarmGitAutoUpdateInterval = ""
	createSwarmGitAutoUpdateWebhook = ""

	payload := buildPayloadFromFlags()

	assert.Equal(t, "custom-compose.yml", payload.ComposeFile)
	assert.Equal(t, "refs/heads/develop", payload.RepositoryReferenceName)
	assert.True(t, payload.TLSSkipVerify)
}
