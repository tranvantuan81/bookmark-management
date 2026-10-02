package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

func TestGeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		length     int
		wantLength int
		wantErr    error
	}{
		{
			name:       "default length",
			length:     service.DefaultCodeLength,
			wantLength: service.DefaultCodeLength,
			wantErr:    nil,
		},
		{
			name:       "custom length",
			length:     20,
			wantLength: 20,
			wantErr:    nil,
		},
		{
			name:       "minimum length",
			length:     1,
			wantLength: 1,
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewGenPass()

			password, err := svc.GeneratePassword(tc.length)

			require.NoError(t, err)
			assert.Len(t, password, tc.wantLength)

			for _, char := range password {
				assert.Contains(t, charset, string(char))
			}
		})
	}
}

// charset mirrors the alphabet used by the implementation; keeping the copy here
// makes the assertion meaningful without exporting internals from the package.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
