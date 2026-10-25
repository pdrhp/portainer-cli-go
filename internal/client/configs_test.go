package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/api/endpoints/7/docker/configs/create", r.URL.Path)
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "nginx-v1", body["Name"])
		require.Equal(t, "aGVsbG8K", body["Data"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ID":"config-123"}`))
	}))
	defer server.Close()
	c := New(server.URL)
	c.SetToken("test-token")
	result, err := c.CreateConfig(context.Background(), 7, "nginx-v1", []byte("hello\n"))
	require.NoError(t, err)
	require.Equal(t, "config-123", result.ID)
}

func TestCreateConfigPreservesAPIErrors(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 409, 500} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"message":"config rejected"}`))
			}))
			defer server.Close()
			_, err := New(server.URL).CreateConfig(context.Background(), 1, "test", []byte("data"))
			var httpErr *HTTPError
			require.ErrorAs(t, err, &httpErr)
			require.Equal(t, status, httpErr.StatusCode)
			require.Contains(t, err.Error(), "config rejected")
		})
	}
}
