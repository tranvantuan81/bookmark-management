package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func TestNewEngine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		cfg    *config.Config
		verify func(t *testing.T, e Engine)
	}{
		{
			name: "engine is created with routes registered",
			cfg: &config.Config{
				ServiceName: "bookmark-management",
				InstanceID:  "test-id",
			},
			verify: func(t *testing.T, e Engine) {
				require.NotNil(t, e)

				// health-check route
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				e.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusOK, rec.Code)

				// shorten route exists (bad input → 400, not 404)
				rec2 := httptest.NewRecorder()
				req2 := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", nil)
				req2.Header.Set("Content-Type", "application/json")
				e.ServeHTTP(rec2, req2)
				assert.Equal(t, http.StatusBadRequest, rec2.Code)

				// redirect route exists (unknown code → 500, not 404)
				rec3 := httptest.NewRecorder()
				req3 := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/unknowncode", nil)
				e.ServeHTTP(rec3, req3)
				assert.Equal(t, http.StatusInternalServerError, rec3.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRedis := redisPkg.InitMockRedis(t)
			e := NewEngine(tc.cfg, mockRedis)

			tc.verify(t, e)
		})
	}
}
