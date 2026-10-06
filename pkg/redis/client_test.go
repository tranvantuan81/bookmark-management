package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	// NOTE: t.Parallel() is intentionally omitted. t.Setenv() in subtests is
	// incompatible with parallel execution at any level of the test hierarchy.

	testCases := []struct {
		name        string
		setupEnv    func(t *testing.T)
		expectedErr error
		verify      func(t *testing.T, addr string)
	}{
		{
			name:        "default config creates a client",
			setupEnv:    func(t *testing.T) {},
			expectedErr: nil,
			verify: func(t *testing.T, addr string) {
				assert.Equal(t, "localhost:6379", addr)
			},
		},
		{
			name: "custom address is applied",
			setupEnv: func(t *testing.T) {
				t.Setenv("REDIS_ADDRESS", "localhost:6380")
			},
			expectedErr: nil,
			verify: func(t *testing.T, addr string) {
				assert.Equal(t, "localhost:6380", addr)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// NOTE: t.Parallel() cannot be used here because subtests call
			// t.Setenv(), which is incompatible with parallel execution.
			tc.setupEnv(t)

			client, err := NewClient("")

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, client)
			tc.verify(t, client.Options().Addr)
		})
	}
}
