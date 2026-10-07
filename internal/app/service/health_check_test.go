package service

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/app/repository"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		cfg         *config.Config
		setupMock   func(ctx context.Context) *redis.Client
		expectedRes Response
		expectedErr error
	}{
		{
			name: "success",
			cfg: &config.Config{
				ServiceName: "test",
				InstanceID:  "xxx",
			},
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedRes: Response{
				Message:     "OK",
				ServiceName: "test",
				InstanceID:  "xxx",
			},
			expectedErr: nil,
		},
		{
			name: "connection err",
			cfg: &config.Config{
				ServiceName: "test",
				InstanceID:  "xxx",
			},
			setupMock: func(ctx context.Context) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				mockR.Close()
				return mockR
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisClient := tc.setupMock(ctx)

			ping := repository.NewPing(redisClient)

			healthCheckSev := NewHealthCheck(tc.cfg, ping)

			res, err := healthCheckSev.HealthCheck(ctx)

			assert.Equal(t, tc.expectedRes, res)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
