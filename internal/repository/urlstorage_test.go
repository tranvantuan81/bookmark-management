package repository

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "abcdef").Result()
				assert.NoError(t, err)
				assert.Equal(t, "https://google.com", res)
			},
		},
		{
			name: "connection error",
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
			mock := tc.setupMock(ctx)

			testRepo := NewURLStorage(mock)
			err := testRepo.StoreURL(ctx, "abcdef", "https://google.com", 0)
			assert.ErrorIs(t, err, tc.expectedErr)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}

func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		setupMock func(ctx context.Context) *redis.Client

		expectedErr error
		expectedUrl string
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				testDB := redisPkg.InitMockRedis(t)
				testDB.Set(ctx, "abcdef", "https://google.com", 60)
				return testDB
			},
			expectedErr: nil,
			expectedUrl: "https://google.com",
		},
		{
			name: "case key == nil",
			setupMock: func(ctx context.Context) *redis.Client {
				testDB := redisPkg.InitMockRedis(t)
				return testDB
			},
			expectedErr: ErrNotFound,
		},
		{
			name: "connection error",
			setupMock: func(ctx context.Context) *redis.Client {
				testDB := redisPkg.InitMockRedis(t)
				testDB.Close()
				return testDB
			},
			expectedErr: redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx)

			testRepo := NewURLStorage(mock)
			url, err := testRepo.GetURL(ctx, "abcdef")
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedUrl, url)
		})
	}
}
