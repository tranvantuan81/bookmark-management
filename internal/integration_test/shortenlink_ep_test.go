package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/api"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestCreateShortenLinkEndpoints(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(engine api.Engine) *httptest.ResponseRecorder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://google.com", "exp":60}`)),
				)
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"message":`,
		},
		{
			name: "missing required field url returns bad request",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					strings.NewReader(`{"exp":3600}`),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":`,
		},
		{
			name: "malformed json returns bad request",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					strings.NewReader(`{"url":`),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"`,
		},
		{
			name: "wrong http method returns not found",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/shorten", nil)
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
		{
			name: "exp exceeds max (604800) returns bad request",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					strings.NewReader(`{"url":"https://google.com","exp":604801}`),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"`,
		},
		{
			name: "exp at max boundary (604800) is accepted",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					strings.NewReader(`{"url":"https://google.com","exp":604800}`),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"message":"Shorten URL generated successfully!"`,
		},
		{
			name: "missing required field exp returns bad request",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/v1/links/shorten",
					strings.NewReader(`{"url":"https://google.com"}`),
				)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRedis := redisPkg.InitMockRedis(t)

			testAPI := api.NewEngine(&config.Config{}, mockRedis)

			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestRedirectEndpoints(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(engine api.Engine) *httptest.ResponseRecorder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "valid code redirects to original url",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				code := createShortenLink(t, engine, "https://google.com", 3600)

				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/links/redirect/%s", code), nil)
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusMovedPermanently,
			expectedResponseBody: "https://google.com",
		},
		{
			name: "code not found returns internal server error",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/notexist", nil)
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `"error":"Internal Server Error"`,
		},
		{
			name: "wrong http method returns not found",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/v1/links/redirect/abcdefg", nil)
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testAPI := api.NewEngine(&config.Config{}, redisPkg.InitMockRedis(t))

			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

// createShortenLink is a test helper that calls the shorten endpoint and returns the generated code.
func createShortenLink(t *testing.T, engine api.Engine, url string, exp int64) string {
	t.Helper()

	body := fmt.Sprintf(`{"url":%q,"exp":%d}`, url, exp)
	req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var res struct {
		Code string `json:"code"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&res))
	require.NotEmpty(t, res.Code)

	return res.Code
}
