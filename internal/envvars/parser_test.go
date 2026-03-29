package envvars

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_PreservesValueAfterEquals(t *testing.T) {
	vars, err := Parse("GREETING=  hello world  ")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "GREETING", vars[0].Name)
	assert.Equal(t, "  hello world  ", vars[0].Value)
}

func TestParse_TrimsOnlyKeyAroundEquals(t *testing.T) {
	vars, err := Parse("  GREETING   =hello")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "GREETING", vars[0].Name)
	assert.Equal(t, "hello", vars[0].Value)
}

func TestParse_QuotedDoubleValue(t *testing.T) {
	vars, err := Parse("GREETING=\"hello world\"")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "GREETING", vars[0].Name)
	assert.Equal(t, "hello world", vars[0].Value)
}

func TestParse_QuotedSingleValue(t *testing.T) {
	vars, err := Parse("GREETING='hello world'")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "GREETING", vars[0].Name)
	assert.Equal(t, "hello world", vars[0].Value)
}

func TestParse_AllowsSingleQuotedAssignment(t *testing.T) {
	vars, err := Parse("'DB_HOST=postgres'")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "DB_HOST", vars[0].Name)
	assert.Equal(t, "postgres", vars[0].Value)
}

func TestParse_AllowsDoubleQuotedAssignment(t *testing.T) {
	vars, err := Parse("\"DB_HOST=postgres\"")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "DB_HOST", vars[0].Name)
	assert.Equal(t, "postgres", vars[0].Value)
}

func TestParse_EmptyValue(t *testing.T) {
	vars, err := Parse("EMPTY=")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "EMPTY", vars[0].Name)
	assert.Equal(t, "", vars[0].Value)
}

func TestParse_ValueCanContainEquals(t *testing.T) {
	vars, err := Parse("URL=https://example.com?a=b")

	require.NoError(t, err)
	require.Len(t, vars, 1)
	assert.Equal(t, "URL", vars[0].Name)
	assert.Equal(t, "https://example.com?a=b", vars[0].Value)
}

func TestParse_MissingEqualsReturnsError(t *testing.T) {
	_, err := Parse("FIRST=ok\nMISSING_EQUALS")

	require.Error(t, err)
	assert.ErrorContains(t, err, "line 2")
}

func TestParse_InvalidKeyReturnsError(t *testing.T) {
	_, err := Parse("1INVALID=value")

	require.Error(t, err)
}

func TestParse_MultipleLines(t *testing.T) {
	input := "FIRST=one\n\nSECOND=two words\nTHIRD=\"three\""

	vars, err := Parse(input)

	require.NoError(t, err)
	require.Len(t, vars, 3)

	assert.Equal(t, "FIRST", vars[0].Name)
	assert.Equal(t, "one", vars[0].Value)
	assert.Equal(t, "SECOND", vars[1].Name)
	assert.Equal(t, "two words", vars[1].Value)
	assert.Equal(t, "THIRD", vars[2].Name)
	assert.Equal(t, "three", vars[2].Value)
}

func TestParse_CRLFInput(t *testing.T) {
	vars, err := Parse("FIRST=one\r\nSECOND=two\r\n")

	require.NoError(t, err)
	require.Len(t, vars, 2)
	assert.Equal(t, "FIRST", vars[0].Name)
	assert.Equal(t, "one", vars[0].Value)
	assert.Equal(t, "SECOND", vars[1].Name)
	assert.Equal(t, "two", vars[1].Value)
}
