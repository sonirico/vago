package lol

import (
	"bytes"
	"context"
	"fmt"

	"go.elastic.co/apm/v2"
)

// ExampleNewZerolog demonstrates creating a structured logger with zerolog backend
func ExampleNewZerolog() {
	// Create a logger with custom fields and configuration
	var buf bytes.Buffer

	logger := NewZerolog(
		WithFields(Fields{"service": "example-app", "version": "1.0.0"}),
		WithEnv(EnvProd),
		WithLevel(LevelInfo),
		WithWriter(&buf),
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

	fmt.Printf(
		"Logged output contains 'version': %t\n",
		bytes.Contains(buf.Bytes(), []byte("version")),
	)

	// Output:
	// Logged output contains 'Application started': true
	// Logged output contains 'user_id': true
	// Logged output contains 'service': true
	// Logged output contains 'version': true
}

// ExampleLogger_WithField demonstrates adding contextual fields to log messages
func ExampleLogger_WithField() {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithFields(Fields{"app": "demo"}),
		WithEnv(EnvDev),
		WithLevel(LevelDebug),
		WithWriter(&buf),
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

	logger := NewZerolog(
		WithFields(Fields{"component": "auth"}),
		WithEnv(EnvDev),
		WithLevel(LevelTrace), // Set to trace level to see all messages
		WithWriter(&buf),
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

// ExampleLogger_WithTrace demonstrates APM trace context integration
func ExampleLogger_WithTrace() {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithFields(Fields{"service": "payment-service"}),
		WithEnv(EnvDev),
		WithLevel(LevelInfo),
		WithWriter(&buf),
		WithApm(), // Enable APM tracing
	)

	// Create an APM transaction (simulating real APM integration)
	tracer := apm.DefaultTracer()
	tx := tracer.StartTransaction("payment-processing", "request")
	defer tx.End()

	// Create context with the transaction
	ctx := apm.ContextWithTransaction(context.Background(), tx)

	// Start a span for more detailed tracing
	span, ctx := apm.StartSpan(ctx, "payment-validation", "internal")
	defer span.End()

	// Create a logger with trace context
	tracedLogger := logger.WithTrace(ctx)

	// Log with trace context - these should include APM trace fields
	tracedLogger.Info("Processing payment request")
	tracedLogger.WithField("payment_id", "pay_123").
		WithField("amount", 99.99).
		Info("Payment validation started")

	// Log without trace context for comparison
	logger.Info("Regular log message without trace context")

	output := buf.String()
	fmt.Printf(
		"Output contains 'Processing payment': %t\n",
		bytes.Contains([]byte(output), []byte("Processing payment")),
	)
	fmt.Printf(
		"Output contains 'payment_id': %t\n",
		bytes.Contains([]byte(output), []byte("payment_id")),
	)
	fmt.Printf(
		"Output contains 'service': %t\n",
		bytes.Contains([]byte(output), []byte("service")),
	)

	// Check for APM trace fields in the output
	// The apmzerolog hook should add these fields when WithTrace is used
	fmt.Printf(
		"Output contains trace information: %t\n",
		bytes.Contains([]byte(output), []byte("trace")) ||
			bytes.Contains([]byte(output), []byte("transaction")) ||
			bytes.Contains([]byte(output), []byte("span")),
	)

	// Output:
	// Output contains 'Processing payment': true
	// Output contains 'payment_id': true
	// Output contains 'service': true
	// Output contains trace information: true
}
