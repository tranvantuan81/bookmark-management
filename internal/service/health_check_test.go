package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tranvantuan81/bookmark-management/internal/config"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		cfg         *config.Config
		expectedRes Response
		expectedErr error
	}{
		{
			name: "success",
			cfg: &config.Config{
				ServiceName: "test",
				InstanceID:  "xxx",
			},
			expectedRes: Response{
				Message:     "OK",
				ServiceName: "test",
				InstanceID:  "xxx",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			healthCheckSev := NewHealthCheck(tc.cfg)

			res, err := healthCheckSev.CheckHealth()

			assert.Equal(t, tc.expectedRes, res)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
