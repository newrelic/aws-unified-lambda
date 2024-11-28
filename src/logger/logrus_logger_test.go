package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

// TestWithLogLevel tests the WithLogLevel function of the LogrusLogger.
// It verifies that the log level is set correctly based on the input level.
func TestWithLogLevel(t *testing.T) {
	tests := []struct {
		name          string    // Name of the test case
		inputLevel    string    // Input log level
		expectedLevel log.Level // Expected log level
	}{
		{"ValidDebug", "debug", log.DebugLevel},
		{"ValidInfo", "info", log.InfoLevel},
		{"ValidWarn", "warn", log.WarnLevel},
		{"ValidError", "error", log.ErrorLevel},
		{"InvalidLevel", "invalid", log.InfoLevel},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewLogrusLogger(WithLogLevel(tc.inputLevel))
			if logger.GetLevel() != tc.expectedLevel {
				t.Errorf("WithLogLevel(%s) got %v, want %v", tc.inputLevel, logger.GetLevel(), tc.expectedLevel)
			}
		})
	}
}

// TestWithDebugLevel tests the WithDebugLevel function of the LogrusLogger.
// It verifies that the log level is set correctly based on the DebugEnabled environment variable.
func TestWithDebugLevel(t *testing.T) {
	tests := []struct {
		name          string    // Name of the test case
		debugEnabled  string    // DEBUG_ENABLED environment variable value
		expectedLevel log.Level // Expected log level
	}{
		{"DebugEnabled", "true", log.DebugLevel},
		{"DebugDisabled", "false", log.InfoLevel},
		{"DebugEmpty", "", log.InfoLevel},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(DebugEnabled, tc.debugEnabled)
			defer os.Unsetenv(DebugEnabled)

			logger := NewLogrusLogger(WithDebugLevel())
			if logger.GetLevel() != tc.expectedLevel {
				t.Errorf("WithDebugLevel() with DEBUG_ENABLED=%v got %v, want %v", tc.debugEnabled, logger.GetLevel(), tc.expectedLevel)
			}
		})
	}
}
