package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthCmd_Flags(t *testing.T) {
	cmd := authCmd

	usernameFlag := cmd.Flag("username")
	passwordFlag := cmd.Flag("password")

	assert.NotNil(t, usernameFlag)
	assert.NotNil(t, passwordFlag)
	assert.Equal(t, "", usernameFlag.DefValue)
	assert.Equal(t, "", passwordFlag.DefValue)
}
