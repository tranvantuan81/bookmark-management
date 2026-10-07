package repository

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestPing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(t *testing.T) *redis.Client
		expectedErr error
	}{
		{
			name: "redis is reachable",
			setupMock: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedErr: nil,
		},
		{
			name: "redis is unreachable",
			setupMock: func(t *testing.T) *redis.Client {
				client := redisPkg.InitMockRedis(t)
				client.Close()
				return client
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			ping := NewPing(tc.setupMock(t))

			err := ping.Ping(ctx)

			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
