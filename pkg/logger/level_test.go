package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		input         string
		expectedLevel zerolog.Level
	}{
		{
			name:          "debug level",
			input:         "debug",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			name:          "info level",
			input:         "info",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			name:          "warn level",
			input:         "warn",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			name:          "error level",
			input:         "error",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			name:          "invalid level falls back to info",
			input:         "invalid",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			name:          "empty string is parsed as NoLevel",
			input:         "",
			expectedLevel: zerolog.NoLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetLogLevel(tc.input)

			assert.Equal(t, tc.expectedLevel, zerolog.GlobalLevel())
		})
	}
}
