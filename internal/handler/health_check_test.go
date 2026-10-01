package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/service"
	"github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		setUpSvc   func(t *testing.T) *mocks.HealthCheck
		wantStatus int
		wantBody   string
	}{
		{
			name: "service reports a healthy instance",
			setUpSvc: func(t *testing.T) *mocks.HealthCheck {
				svcMock := mocks.NewHealthCheck(t)
				svcMock.On("HealthCheck", mock.Anything).Return(service.Response{
					Message:     "OK",
					ServiceName: "bookmark-management",
					InstanceID:  "0f1092e7-38ed-4701-b500-c4697c9dc122",
				}, nil).Once()
				return svcMock
			},
			wantStatus: http.StatusOK,
			wantBody: `{
				"message": "OK",
				"service_name": "bookmark-management",
				"instance_id": "0f1092e7-38ed-4701-b500-c4697c9dc122"
			}`,
		},
		{
			name: "service returns an error",
			setUpSvc: func(t *testing.T) *mocks.HealthCheck {
				svcMock := mocks.NewHealthCheck(t)
				svcMock.On("HealthCheck", mock.Anything).Return(service.Response{}, errSomething).Once()
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
			ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

			handler := NewHealthCheck(tc.setUpSvc(t))

			handler.HealthCheck(ctx)

			require.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}
