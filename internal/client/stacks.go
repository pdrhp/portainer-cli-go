package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pdrhp/portainer-go-cli/pkg/types"
)

func (c *Client) ListStacks(ctx context.Context, filters *types.StackFilters) ([]types.Stack, error) {
	path := "/api/stacks"

	if filters != nil && (filters.EndpointID > 0 || filters.SwarmID != "") {
		params := url.Values{}
		if filters.EndpointID > 0 {
			params.Add("EndpointID", fmt.Sprintf("%d", filters.EndpointID))
		}
		if filters.SwarmID != "" {
			params.Add("SwarmID", filters.SwarmID)
		}
		path += "?" + params.Encode()
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

func (c *Client) RedeployStackFromGit(ctx context.Context, stackID int, endpointID int, payload types.StackGitRedeployPayload) (*types.Stack, error) {
	path := fmt.Sprintf("/api/stacks/%d/git/redeploy", stackID)
	if endpointID > 0 {
		path += fmt.Sprintf("?endpointId=%d", endpointID)
	}

	var stack types.Stack
	err := c.doRequest(ctx, "PUT", path, payload, &stack)
	if err != nil {
		return nil, fmt.Errorf("failed to redeploy stack from git: %w", err)
	}

	return &stack, nil
}
