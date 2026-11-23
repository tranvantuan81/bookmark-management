package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedRes Response
		expectedErr error
	}{
		{
			name: "test case 1",
			expectedRes: Response{
				Message:     "OK",
				ServiceName: "bookmark_service",
				InstanceID:  "xxx",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appName := "bookmark_service"
			instanceID := "xxx"

			svc := NewHealthCheck(appName, instanceID)

			result, err := svc.Check()

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedRes, result)

		})
	}
}
