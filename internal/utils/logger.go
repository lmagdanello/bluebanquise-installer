package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var Logger *slog.Logger

// InitLogger initializes the logger for BlueBanquise installer
func InitLogger() error {
	// Try to use LOG_DIR environment variable first
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "/var/log/bluebanquise"
	}

	// Try to create log directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If we can't create /var/log/bluebanquise, try a temporary directory
		if logDir == "/var/log/bluebanquise" {
			logDir = os.TempDir()
		} else {
			return err
		}
	}

	// Create log file
	logFile := filepath.Join(logDir, "bluebanquise-installer.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Create multi-writer for both file and console
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Create logger with multi-writer
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Logger = slog.New(handler)

	// Set as default logger
	slog.SetDefault(Logger)

	// Log startup
	Logger.Info("BlueBanquise installer started",
		"version", "3.2.0",
		"log_file", logFile)

	return nil
}

// InitTestLogger initializes the logger for testing
func InitTestLogger() {
	// Create logger that writes to io.Discard for tests
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// LogCommand logs a command execution
func LogCommand(command string, args ...string) {
	Logger.Info("Executing command",
		"command", command,
		"args", args)
}

// LogError logs an error with context
func LogError(msg string, err error, context ...any) {
	Logger.Error(msg, append([]any{"error", err}, context...)...)
}

// LogInfo logs an info message
func LogInfo(msg string, context ...any) {
	Logger.Info(msg, context...)
}

// LogWarning logs a warning message
func LogWarning(msg string, context ...any) {
	Logger.Warn(msg, context...)
}
