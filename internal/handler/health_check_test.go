package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/service"
	"github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

var response = service.Response{
	Message:     "OK",
	ServiceName: "bookmark-management",
	InstanceID:  "",
}

var testErr = "test error"

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.HealthCheck

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("HealthCheck", ctx).Return(response, nil)
				return serviceMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"OK","service_name":"bookmark-management","instance_id":""}`,
		},

		{
			name: "service fail",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)
				serviceMock.On("HealthCheck", ctx).Return(service.Response{}, testErr)
				return serviceMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"Internal Server Error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewHealthCheck(mockSvc)

			testHandler.HealthCheck(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}
