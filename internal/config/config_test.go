package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// NOTE: t.Parallel() is intentionally omitted. t.Setenv() in subtests is
	// incompatible with parallel execution at any level of the test hierarchy.

	testCases := []struct {
		name        string
		setupEnv    func(t *testing.T)
		expectedErr error
		verify      func(t *testing.T, cfg *Config)
	}{
		{
			name:        "default values are applied when no env vars set",
			setupEnv:    func(t *testing.T) {},
			expectedErr: nil,
			verify: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "8080", cfg.AppPort)
				assert.Equal(t, "bookmark-management", cfg.ServiceName)
				assert.Equal(t, "info", cfg.LogLevel)
				// InstanceID is auto-generated when empty
				assert.NotEmpty(t, cfg.InstanceID)
			},
		},
		{
			name: "env vars override defaults",
			setupEnv: func(t *testing.T) {
				t.Setenv("APP_PORT", "9090")
				t.Setenv("SERVICE_NAME", "my-service")
				t.Setenv("LOG_LEVEL", "debug")
				t.Setenv("INSTANCE_ID", "test-instance-id")
			},
			expectedErr: nil,
			verify: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "9090", cfg.AppPort)
				assert.Equal(t, "my-service", cfg.ServiceName)
				assert.Equal(t, "debug", cfg.LogLevel)
				assert.Equal(t, "test-instance-id", cfg.InstanceID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// NOTE: t.Parallel() cannot be used here because subtests call
			// t.Setenv(), which is incompatible with parallel execution.
			tc.setupEnv(t)

			cfg, err := NewConfig()

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			tc.verify(t, cfg)
		})
	}
}
