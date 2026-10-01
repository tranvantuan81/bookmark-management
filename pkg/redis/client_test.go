package redis_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

// NewClient reads its options from the environment, so this test must not run in
// parallel: t.Setenv panics when the test has called t.Parallel.
func TestNewClient(t *testing.T) {
	testCases := []struct {
		name     string
		setUpEnv func(t *testing.T, addr string)
		wantErr  bool
	}{
		{
			name: "client connects to the configured address",
			setUpEnv: func(t *testing.T, addr string) {
				t.Setenv("REDIS_ADDRESS", addr)
				t.Setenv("REDIS_PASSWORD", "")
				t.Setenv("REDIS_DB", "0")
			},
			wantErr: false,
		},
		{
			name: "invalid REDIS_DB value is rejected",
			setUpEnv: func(t *testing.T, addr string) {
				t.Setenv("REDIS_ADDRESS", addr)
				t.Setenv("REDIS_DB", "not-a-number")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRedis := miniredis.RunT(t)
			tc.setUpEnv(t, mockRedis.Addr())

			client, err := redisPkg.NewClient("")

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, client)
				return
			}

			require.NoError(t, err)
			t.Cleanup(func() { _ = client.Close() })

			ctx := t.Context()

			require.NoError(t, client.Set(ctx, "key", "value", 0).Err())

			got, err := client.Get(ctx, "key").Result()
			assert.NoError(t, err)
			assert.Equal(t, "value", got)
		})
	}
}
