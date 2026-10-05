package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

var ErrNotFound = errors.New("url not found")

func TestShortenLink_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShortenUrl

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://google.com", "exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("CreateShortenLink", ctx, "https://google.com", int64(60)).
					Return("1234567", nil)
				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"1234567","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "fail case - service failed",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://google.com", "exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("CreateShortenLink", ctx, "https://google.com", int64(60)).
					Return("", testErr)
				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"Internal Server Error"}`,
		},
		{
			name: "fail case - bad input",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"not a link", "exp":60'}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				return serviceMock
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid input"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewShortenLink(mockSvc)

			testHandler.CreateShortenLink(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}

func TestShortenLink_Redirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShortenUrl
		expectedStatus   int
		expectedUrl      string
	}{
		{
			name: "normal",
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				serviceMock := mocks.NewShortenUrl(t)
				serviceMock.On("GetURLFromCode", ctx, "1234567").Return("https://google.com", nil)
				return serviceMock
			},
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			expectedStatus: http.StatusMovedPermanently,
			expectedUrl:    "https://google.com",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewShortenLink(mockSvc)

			testHandler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedUrl, rec.Header().Get("Location"))

		})
	}
}
