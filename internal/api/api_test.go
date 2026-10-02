package api_test

import (
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

func TestEngine_ServeHTTP(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "health check route is registered",
			method:     http.MethodGet,
			path:       "/health-check",
			wantStatus: http.StatusOK,
		},
		{
			name:       "genpass route is registered",
			method:     http.MethodGet,
			path:       "/genpass",
			wantStatus: http.StatusOK,
		},
		{
			name:       "shorten route is registered",
			method:     http.MethodPost,
			path:       "/v1/links/shorten",
			body:       `{"url":"https://google.com","exp":3600}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "unknown route returns not found",
			method:     http.MethodGet,
			path:       "/not-exist",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "wrong method returns not found",
			method:     http.MethodDelete,
			path:       "/health-check",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine := api.NewEngine(&config.Config{
				AppPort:     "8080",
				ServiceName: "bookmark-management",
				InstanceID:  "test-instance",
			}, redisPkg.InitMockRedis(t))

			var req *http.Request
			if tc.body == "" {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			} else {
				req = httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()

			engine.ServeHTTP(rec, req)

			require.Equal(t, tc.wantStatus, rec.Code)
			assert.NotEmpty(t, rec.Body.String())
		})
	}
}
