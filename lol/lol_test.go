package lol

import (
	"bytes"
	"fmt"
)

// ExampleNewZerologLogger demonstrates creating a structured logger with zerolog backend
func ExampleNewZerologLogger() {
	// Create a logger with custom fields and configuration
	var buf bytes.Buffer

	logger := NewZerologLogger(
		Fields{"service": "example-app", "version": "1.0.0"},
		"production",
		"info",
		&buf,
		APMConfig{Enabled: false}, // Disable APM for this example
	)

	// Log some messages
	logger.Info("Application started successfully")
	logger.WithField("user_id", 123).WithField("action", "login").Info("User logged in")
	logger.Warn("This is a warning message")

	fmt.Printf(
		"Logged output contains 'Application started': %t\n",
		bytes.Contains(buf.Bytes(), []byte("Application started")),
	)
	fmt.Printf(
		"Logged output contains 'user_id': %t\n",
		bytes.Contains(buf.Bytes(), []byte("user_id")),
	)
	fmt.Printf(
		"Logged output contains 'service': %t\n",
		bytes.Contains(buf.Bytes(), []byte("service")),
	)

	// Output:
	// Logged output contains 'Application started': true
	// Logged output contains 'user_id': true
	// Logged output contains 'service': true
}

// ExampleNewTest demonstrates creating a test logger that doesn't output anything
func ExampleNewTest() {
	// Create a test logger for unit tests
	testLogger := NewTest()

	// These messages won't be printed to stdout/stderr
	testLogger.Info("This is a test message")
	testLogger.Error("This error won't be shown")
	testLogger.WithField("test_field", "test_value").Debug("Debug message")

	fmt.Println("Test logger created successfully")
	fmt.Println("Messages logged silently for testing")

	// Output:
	// Test logger created successfully
	// Messages logged silently for testing
}

// ExampleLogger_WithField demonstrates adding contextual fields to log messages
func ExampleLogger_WithField() {
	var buf bytes.Buffer

	logger := NewZerologLogger(
		Fields{"app": "demo"},
		"development",
		"debug",
		&buf,
		APMConfig{Enabled: false},
	)

	// Chain multiple fields
	enrichedLogger := logger.WithField("request_id", "req-123").
		WithField("user_agent", "test-client")
	enrichedLogger.Info("Processing request")

	// Add more context
	enrichedLogger.WithField("duration_ms", 45).Info("Request completed")

	output := buf.String()
	fmt.Printf(
		"Output contains 'request_id': %t\n",
		bytes.Contains([]byte(output), []byte("request_id")),
	)
	fmt.Printf(
		"Output contains 'user_agent': %t\n",
		bytes.Contains([]byte(output), []byte("user_agent")),
	)
	fmt.Printf(
		"Output contains 'duration_ms': %t\n",
		bytes.Contains([]byte(output), []byte("duration_ms")),
	)

	// Output:
	// Output contains 'request_id': true
	// Output contains 'user_agent': true
	// Output contains 'duration_ms': true
}

// ExampleLogger_LogLevels demonstrates different log levels
func ExampleLogger_LogLevels() {
	var buf bytes.Buffer

	logger := NewZerologLogger(
		Fields{"component": "auth"},
		"development",
		"trace", // Set to trace level to see all messages
		&buf,
		APMConfig{Enabled: false},
	)

	// Log at different levels
	logger.Trace("Entering authentication function")
	logger.Debug("Validating user credentials")
	logger.Info("User authentication successful")
	logger.Warn("Rate limit approaching")
	logger.Error("Authentication failed")

	output := buf.String()
	fmt.Printf(
		"Contains trace message: %t\n",
		bytes.Contains([]byte(output), []byte("Entering authentication")),
	)
	fmt.Printf(
		"Contains debug message: %t\n",
		bytes.Contains([]byte(output), []byte("Validating user")),
	)
	fmt.Printf(
		"Contains info message: %t\n",
		bytes.Contains([]byte(output), []byte("authentication successful")),
	)
	fmt.Printf("Contains warn message: %t\n", bytes.Contains([]byte(output), []byte("Rate limit")))
	fmt.Printf(
		"Contains error message: %t\n",
		bytes.Contains([]byte(output), []byte("Authentication failed")),
	)

	// Output:
	// Contains trace message: true
	// Contains debug message: true
	// Contains info message: true
	// Contains warn message: true
	// Contains error message: true
}

// ExampleParseLevel demonstrates parsing log levels from strings
func ExampleParseLevel() {
	// Parse different log levels
	levels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}

	for _, levelStr := range levels {
		level := ParseLevel(levelStr)
		fmt.Printf("Level '%s' parsed as: %d\n", levelStr, level)
	}

	// Output:
	// Level 'trace' parsed as: 6
	// Level 'debug' parsed as: 5
	// Level 'info' parsed as: 4
	// Level 'warn' parsed as: 3
	// Level 'error' parsed as: 2
	// Level 'fatal' parsed as: 1
	// Level 'panic' parsed as: 0
}

// ExampleParseEnv demonstrates parsing environment strings
func ExampleParseEnv() {
	// Parse different environments
	envs := []string{"test", "local", "development", "staging", "production"}

	for _, envStr := range envs {
		env := ParseEnv(envStr)
		fmt.Printf("Environment '%s' parsed as: %d\n", envStr, env)
	}

	// Output:
	// Environment 'test' parsed as: 0
	// Environment 'local' parsed as: 1
	// Environment 'development' parsed as: 2
	// Environment 'staging' parsed as: 3
	// Environment 'production' parsed as: 4
}
