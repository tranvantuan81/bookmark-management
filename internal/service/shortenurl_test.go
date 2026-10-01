package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/repository"
	repoMocks "github.com/tranvantuan81/bookmark-management/internal/repository/mocks"
	"github.com/tranvantuan81/bookmark-management/internal/service"
	svcMocks "github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

// errSomething is a sentinel error used by the dependency mocks.
var errSomething = errors.New("something went wrong")

func TestShortenUrl_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		setUpMocks func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.GenPass)
		verifyFunc func(t *testing.T, r *repoMocks.URLStorage)
		wantCode   string
		wantErr    error
	}{
		{
			name: "code is generated and stored",
			setUpMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.GenPass) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewGenPass(t)

				mockCodeGen.On("GeneratePassword", service.DefaultCodeLength).Return("abcdefg", nil).Once()
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", repository.ErrNotFound).Once()
				mockRepo.On("StoreURL", ctx, "abcdefg", "https://google.com", 3600).Return(nil).Once()

				return mockRepo, mockCodeGen
			},
			wantCode: "abcdefg",
			wantErr:  nil,
		},
		{
			name: "code already exists so a new one is generated",
			setUpMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.GenPass) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewGenPass(t)

				mockCodeGen.On("GeneratePassword", service.DefaultCodeLength).Return("dupe123", nil).Once()
				mockCodeGen.On("GeneratePassword", service.DefaultCodeLength).Return("fresh45", nil).Once()
				mockRepo.On("GetURL", ctx, "dupe123").Return("https://old.com", nil).Once()
				mockRepo.On("GetURL", ctx, "fresh45").Return("", repository.ErrNotFound).Once()
				mockRepo.On("StoreURL", ctx, "fresh45", "https://google.com", 3600).Return(nil).Once()

				return mockRepo, mockCodeGen
			},
			wantCode: "fresh45",
			wantErr:  nil,
		},
		{
			name: "code generation fails",
			setUpMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.GenPass) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewGenPass(t)

				mockCodeGen.On("GeneratePassword", service.DefaultCodeLength).Return("", errSomething).Once()

				return mockRepo, mockCodeGen
			},
			wantCode: "",
			wantErr:  errSomething,
		},
		{
			name: "storage lookup fails",
			setUpMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.GenPass) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewGenPass(t)

				mockCodeGen.On("GeneratePassword", service.DefaultCodeLength).Return("abcdefg", nil).Once()
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", errSomething).Once()

				return mockRepo, mockCodeGen
			},
			verifyFunc: func(t *testing.T, r *repoMocks.URLStorage) {
				r.AssertNotCalled(t, "StoreURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
			wantCode: "",
			wantErr:  errSomething,
		},
		{
			name: "storing the url fails",
			setUpMocks: func(t *testing.T, ctx context.Context) (*repoMocks.URLStorage, *svcMocks.GenPass) {
				mockRepo := repoMocks.NewURLStorage(t)
				mockCodeGen := svcMocks.NewGenPass(t)

				mockCodeGen.On("GeneratePassword", service.DefaultCodeLength).Return("abcdefg", nil).Once()
				mockRepo.On("GetURL", ctx, "abcdefg").Return("", repository.ErrNotFound).Once()
				mockRepo.On("StoreURL", ctx, "abcdefg", "https://google.com", 3600).Return(errSomething).Once()

				return mockRepo, mockCodeGen
			},
			wantCode: "",
			wantErr:  errSomething,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockRepo, mockCodeGen := tc.setUpMocks(t, ctx)
			svc := service.NewShortenUrl(mockRepo, mockCodeGen)

			code, err := svc.CreateShortenLink(ctx, "https://google.com", 3600)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Len(t, code, service.DefaultCodeLength)
			}

			assert.Equal(t, tc.wantCode, code)

			if tc.verifyFunc != nil {
				tc.verifyFunc(t, mockRepo)
			}
		})
	}
}
