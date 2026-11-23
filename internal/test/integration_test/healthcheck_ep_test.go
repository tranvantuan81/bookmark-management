package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/api"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupTestRequest func() (*httptest.ResponseRecorder, *http.Request)

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case",
			setupTestRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health-check", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":`,
		},

		{
			name: "wrong endpoint",
			setupTestRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health-check123", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg, _ := api.NewConfig()

			testAPIEngine := api.NewEngine(cfg)
			rec, req := tc.setupTestRequest()
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)

		})
	}
}
