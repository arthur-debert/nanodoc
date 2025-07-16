/*
 * This is the logging setup for ttxt:
 * - Uses zerolog for logging
 * - Uses console writer for human-readable output
 * - Maps command line flags to log levels:
 *   - -v: Info
 *   - -vv: Debug
 *   - -vvv: Trace
 * - Level Guideline:
 *   - Warn: is the default level.
 *   - Info: short messages signaling entry / exit of relevant execution points.
 *     may include some details about the execution, but not significant
 *     amount of data
 *   - Debug: detailed messages for debugging, which include some data.
 *   - Trace: very detailed messages for debugging, which include quiet a bit
 *     of data.
 *   - The core idea is :
 *     what's going on? info.
 *     what is this returning that? debug.
 *     what is happening to this variable? trace.
 */
package logging

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogLevel represents the available log levels
type LogLevel int

const (
	WarnLevel LogLevel = iota
	InfoLevel
	DebugLevel
	TraceLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	default:
		return "warn"
	}
}

// ToZerologLevel converts our LogLevel to zerolog.Level
func (l LogLevel) ToZerologLevel() zerolog.Level {
	switch l {
	case TraceLevel:
		return zerolog.TraceLevel
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	default:
		return zerolog.WarnLevel
	}
}

// Config holds the logging configuration
type Config struct {
	Level  LogLevel
	Writer io.Writer
}

// DefaultConfig returns the default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:  WarnLevel,
		Writer: os.Stderr,
	}
}

// Setup configures the global zerolog logger with the given configuration
func Setup(config *Config) {
	// Set global log level
	zerolog.SetGlobalLevel(config.Level.ToZerologLevel())

	// Configure the global logger with console writer for human-readable output
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        config.Writer,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	// Enable stack trace for error level
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return err.Error()
	}
}

// ParseVerbosityFlags converts CLI verbosity flags to LogLevel
func ParseVerbosityFlags(verbose, v, vv, vvv bool) LogLevel {
	// Priority order: vvv > vv > v > verbose
	if vvv {
		return TraceLevel
	}
	if vv {
		return DebugLevel
	}
	if v || verbose {
		return InfoLevel
	}
	return WarnLevel
}

// NewTestLogger creates a logger for testing that integrates with Go's testing framework
func NewTestLogger(t *testing.T) zerolog.Logger {
	level := InfoLevel
	if testing.Verbose() {
		level = TraceLevel
	}

	// Use zerolog's test writer to integrate with t.Log()
	testWriter := zerolog.NewTestWriter(t)
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        testWriter,
		TimeFormat: "15:04:05.000",
	}).Level(level.ToZerologLevel()).With().Timestamp().Logger()
}
