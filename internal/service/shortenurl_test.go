package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

func TestGeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		expectedLength int
		expectedError  error
	}{
		{
			name:           "success",
			expectedLength: 12,
			expectedError:  nil,
		},
		{
			name:           "success with custom length",
			expectedLength: 1,
			expectedError:  nil,
		},
		{
			name:           "success with custom length",
			expectedLength: 100000,
			expectedError:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			password, err := service.GeneratePassword(tc.expectedLength)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedLength, len(password))
		})
	}
}

//func TestShortenUrl_CreateShortenLink(t *testing.T) {
//	t.Parallel()
//
//	testCases := []struct {
//		name         string
//		setupMock    func(ctx context.Context) *redis.Client
//		expectedCode string
//		expectedErr  error
//	}{
//		{
//			name: "code is generated and stored",
//			setupMock: func(ctx context.Context) *redis.Client {
//				testDB := redisPkg.InitMockRedis(t)
//				testDB.Set(ctx, "abcdef", "https://google.com", 60)
//				return testDB
//			},
//			expectedCode: "",
//			expectedErr:  nil,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			t.Parallel()
//			ctx := t.Context()
//
//			code, _ := service.GeneratePassword(7)
//			mockRedis := tc.setupMock(ctx)
//
//			testRepo := repository.NewURLStorage(mockRedis)
//
//			testSvc := service.NewShortenUrl(testRepo)
//
//			result, err := testSvc.CreateShortenLink(ctx, code, 60)
//			assert.NoError(t, tc.expectedErr, err)
//			assert.Equal(t, tc.expectedCode, result)
//		})
//	}
//}
