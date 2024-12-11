package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := &App{
		Listen:          ":8080",
		BaseShortURL:    "testBase",
		FileStoragePath: "fileStore",
		DBDsn:           "dbDsn",
		TrustedSubnet:   "subnet",
	}

	assert.Equal(t, "dbDsn", cfg.GetDsn())
	assert.Equal(t, "fileStore", cfg.GetStoragePath())
	assert.Equal(t, "testBase", cfg.GetBaseShortURL())
	assert.Equal(t, "subnet", cfg.GetTrustedSubnet())
}

func TestNewConfig(t *testing.T) {
	os.Clearenv()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		args    []string
		expect  struct {
			HTTPS   bool
			Listen  string
			BaseURL string
		}
	}{
		{
			name: "default_values",
			envVars: map[string]string{
				"SERVER_ADDRESS": "",
				"BASE_URL":       "",
			},
			args: []string{},
			expect: struct {
				HTTPS   bool
				Listen  string
				BaseURL string
			}{HTTPS: false, Listen: ":8080", BaseURL: "http://localhost:8080"},
		},
		{
			name: "environment_variables",
			envVars: map[string]string{
				"SERVER_ADDRESS": ":9090",
				"BASE_URL":       "http://ya.ru",
			},
			args: []string{},
			expect: struct {
				HTTPS   bool
				Listen  string
				BaseURL string
			}{HTTPS: false, Listen: ":9090", BaseURL: "http://ya.ru"},
		},
		{
			name: "environment_variables_with_https",
			envVars: map[string]string{
				"SERVER_ADDRESS": ":9090",
				"BASE_URL":       "http://ya.ru",
				"ENABLE_HTTPS":   "1",
			},
			args: []string{},
			expect: struct {
				HTTPS   bool
				Listen  string
				BaseURL string
			}{HTTPS: true, Listen: ":9090", BaseURL: "http://ya.ru"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("failed to set environment variable %s: %v", key, err)
				}
			}

			os.Args = append([]string{"cmd"}, tt.args...)

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := NewConfig()
			require.NoError(t, err)

			assert.Equal(t, "", cfg.GetDsn())
			assert.NotEmpty(t, cfg.GetStoragePath())
			assert.Equal(t, tt.expect.HTTPS, cfg.HTTPS.Enable)
			assert.Equal(t, tt.expect.Listen, cfg.Listen)
			assert.Equal(t, tt.expect.BaseURL, cfg.BaseShortURL)

			if tt.expect.HTTPS {
				assert.FileExists(t, cfg.HTTPS.Key)
				assert.FileExists(t, cfg.HTTPS.Pem)

				_ = os.Remove(cfg.HTTPS.Key)
				_ = os.Remove(cfg.HTTPS.Pem)
			}
		})
	}
}
