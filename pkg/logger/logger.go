package logger

import (
	"boiler_plate_be_golang/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// Init initializes the logger
func Init() {
	// Set log level based on environment
	var logLevel zerolog.Level
	switch config.App.App.Env {
	case "production":
		logLevel = zerolog.InfoLevel
	case "development":
		logLevel = zerolog.DebugLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	// Configure output format
	if config.App.App.Env == "development" {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}

	// Set as global logger
	log.Logger = Logger

	Logger.Info().
		Str("env", config.App.App.Env).
		Str("level", logLevel.String()).
		Msg("Logger initialized")
}

// Info logs info level message
func Info(msg string) *zerolog.Event {
	return Logger.Info()
}

// Debug logs debug level message
func Debug(msg string) *zerolog.Event {
	return Logger.Debug()
}

// Warn logs warning level message
func Warn(msg string) *zerolog.Event {
	return Logger.Warn()
}

// Error logs error level message
func Error(msg string) *zerolog.Event {
	return Logger.Error()
}

// Fatal logs fatal level message and exits
func Fatal(msg string) *zerolog.Event {
	return Logger.Fatal()
}

// WithContext returns logger with context fields
func WithContext() zerolog.Context {
	return Logger.With()
}
