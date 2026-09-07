package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/tranvantuan81/bookmark-management/internal/service"
	"github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		mockResponse     service.Response
		mockError        error
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			mockResponse: service.Response{
				Message:     "OK",
				ServiceName: "bookmark-management",
				InstanceID:  "0f1092e7-38ed-4701-b500-c4697c9dc122",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedResponse: `{
				"message": "OK",
				"service_name": "bookmark-management",
				"instance_id": "0f1092e7-38ed-4701-b500-c4697c9dc122"
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = httptest.NewRequest(
				http.MethodGet,
				"/health-check",
				nil,
			)

			mockService := mocks.NewHealthCheck(t)

			mockService.
				On("CheckHealth").
				Return(tc.mockResponse, tc.mockError)

			handler := NewHealthCheck(mockService)

			handler.HealthCheck(ctx)
			assert.Equal(t, tc.expectedStatus, w.Code)

			assert.JSONEq(
				t,
				tc.expectedResponse,
				w.Body.String(),
			)
		})
	}
}
