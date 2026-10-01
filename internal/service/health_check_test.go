package service_test

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	"github.com/tranvantuan81/bookmark-management/internal/repository"
	"github.com/tranvantuan81/bookmark-management/internal/service"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestHealthCheck_Service(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		setUpRedis func(t *testing.T) *redis.Client
		wantRes    service.Response
		wantErr    error
	}{
		{
			name: "redis is reachable",
			setUpRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			wantRes: service.Response{
				Message:     "OK",
				ServiceName: "bookmark-management",
				InstanceID:  "0f1092e7-38ed-4701-b500-c4697c9dc122",
			},
			wantErr: nil,
		},
		{
			name: "redis is unreachable",
			setUpRedis: func(t *testing.T) *redis.Client {
				client := redisPkg.InitMockRedis(t)
				require.NoError(t, client.Close())
				return client
			},
			wantRes: service.Response{},
			wantErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			cfg := &config.Config{
				ServiceName: "bookmark-management",
				InstanceID:  "0f1092e7-38ed-4701-b500-c4697c9dc122",
			}

			ping := repository.NewPing(tc.setUpRedis(t))
			svc := service.NewHealthCheck(cfg, ping)

			res, err := svc.HealthCheck(ctx)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
