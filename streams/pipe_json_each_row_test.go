package streams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of json.Marshaler for testing
type mockJsonEachRowMarshaler struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Err   error
}

func (m mockJsonEachRowMarshaler) MarshalJSON() ([]byte, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return json.Marshal(struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	})
}

func TestPipeJSONEachRowTransform_WriteTo(t *testing.T) {
	tests := []struct {
		name          string
		stream        ReadStream[mockJsonEachRowMarshaler]
		expectedJSON  string
		expectedError error
	}{
		{
			name: "valid stream with multiple records",
			stream: MemReader([]mockJsonEachRowMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedJSON:  `{"id":1,"name":"Alice","email":"alice@example.com"}` + "\n" + `{"id":2,"name":"Bob","email":"bob@example.com"}` + "\n",
			expectedError: nil,
		},
		{
			name:          "empty stream",
			stream:        MemReader([]mockJsonEachRowMarshaler{}, nil),
			expectedJSON:  ``,
			expectedError: nil,
		},
		{
			name: "stream with single record",
			stream: MemReader([]mockJsonEachRowMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
			}, nil),
			expectedJSON:  `{"id":1,"name":"Alice","email":"alice@example.com"}` + "\n",
			expectedError: nil,
		},
		{
			name:          "stream with error",
			stream:        MemReader([]mockJsonEachRowMarshaler{}, errors.New("stream error")),
			expectedJSON:  ``,
			expectedError: errors.New("stream error"),
		},
		{
			name:          "stream with error io.EOF",
			stream:        MemReader([]mockJsonEachRowMarshaler{}, io.EOF),
			expectedJSON:  ``,
			expectedError: nil,
		},
		{
			name: "marshaler with error",
			stream: MemReader([]mockJsonEachRowMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{
					ID:    2,
					Name:  "Error",
					Email: "error@example.com",
					Err:   errors.New("marshal error"),
				},
			}, nil),
			expectedJSON:  ``,
			expectedError: errors.New("marshal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			transform := JSONEachRowTransform(tt.stream)

			written, err := transform.WriteTo(buffer)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedJSON, buffer.String())
				assert.Equal(t, int64(len(tt.expectedJSON)), written)
			}
		})
	}
}

func TestPipeJSONEachRowTransform_WriteToFile(t *testing.T) {
	stream := MemReader([]mockJsonEachRowMarshaler{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}, nil)

	_, err := PipeJSONEachRow(stream, io.Discard)
	assert.NoError(t, err)
}

// ExampleJSONEachRowTransform demonstrates converting a stream to JSON lines format.
func ExampleJSONEachRowTransform() {
	// Define a simple structure for demonstration
	type LogEntry struct {
		Timestamp string `json:"timestamp"`
		Level     string `json:"level"`
		Message   string `json:"message"`
	}

	// Create a stream of log entries
	logs := []LogEntry{
		{Timestamp: "2025-06-28T10:00:00Z", Level: "INFO", Message: "Application started"},
		{Timestamp: "2025-06-28T10:01:00Z", Level: "WARN", Message: "High memory usage detected"},
		{Timestamp: "2025-06-28T10:02:00Z", Level: "ERROR", Message: "Database connection failed"},
	}
	stream := MemReader(logs, nil)

	// Transform to JSON lines format and write to stdout
	transform := JSONEachRowTransform(stream)
	transform.WriteTo(os.Stdout)

	// Output:
	// {"timestamp":"2025-06-28T10:00:00Z","level":"INFO","message":"Application started"}
	// {"timestamp":"2025-06-28T10:01:00Z","level":"WARN","message":"High memory usage detected"}
	// {"timestamp":"2025-06-28T10:02:00Z","level":"ERROR","message":"Database connection failed"}
}

// ExamplePipeJSONEachRow demonstrates using the PipeJSONEachRow convenience function.
func ExamplePipeJSONEachRow() {
	// Define a simple metric structure
	type Metric struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	}

	// Create a stream of metrics
	metrics := []Metric{
		{Name: "cpu_usage", Value: 85.5, Unit: "percent"},
		{Name: "memory_usage", Value: 1024, Unit: "MB"},
		{Name: "disk_usage", Value: 75.2, Unit: "percent"},
	}
	stream := MemReader(metrics, nil)

	// Use PipeJSONEachRow to write to stdout
	bytesWritten, err := PipeJSONEachRow(stream, os.Stdout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Bytes written: %d\n", bytesWritten)

	// Output:
	// {"name":"cpu_usage","value":85.5,"unit":"percent"}
	// {"name":"memory_usage","value":1024,"unit":"MB"}
	// {"name":"disk_usage","value":75.2,"unit":"percent"}
	// Bytes written: 152
}
