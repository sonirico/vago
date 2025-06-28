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
type mockJsonMarshaler struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Err   error
}

func (m mockJsonMarshaler) MarshalJSON() ([]byte, error) {
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

func TestPipeJSONTransform_WriteTo(t *testing.T) {
	tests := []struct {
		name          string
		stream        ReadStream[mockJsonMarshaler]
		expectedJSON  string
		expectedError error
	}{
		{
			name: "valid stream with multiple records",
			stream: MemReader([]mockJsonMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedJSON:  `[{"id":1,"name":"Alice","email":"alice@example.com"},{"id":2,"name":"Bob","email":"bob@example.com"}]`,
			expectedError: nil,
		},
		{
			name:          "empty stream",
			stream:        MemReader([]mockJsonMarshaler{}, nil),
			expectedJSON:  `[]`,
			expectedError: nil,
		},
		{
			name: "stream with single record",
			stream: MemReader([]mockJsonMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
			}, nil),
			expectedJSON:  `[{"id":1,"name":"Alice","email":"alice@example.com"}]`,
			expectedError: nil,
		},
		{
			name:          "stream with error",
			stream:        MemReader([]mockJsonMarshaler{}, errors.New("stream error")),
			expectedJSON:  ``,
			expectedError: errors.New("stream error"),
		},
		{
			name:          "stream with error io.EOF",
			stream:        MemReader([]mockJsonMarshaler{}, io.EOF),
			expectedJSON:  `[]`,
			expectedError: nil,
		},
		{
			name: "marshaler with error",
			stream: MemReader([]mockJsonMarshaler{
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
			transform := JSONTransform(tt.stream)

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

func TestPipeJSONTransform_WriteToFile(t *testing.T) {
	stream := MemReader([]mockJsonMarshaler{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}, nil)

	_, err := JSONTransform(stream).WriteTo(os.Stdout)
	assert.NoError(t, err)
}

// ExampleJSONTransform demonstrates converting a stream to JSON array format.
func ExampleJSONTransform() {
	// Define a simple structure for demonstration
	type Person struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Create a stream from a slice of persons
	people := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	stream := MemReader(people, nil)

	// Transform to JSON and write to stdout
	transform := JSONTransform(stream)
	transform.WriteTo(os.Stdout)

	// Output:
	// [{"id":1,"name":"Alice"},{"id":2,"name":"Bob"},{"id":3,"name":"Charlie"}]
}

// ExamplePipeJSON demonstrates using the PipeJSON convenience function.
func ExamplePipeJSON() {
	// Define a simple structure
	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	// Create a stream of products
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99},
		{ID: 2, Name: "Mouse", Price: 29.99},
	}
	stream := MemReader(products, nil)

	// Use PipeJSON to write directly to stdout
	bytesWritten, err := PipeJSON(stream, os.Stdout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\nBytes written: %d\n", bytesWritten)

	// Output:
	// [{"id":1,"name":"Laptop","price":999.99},{"id":2,"name":"Mouse","price":29.99}]
	// Bytes written: 79
}
