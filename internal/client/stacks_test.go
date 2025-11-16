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

func TestClient_CreateSwarmStackFromGit_Success(t *testing.T) {
	expectedStack := types.Stack{
		ID:         1,
		Name:       "test-stack",
		Type:       types.StackTypeDockerSwarm,
		EndpointID: 1,
		SwarmID:    "jpofkc0i9uo9wtx1zesuk649w",
		Status:     1,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/stacks/create/swarm/repository" {
			t.Errorf("expected path /api/stacks/create/swarm/repository, got %s", r.URL.Path)
		}

		// Check query parameter
		endpointID := r.URL.Query().Get("endpointId")
		if endpointID != "1" {
			t.Errorf("expected endpointId=1, got %s", endpointID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedStack)
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("test-token")

	payload := types.StackCreateSwarmGitPayload{
		Name:                     "test-stack",
		RepositoryURL:            "https://github.com/user/repo",
		SwarmID:                  "jpofkc0i9uo9wtx1zesuk649w",
		ComposeFile:              "docker-compose.yml",
		RepositoryReferenceName:  "refs/heads/master",
		TLSSkipVerify:            false,
		RepositoryAuthentication: false,
		FromAppTemplate:          false,
	}

	result, err := client.CreateSwarmStackFromGit(context.Background(), 1, payload)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStack.ID, result.ID)
	assert.Equal(t, expectedStack.Name, result.Name)
	assert.Equal(t, expectedStack.Type, result.Type)
}

func TestClient_CreateSwarmStackFromGit_WithAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload types.StackCreateSwarmGitPayload
		json.NewDecoder(r.Body).Decode(&payload)

		assert.True(t, payload.RepositoryAuthentication)
		assert.Equal(t, "testuser", payload.RepositoryUsername)
		assert.Equal(t, "testpass", payload.RepositoryPassword)

		stack := types.Stack{
			ID:   1,
			Name: payload.Name,
			Type: types.StackTypeDockerSwarm,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stack)
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("test-token")

	payload := types.StackCreateSwarmGitPayload{
		Name:                     "test-stack",
		RepositoryURL:            "https://github.com/user/repo",
		SwarmID:                  "jpofkc0i9uo9wtx1zesuk649w",
		RepositoryAuthentication: true,
		RepositoryUsername:       "testuser",
		RepositoryPassword:       "testpass",
	}

	result, err := client.CreateSwarmStackFromGit(context.Background(), 1, payload)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestClient_CreateSwarmStackFromGit_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid swarm ID"}`))
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("test-token")

	payload := types.StackCreateSwarmGitPayload{
		Name:          "test-stack",
		RepositoryURL: "https://github.com/user/repo",
		SwarmID:       "invalid",
	}

	_, err := client.CreateSwarmStackFromGit(context.Background(), 1, payload)

	require.Error(t, err)
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}

func TestClient_CreateSwarmStackFromGit_Conflict(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`Stack name already exists`))
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("test-token")

	payload := types.StackCreateSwarmGitPayload{
		Name:          "existing-stack",
		RepositoryURL: "https://github.com/user/repo",
		SwarmID:       "jpofkc0i9uo9wtx1zesuk649w",
	}

	_, err := client.CreateSwarmStackFromGit(context.Background(), 1, payload)

	require.Error(t, err)
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 409, httpErr.StatusCode)
}

func TestClient_CreateSwarmStackFromGit_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := New(server.URL)
	client.SetToken("invalid-token")

	payload := types.StackCreateSwarmGitPayload{
		Name:          "test-stack",
		RepositoryURL: "https://github.com/user/repo",
		SwarmID:       "jpofkc0i9uo9wtx1zesuk649w",
	}

	_, err := client.CreateSwarmStackFromGit(context.Background(), 1, payload)

	require.Error(t, err)
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 401, httpErr.StatusCode)
}
