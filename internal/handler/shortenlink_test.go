package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

func TestShortenLinkHandler_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		requestBody string
		setUpSvc    func(t *testing.T) *mocks.ShortenUrl
		wantStatus  int
		wantBody    string
	}{
		{
			name:        "valid payload returns the generated code",
			requestBody: `{"url":"https://google.com","exp":3600}`,
			setUpSvc: func(t *testing.T) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("CreateShortenLink", mock.Anything, "https://google.com", 3600).
					Return("abcdefg", nil).Once()
				return svcMock
			},
			wantStatus: http.StatusOK,
			wantBody: `{
				"code": "abcdefg",
				"message": "Shorten URL generated successfully!"
			}`,
		},
		{
			name:        "malformed payload is rejected before the service is called",
			requestBody: `{"url":`,
			setUpSvc: func(t *testing.T) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.AssertNotCalled(t, "CreateShortenLink", mock.Anything, mock.Anything, mock.Anything)
				return svcMock
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Invalid input"}`,
		},
		{
			name:        "service returns an error",
			requestBody: `{"url":"https://google.com","exp":3600}`,
			setUpSvc: func(t *testing.T) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("CreateShortenLink", mock.Anything, "https://google.com", 3600).
					Return("", errSomething).Once()
				return svcMock
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/v1/links/shorten",
				strings.NewReader(tc.requestBody),
			)
			ctx.Request.Header.Set("Content-Type", "application/json")

			handler := NewShortenLink(tc.setUpSvc(t))

			handler.CreateShortenLink(ctx)

			require.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}
