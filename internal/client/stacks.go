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

		path += "?filters=" + string(filtersJSON)
	}

	var stacks []types.Stack
	err := c.doRequest(ctx, "GET", path, nil, &stacks)
	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}

	return stacks, nil
}

func (c *Client) CreateSwarmStackFromGit(ctx context.Context, endpointID int, payload types.StackCreateSwarmGitPayload) (*types.Stack, error) {
	path := fmt.Sprintf("/api/stacks/create/swarm/repository?endpointId=%d", endpointID)

	var stack types.Stack
	err := c.doRequest(ctx, "POST", path, payload, &stack)
	if err != nil {
		return nil, fmt.Errorf("failed to create swarm stack from git: %w", err)
	}

	return &stack, nil
}
