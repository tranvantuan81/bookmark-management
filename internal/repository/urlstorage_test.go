package repository_test

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/repository"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	const (
		key   = "abcdef"
		value = "https://google.com"
	)

	testCases := []struct {
		name       string
		setUpRedis func(t *testing.T) *redis.Client
		wantErr    error
	}{
		{
			name: "key is stored with the given ttl",
			setUpRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			wantErr: nil,
		},
		{
			name: "redis connection is closed",
			setUpRedis: func(t *testing.T) *redis.Client {
				client := redisPkg.InitMockRedis(t)
				require.NoError(t, client.Close())
				return client
			},
			wantErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			client := tc.setUpRedis(t)
			repo := repository.NewURLStorage(client)

			err := repo.StoreURL(ctx, key, value, 0)

			require.ErrorIs(t, err, tc.wantErr)

			if tc.wantErr != nil {
				return
			}

			// the value must be readable through the same storage
			got, err := repo.GetURL(ctx, key)
			require.NoError(t, err)
			assert.Equal(t, value, got)
		})
	}
}

func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()

	const (
		key   = "abcdef"
		value = "https://google.com"
	)

	testCases := []struct {
		name       string
		lookupKey  string
		setUpRedis func(t *testing.T, ctx context.Context) *redis.Client
		wantURL    string
		wantErr    error
	}{
		{
			name:      "key exists",
			lookupKey: key,
			setUpRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				client := redisPkg.InitMockRedis(t)
				require.NoError(t, client.Set(ctx, key, value, 0).Err())
				return client
			},
			wantURL: value,
			wantErr: nil,
		},
		{
			name:      "key does not exist",
			lookupKey: "not-exist",
			setUpRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			wantURL: "",
			wantErr: repository.ErrNotFound,
		},
		{
			name:      "redis connection is closed",
			lookupKey: key,
			setUpRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				client := redisPkg.InitMockRedis(t)
				require.NoError(t, client.Close())
				return client
			},
			wantURL: "",
			wantErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			client := tc.setUpRedis(t, ctx)
			repo := repository.NewURLStorage(client)

			got, err := repo.GetURL(ctx, tc.lookupKey)

			require.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantURL, got)
		})
	}
}
