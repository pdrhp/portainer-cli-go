package client

import (
	"context"
	"encoding/base64"
	"fmt"
)

// ConfigCreateResponse identifies the newly created Swarm config.
type ConfigCreateResponse struct {
	ID string `json:"ID" yaml:"ID"`
}

// CreateConfig creates an immutable Docker Swarm config through Portainer.
func (c *Client) CreateConfig(ctx context.Context, endpointID int, name string, data []byte) (*ConfigCreateResponse, error) {
	payload := struct {
		Name string `json:"Name"`
		Data string `json:"Data"`
	}{name, base64.StdEncoding.EncodeToString(data)}
	var result ConfigCreateResponse
	path := fmt.Sprintf("/api/endpoints/%d/docker/configs/create", endpointID)
	if err := c.doRequest(ctx, "POST", path, payload, &result); err != nil {
		return nil, fmt.Errorf("failed to create Docker config: %w", err)
	}
	return &result, nil
}
