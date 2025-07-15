package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"zerofy.pro/rbac-collector/src/constants"
)

// "console" -> pretty, human-readable output for development.
// "json"    -> standard JSON output for production.
func New() zerolog.Logger {
	logFormat := strings.ToLower(os.Getenv(constants.EnvLogFormat))

	switch logFormat {
	case constants.LogFormatConsole:
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		return zerolog.New(output).With().Timestamp().Logger()
	case constants.LogFormatJSON:
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	default:
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
