package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/service/mocks"
)

// errSomething is a sentinel error used by the service mocks.
var errSomething = errors.New("something went wrong")

func TestGenPassHandler_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		setUpSvc   func(t *testing.T) *mocks.GenPass
		wantStatus int
		wantBody   string
	}{
		{
			name: "service returns a password",
			setUpSvc: func(t *testing.T) *mocks.GenPass {
				svcMock := mocks.NewGenPass(t)
				svcMock.On("GeneratePassword", passwordLength).Return("123456789012", nil).Once()
				return svcMock
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"password":"123456789012"}`,
		},
		{
			name: "service returns an error",
			setUpSvc: func(t *testing.T) *mocks.GenPass {
				svcMock := mocks.NewGenPass(t)
				svcMock.On("GeneratePassword", passwordLength).Return("", errSomething).Once()
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
			ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)

			handler := NewGenPass(tc.setUpSvc(t))

			handler.GeneratePassword(ctx)

			require.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}
