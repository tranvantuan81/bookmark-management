package logger

import "github.com/rs/zerolog"

// SetLogLevel sets the global logging level based on the provided string. Defaults to info level if parsing fails.
func SetLogLevel(levelStr string) {
	level := zerolog.NoLevel
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
