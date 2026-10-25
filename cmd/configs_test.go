package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfigsCreate(t *testing.T) {
	for _, source := range []string{"file", "stdin"} {
		t.Run(source, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			viper.Set("token", "test-token")
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/endpoints/1/docker/configs/create", r.URL.Path)
				var body map[string]interface{}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, "nginx-v1", body["Name"])
				require.Equal(t, "AP9oZWxsbwo=", body["Data"])
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"ID":"created-id"}`))
			}))
			defer server.Close()
			file := "-"
			data := []byte{0, 255, 'h', 'e', 'l', 'l', 'o', '\n'}
			if source == "file" {
				file = filepath.Join(t.TempDir(), "config.bin")
				require.NoError(t, os.WriteFile(file, data, 0600))
			}
			root := &cobra.Command{Use: "test"}
			root.PersistentFlags().String("server-url", server.URL, "")
			root.PersistentFlags().String("output", "json", "")
			root.AddCommand(newConfigsCommand())
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetIn(bytes.NewReader(data))
			root.SetArgs([]string{"configs", "create", "nginx-v1", "--endpoint-id", "1", "--file", file})
			require.NoError(t, root.Execute())
			require.JSONEq(t, `{"ID":"created-id"}`, out.String())
		})
	}
}

func TestConfigsCreateRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		message string
	}{
		{"missing name", []string{"--endpoint-id", "1", "--file", "-"}, "accepts 1 arg"},
		{"blank name", []string{" ", "--endpoint-id", "1", "--file", "-"}, "name must not be empty"},
		{"invalid endpoint", []string{"test", "--endpoint-id", "0", "--file", "-"}, "endpoint-id must be positive"},
		{"missing file flag", []string{"test", "--endpoint-id", "1"}, "required flag"},
		{"missing file", []string{"test", "--endpoint-id", "1", "--file", "/nonexistent/config-file"}, "read config file"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			cmd.PersistentFlags().String("output", "table", "")
			cmd.PersistentFlags().String("server-url", "", "")
			cmd.AddCommand(newConfigsCommand())
			cmd.SetArgs(append([]string{"configs", "create"}, tc.args...))
			cmd.SetIn(strings.NewReader("hello"))
			err := cmd.Execute()
			require.ErrorContains(t, err, tc.message)
		})
	}
}

func TestConfigsCreateReportsAPIErrors(t *testing.T) {
	for _, status := range []int{401, 403, 409} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			viper.Set("token", "test-token")
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"message":"creation rejected"}`))
			}))
			defer server.Close()
			root := &cobra.Command{Use: "test", SilenceUsage: true, SilenceErrors: true}
			root.PersistentFlags().String("server-url", server.URL, "")
			root.PersistentFlags().String("output", "json", "")
			root.AddCommand(newConfigsCommand())
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetIn(strings.NewReader("hello"))
			root.SetArgs([]string{"configs", "create", "test", "--endpoint-id", "1", "--file", "-"})
			require.ErrorContains(t, root.Execute(), "creation rejected")
			require.Empty(t, out.String())
		})
	}
}
