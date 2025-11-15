package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pdrhp/portainer-go-cli/pkg/types"
)

func (c *Client) ListStacks(ctx context.Context, filters *types.StackFilters) ([]types.Stack, error) {
	path := "/api/stacks"

	if filters != nil && (filters.EndpointID > 0 || filters.SwarmID != "") {
		// Criar map apenas com campos não-zero
		filterMap := make(map[string]interface{})
		if filters.EndpointID > 0 {
			filterMap["EndpointId"] = filters.EndpointID
		}
		if filters.SwarmID != "" {
			filterMap["SwarmId"] = filters.SwarmID
		}

		filtersJSON, err := json.Marshal(filterMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal filters: %w", err)
		}

		// Não fazer URL encoding - Portainer pode aceitar JSON puro
		path += "?filters=" + string(filtersJSON)
	}

	var stacks []types.Stack
	err := c.doRequest(ctx, "GET", path, nil, &stacks)
	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}

	return stacks, nil
}
