package envvars

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pdrhp/portainer-go-cli/pkg/types"
)

var variableKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func Parse(input string) ([]types.Pair, error) {
	if input == "" {
		return nil, nil
	}

	lines := strings.Split(input, "\n")
	variables := make([]types.Pair, 0, len(lines))

	for lineNumber, rawLine := range lines {
		line := strings.TrimSuffix(rawLine, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}

		lineToParse := line
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) >= 2 {
			if (trimmedLine[0] == '"' && trimmedLine[len(trimmedLine)-1] == '"') || (trimmedLine[0] == '\'' && trimmedLine[len(trimmedLine)-1] == '\'') {
				lineToParse = trimmedLine[1 : len(trimmedLine)-1]
			}
		}

		equalsIndex := strings.Index(lineToParse, "=")
		if equalsIndex < 0 {
			return nil, fmt.Errorf("invalid env var format at line %d: missing '='", lineNumber+1)
		}

		key := strings.TrimSpace(lineToParse[:equalsIndex])
		if !variableKeyPattern.MatchString(key) {
			return nil, fmt.Errorf("invalid env var key at line %d: %q", lineNumber+1, key)
		}

		value := lineToParse[equalsIndex+1:]
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		variables = append(variables, types.Pair{Name: key, Value: value})
	}

	return variables, nil
}
