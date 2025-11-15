package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pdrhp/portainer-go-cli/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListStacks_Success(t *testing.T) {
	stacks := []types.Stack{
		{
			ID:         1,
			Name:       "web-stack",
			Type:       types.StackTypeDockerCompose,
			EndpointID: 1,
			SwarmID:    "jpofkc0i9uo9wtx1zesuk649w",
			Status:     1,
		},
		{
			ID:         2,
			Name:       "api-stack",
			Type:       types.StackTypeDockerSwarm,
			EndpointID: 1,
			SwarmID:    "jpofkc0i9uo9wtx1zesuk649w",
			Status:     1,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/api/stacks" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stacks)
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("test-token")

	result, err := client.ListStacks(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "web-stack", result[0].Name)
	assert.Equal(t, "api-stack", result[1].Name)
	assert.Equal(t, types.StackTypeDockerCompose, result[0].Type)
	assert.Equal(t, types.StackTypeDockerSwarm, result[1].Type)
}

func TestClient_ListStacks_WithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedFilters := `{"EndpointId":1,"SwarmId":"test-swarm-id"}`
		if r.URL.RawQuery != "filters="+expectedFilters {
			t.Errorf("unexpected filters: %s", r.URL.RawQuery)
		}

		stacks := []types.Stack{
			{
				ID:         1,
				Name:       "filtered-stack",
				Type:       types.StackTypeDockerSwarm,
				EndpointID: 1,
				SwarmID:    "test-swarm-id",
				Status:     1,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stacks)
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("test-token")

	filters := &types.StackFilters{
		EndpointID: 1,
		SwarmID:    "test-swarm-id",
	}

	result, err := client.ListStacks(context.Background(), filters)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "filtered-stack", result[0].Name)
}

func TestClient_ListStacks_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("invalid-token")

	_, err := client.ListStacks(context.Background(), nil)

	require.Error(t, err)
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 401, httpErr.StatusCode)
}
