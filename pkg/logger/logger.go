package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// LoggerContext wraps zerolog.Logger to provide logrus-style methods
type LoggerContext struct {
	logger zerolog.Logger
}

// Errorf logs error message with format
func (l *LoggerContext) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Warnf logs warning message with format
func (l *LoggerContext) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Infof logs info message with format
func (l *LoggerContext) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Debugf logs debug message with format
func (l *LoggerContext) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Error returns error event
func (l *LoggerContext) Error(msg string) *zerolog.Event {
	return l.logger.Error()
}

// Warn returns warning event
func (l *LoggerContext) Warn(msg string) *zerolog.Event {
	return l.logger.Warn()
}

// Info returns info event
func (l *LoggerContext) Info(msg string) *zerolog.Event {
	return l.logger.Info()
}

// Debug returns debug event
func (l *LoggerContext) Debug(msg string) *zerolog.Event {
	return l.logger.Debug()
}

// Init initializes the logger
func Init(env string) {
	// Set log level based on environment
	var logLevel zerolog.Level
	switch env {
	case "production":
		logLevel = zerolog.InfoLevel
	case "development":
		logLevel = zerolog.DebugLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	// Configure output format
	if env == "development" {
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
		Str("env", env).
		Str("level", logLevel.String()).
		Msg("Logger initialized")
}

// GetLoggerContext returns logger with context information (Localoka_V2 pattern)
func GetLoggerContext(ctx context.Context, functionName string) *LoggerContext {
	logger := Logger.With().
		Str("function", functionName).
		Logger()
	
	// Add request ID if available in context
	if requestID := ctx.Value("request_id"); requestID != nil {
		if reqID, ok := requestID.(string); ok {
			logger = logger.With().Str("request_id", reqID).Logger()
		}
	}
	
	// Add user ID if available in context
	if userID := ctx.Value("user_id"); userID != nil {
		if uid, ok := userID.(string); ok {
			logger = logger.With().Str("user_id", uid).Logger()
		}
	}
	
	return &LoggerContext{logger: logger}
}

// GetLogger returns logger with custom context (Localoka_V2 pattern)
func GetLogger(component, action string) *LoggerContext {
	logger := Logger.With().
		Str("component", component).
		Str("action", action).
		Logger()
	return &LoggerContext{logger: logger}
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
