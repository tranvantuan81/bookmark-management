package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/app/repository"
	repoMocks "github.com/tranvantuan81/bookmark-management/internal/app/repository/mocks"
	"github.com/tranvantuan81/bookmark-management/internal/app/service"
	svcMocks "github.com/tranvantuan81/bookmark-management/internal/app/service/mocks"
)

var errSomething = errors.New("something went wrong")

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		length         int
		expectedLength int
		expectedErr    error
	}{
		{
			name:           "default length",
			length:         service.DefaultCodeLength,
			expectedLength: service.DefaultCodeLength,
			expectedErr:    nil,
		},
		{
			name:           "minimum length",
			length:         1,
			expectedLength: 1,
			expectedErr:    nil,
		},
		{
			name:           "large length",
			length:         100,
			expectedLength: 100,
			expectedErr:    nil,
		},
		{
			name:           "zero length produces empty string",
			length:         0,
			expectedLength: 0,
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// test the exported GenerateCode function
			code, err := service.GenerateCode(tc.length)
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Len(t, code, tc.expectedLength)

			// test the CodeGen interface implementation (NewCodeGen)
			codeGen := service.NewCodeGen()
			code2, err2 := codeGen.GenerateCode(tc.length)
			assert.ErrorIs(t, err2, tc.expectedErr)
			assert.Len(t, code2, tc.expectedLength)
		})
	}
}

// TestGenerateCode_Charset verifies that every character in the generated code
// belongs to the allowed charset (alphanumeric).
func TestGenerateCode_Charset(t *testing.T) {
	t.Parallel()

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetSet := make(map[byte]struct{}, len(charset))
	for i := range len(charset) {
		charsetSet[charset[i]] = struct{}{}
	}

	code, err := service.GenerateCode(50)
	assert.NoError(t, err)
	assert.Len(t, code, 50)
	for _, ch := range []byte(code) {
		_, ok := charsetSet[ch]
		assert.True(t, ok, "unexpected character %q in generated code", string(ch))
	}
}

// TestGenerateCode_Unique verifies that two consecutive calls produce different codes.
func TestGenerateCode_Unique(t *testing.T) {
	t.Parallel()

	code1, err1 := service.GenerateCode(service.DefaultCodeLength)
	code2, err2 := service.GenerateCode(service.DefaultCodeLength)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// While theoretically possible for two codes to collide, the probability is
	// negligible (62^7 ≈ 3.5 × 10^12) and would indicate a broken RNG.
	assert.NotEqual(t, code1, code2)
}

func TestShortenUrl_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMocks  func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.CodeGen)
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.CodeGen) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewCodeGen(t)

				mockCodeGen.On("GenerateCode", service.DefaultCodeLength).Return("abcdefg", nil).Once()
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", repository.ErrNotFound).Once()
				mockRepo.On("StoreURL", ctx, "abcdefg", "https://google.com", int64(3600)).Return(nil).Once()

				return mockRepo, mockCodeGen
			},
			expectedErr: nil,
		},
		{
			name: "code already exists so a new one is generated",
			setupMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.CodeGen) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewCodeGen(t)

				mockCodeGen.On("GenerateCode", service.DefaultCodeLength).Return("dupe123", nil).Once()
				mockCodeGen.On("GenerateCode", service.DefaultCodeLength).Return("fresh45", nil).Once()
				mockRepo.On("GetURL", ctx, "dupe123").Return("https://old.com", nil).Once()
				mockRepo.On("GetURL", ctx, "fresh45").Return("", repository.ErrNotFound).Once()
				mockRepo.On("StoreURL", ctx, "fresh45", "https://google.com", int64(3600)).Return(nil).Once()

				return mockRepo, mockCodeGen
			},
			expectedErr: nil,
		},
		{
			name: "code generation fails",
			setupMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.CodeGen) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewCodeGen(t)

				mockCodeGen.On("GenerateCode", service.DefaultCodeLength).Return("", errSomething).Once()

				return mockRepo, mockCodeGen
			},
			expectedErr: errSomething,
		},
		{
			name: "storage lookup fails",
			setupMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.CodeGen) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewCodeGen(t)

				mockCodeGen.On("GenerateCode", service.DefaultCodeLength).Return("abcdefg", nil).Once()
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", errSomething).Once()

				return mockRepo, mockCodeGen
			},
			expectedErr: errSomething,
		},
		{
			name: "storing the url fails",
			setupMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.CodeGen) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewCodeGen(t)

				mockCodeGen.On("GenerateCode", service.DefaultCodeLength).Return("abcdefg", nil).Once()
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", repository.ErrNotFound).Once()
				mockRepo.On("StoreURL", ctx, "abcdefg", "https://google.com", int64(3600)).Return(errSomething).Once()

				return mockRepo, mockCodeGen
			},
			expectedErr: errSomething,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockRepo, mockCodeGen := tc.setupMocks(t, ctx)
			svc := service.NewShortenUrl(mockRepo, mockCodeGen)

			code, err := svc.CreateShortenLink(ctx, "https://google.com", 3600)

			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				assert.Len(t, code, service.DefaultCodeLength)
			} else {
				assert.Empty(t, code)
			}
		})
	}
}

func TestShortenUrl_GetURLFromCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		code        string
		setupMocks  func(t *testing.T, ctx context.Context) *repoMocks.URLStorage
		expectedURL string
		expectedErr error
	}{
		{
			name: "code exists",
			code: "abcdefg",
			setupMocks: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockRepo := repoMocks.NewURLStorage(t)
				mockRepo.On("GetURL", ctx, "abcdefg").Return("https://google.com", nil).Once()
				return mockRepo
			},
			expectedURL: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "code does not exist",
			code: "notfound",
			setupMocks: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockRepo := repoMocks.NewURLStorage(t)
				mockRepo.On("GetURL", ctx, "notfound").Return("", repository.ErrNotFound).Once()
				return mockRepo
			},
			expectedURL: "",
			expectedErr: repository.ErrNotFound,
		},
		{
			name: "storage returns an error",
			code: "abcdefg",
			setupMocks: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockRepo := repoMocks.NewURLStorage(t)
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", errSomething).Once()
				return mockRepo
			},
			expectedURL: "",
			expectedErr: errSomething,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockRepo := tc.setupMocks(t, ctx)
			svc := service.NewShortenUrl(mockRepo, svcMocks.NewCodeGen(t))

			url, err := svc.GetURLFromCode(ctx, tc.code)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}
