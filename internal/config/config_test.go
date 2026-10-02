package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/config"
)

// NewConfig reads from the environment, so this test must not run in parallel:
// t.Setenv panics when the test has called t.Parallel.
func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name     string
		setUpEnv func(t *testing.T)
		wantApp  string
		wantSvc  string
		wantUUID bool
	}{
		{
			name: "falls back to defaults when nothing is set",
			setUpEnv: func(t *testing.T) {
				unsetEnv(t, "APP_PORT")
				unsetEnv(t, "SERVICE_NAME")
				unsetEnv(t, "INSTANCE_ID")
			},
			wantApp:  "8080",
			wantSvc:  "bookmark-management",
			wantUUID: true,
		},
		{
			name: "reads values from the environment",
			setUpEnv: func(t *testing.T) {
				t.Setenv("APP_PORT", "9090")
				t.Setenv("SERVICE_NAME", "bookmark-management-test")
				t.Setenv("INSTANCE_ID", "0f1092e7-38ed-4701-b500-c4697c9dc122")
			},
			wantApp:  "9090",
			wantSvc:  "bookmark-management-test",
			wantUUID: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setUpEnv(t)

			cfg, err := config.NewConfig()

			require.NoError(t, err)
			assert.Equal(t, tc.wantApp, cfg.AppPort)
			assert.Equal(t, tc.wantSvc, cfg.ServiceName)

			if tc.wantUUID {
				// a uuid is generated once per process when INSTANCE_ID is absent
				assert.NotEmpty(t, cfg.InstanceID)
				return
			}
			assert.Equal(t, os.Getenv("INSTANCE_ID"), cfg.InstanceID)
		})
	}
}

// unsetEnv removes an env var for the duration of the test and restores it afterwards.
func unsetEnv(t *testing.T, key string) {
	t.Helper()

	if old, ok := os.LookupEnv(key); ok {
		t.Cleanup(func() { _ = os.Setenv(key, old) })
	}

	require.NoError(t, os.Unsetenv(key))
}
