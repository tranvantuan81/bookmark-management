package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/api"
	"github.com/tranvantuan81/bookmark-management/internal/config"
)

func TestGenPassEndPoints(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/genpass", nil)
				resRecorder := httptest.NewRecorder()
				engine.ServeHTTP(resRecorder, req)
				return resRecorder
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"password":`,
		},
		{
			name: "wrong endpoint",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/genpass", nil)
				resRecorder := httptest.NewRecorder()
				engine.ServeHTTP(resRecorder, req)
				return resRecorder
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testAPI := api.NewEngine(&config.Config{})
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}
