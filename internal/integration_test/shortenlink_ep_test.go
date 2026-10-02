package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/api"
)

func TestShortenLinkEndPoints(t *testing.T) {
	t.Parallel()

	const (
		originalURL = "https://google.com"
		ttlSeconds  = 3600
	)

	testCases := []struct {
		name string
		body string
		call func(engine api.Engine, body string) *httptest.ResponseRecorder

		wantStatus int
		wantBody   string
		wantStored bool
	}{
		{
			name:       "valid payload is shortened and stored",
			body:       `{"url":"https://google.com","exp":3600}`,
			call:       postShortenLink,
			wantStatus: http.StatusOK,
			wantStored: true,
		},
		{
			name:       "malformed payload is rejected and nothing is stored",
			body:       `{"url":`,
			call:       postShortenLink,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Invalid input"}`,
		},
		{
			name: "unsupported method is not routed",
			call: func(engine api.Engine, _ string) *httptest.ResponseRecorder {
				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/links/shorten", nil))
				return rec
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine, redisClient := setupTestAPI(t)

			rec := tc.call(engine, tc.body)

			require.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantStatus == http.StatusNotFound {
				// gin replies 404 with plain text, not JSON
				assert.Equal(t, "404 page not found", rec.Body.String())
				return
			}

			if tc.wantBody != "" {
				assert.JSONEq(t, tc.wantBody, rec.Body.String())
			}

			ctx := t.Context()
			keys, err := redisClient.Keys(ctx, "*").Result()
			require.NoError(t, err)

			if !tc.wantStored {
				assert.Empty(t, keys, "nothing must be written for a rejected request")
				return
			}

			// the generated code is random, so assert on its shape and on what was stored
			var res struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}
			decodeBody(t, rec, &res)

			assert.Len(t, res.Code, expectedCodeLength)
			assert.Equal(t, "Shorten URL generated successfully!", res.Message)

			storedURL, err := redisClient.Get(ctx, res.Code).Result()
			assert.NoError(t, err)
			assert.Equal(t, originalURL, storedURL)

			ttlLeft, err := redisClient.TTL(ctx, res.Code).Result()
			assert.NoError(t, err)
			assert.Equal(t, ttlSeconds*time.Second, ttlLeft)
		})
	}
}

func postShortenLink(engine api.Engine, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(rec, req)

	return rec
}
