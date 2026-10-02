package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/api"
)

func TestGenPassEndPoints(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		call       func(engine api.Engine) *httptest.ResponseRecorder
		wantStatus int
	}{
		{
			name: "GET returns a generated password",
			call: func(engine api.Engine) *httptest.ResponseRecorder {
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/genpass", nil))
				return rec
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "unsupported method is not routed",
			call: func(engine api.Engine) *httptest.ResponseRecorder {
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/genpass", nil))
				return rec
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine, _ := setupTestAPI(t)

			rec := tc.call(engine)

			require.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantStatus != http.StatusOK {
				return
			}

			// the password is random, so assert on its shape instead of its value
			var res struct {
				Password string `json:"password"`
			}
			decodeBody(t, rec, &res)

			assert.Len(t, res.Password, expectedPasswordLength)
		})
	}
}
