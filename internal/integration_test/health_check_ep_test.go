package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/api"
)

func TestHealthCheckEndPoints(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		call       func(engine api.Engine) *httptest.ResponseRecorder
		wantStatus int
		wantBody   string
	}{
		{
			name: "GET returns the running instance",
			call: func(engine api.Engine) *httptest.ResponseRecorder {
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health-check", nil))
				return rec
			},
			wantStatus: http.StatusOK,
			wantBody: `{
				"message": "OK",
				"service_name": "bookmark-management",
				"instance_id": "0f1092e7-38ed-4701-b500-c4697c9dc122"
			}`,
		},
		{
			name: "unsupported method is not routed",
			call: func(engine api.Engine) *httptest.ResponseRecorder {
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/health-check", nil))
				return rec
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `404 page not found`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine, _ := setupTestAPI(t)

			rec := tc.call(engine)

			require.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantStatus == http.StatusNotFound {
				// gin's 404 replies with plain text, not JSON
				assert.Equal(t, tc.wantBody, rec.Body.String())
				return
			}
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}
