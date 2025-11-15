package client

import (
	"context"
	"fmt"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	JWT string `json:"jwt"`
}

func (c *Client) Authenticate(ctx context.Context, username, password string) (string, error) {
	req := AuthRequest{
		Username: username,
		Password: password,
	}

	var resp AuthResponse
	err := c.doRequest(ctx, "POST", "/api/auth", req, &resp)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	if resp.JWT == "" {
		return "", fmt.Errorf("no JWT token received in response")
	}

	return resp.JWT, nil
}
