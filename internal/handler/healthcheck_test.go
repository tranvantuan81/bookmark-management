package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/service"
	healthcheck_svc_mocks "github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

func TestHealthCheck_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		setupMockSvc func(ctx context.Context) *healthcheck_svc_mocks.HealthCheck

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case",
			setupMockSvc: func(ctx context.Context) *healthcheck_svc_mocks.HealthCheck {
				mockSvc := healthcheck_svc_mocks.NewHealthCheck(t)
				mockSvc.On("Check").Return(service.Response{
					Message:     "OK",
					ServiceName: "bookmark_service",
					InstanceID:  "xxx",
				}, nil)
				return mockSvc
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"OK","service_name":"bookmark_service","instance_id":"xxx"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			mockSvc := tc.setupMockSvc(ginCtx)

			healthcheckHandler := NewHealthCheck(mockSvc)
			healthcheckHandler.HealthCheck(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
