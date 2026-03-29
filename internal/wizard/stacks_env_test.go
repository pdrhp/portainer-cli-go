package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWizardEnvParse_CronWithSpaces(t *testing.T) {
	vars, err := parseWizardEnvText("CRON=0 0 * * *\nTZ=UTC")

	require.NoError(t, err)
	require.Len(t, vars, 2)
	assert.Equal(t, "CRON", vars[0].Name)
	assert.Equal(t, "0 0 * * *", vars[0].Value)
	assert.Equal(t, "TZ", vars[1].Name)
	assert.Equal(t, "UTC", vars[1].Value)
}

func TestWizardEnvParse_ErrorPreservesLineContext(t *testing.T) {
	_, err := parseWizardEnvText("VALID=1\nINVALID")

	require.Error(t, err)
	assert.ErrorContains(t, err, "line 2")
}

func TestValidateGitAuthInput_AllowsNonEmptyCredentials(t *testing.T) {
	err := validateGitAuthInput("user", "pass")
	require.NoError(t, err)
}

func TestValidateGitAuthInput_RejectsMissingUsername(t *testing.T) {
	err := validateGitAuthInput("", "pass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "required")
}

func TestValidateGitAuthInput_RejectsMissingPassword(t *testing.T) {
	err := validateGitAuthInput("user", "")
	require.Error(t, err)
	assert.ErrorContains(t, err, "required")
}
