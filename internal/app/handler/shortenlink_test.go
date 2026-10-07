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
	"github.com/tranvantuan81/bookmark-management/internal/app/service/mocks"
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
		{
			name: "fail case - url field is not a valid URL",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"not-a-url","exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid input"}`,
		},
		{
			name: "fail case - exp exceeds max allowed value (604800)",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://google.com","exp":604801}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid input"}`,
		},
		{
			name: "fail case - missing exp field",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://google.com"}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid input"}`,
		},
		{
			name: "fail case - missing url field",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"exp":60}`)),
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid input"}`,
		},
		{
			name: "fail case - empty request body",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
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
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("GetURLFromCode", ctx, "1234567").
					Return("https://google.com", nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusMovedPermanently,
			expectedUrl:    "https://google.com",
		},
		{
			name: "service fails",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("GetURLFromCode", ctx, "1234567").
					Return("", testErr).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedUrl:    "",
		},
		{
			name: "empty code returns bad request",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				ctx.Params = gin.Params{{Key: "code", Value: ""}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedUrl:    "",
		},
		{
			name: "code not found in storage returns internal server error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/notexist", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "notexist"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("GetURLFromCode", ctx, "notexist").
					Return("", ErrNotFound).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedUrl:    "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			testHandler := NewShortenLink(tc.setupMockService(ctx))

			testHandler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedUrl, rec.Header().Get("Location"))
		})
	}
}
