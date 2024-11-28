// Package logger provides a Logrus-based logger implementation for unified logging.
package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// DebugEnabled is the environment variable name used to enable debug level logging.
const (
	DebugEnabled = "DEBUG_ENABLED"
)

// ConfigOption is a function type used to configure the logger.
type ConfigOption func(*log.Logger)

// NewLogrusLogger creates a new instance of logrus.Logger with the provided configuration options.
func NewLogrusLogger(opts ...ConfigOption) *log.Logger {
	l := log.New()
	for _, fn := range opts {
		if nil != fn {
			fn(l)
		}
	}

	return l
}

// WithLogLevel is a configuration option that sets the log level of the logger.
func WithLogLevel(level string) ConfigOption {
	return func(l *log.Logger) {
		parsedLevel, err := log.ParseLevel(level)
		if err != nil {
			l.Errorf("Invalid log level '%s'. Using default 'info' level.", level)
			parsedLevel = log.InfoLevel
		}
		l.SetLevel(parsedLevel)
	}
}

// WithDebugLevel is a configuration option that sets the log level to debug if the DebugEnabled environment variable is set to "true", otherwise sets it to info.
func WithDebugLevel() ConfigOption {
	if os.Getenv(DebugEnabled) == "true" {
		return WithLogLevel("debug")
	}
	return WithLogLevel("info")
}
