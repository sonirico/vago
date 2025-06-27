package streams

import (
	"bytes"
	"encoding/json"
	"errors"
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
			stream: NewMemory([]mockJsonMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedJSON:  `[{"id":1,"name":"Alice","email":"alice@example.com"},{"id":2,"name":"Bob","email":"bob@example.com"}]`,
			expectedError: nil,
		},
		{
			name:          "empty stream",
			stream:        NewMemory([]mockJsonMarshaler{}, nil),
			expectedJSON:  `[]`,
			expectedError: nil,
		},
		{
			name: "stream with single record",
			stream: NewMemory([]mockJsonMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
			}, nil),
			expectedJSON:  `[{"id":1,"name":"Alice","email":"alice@example.com"}]`,
			expectedError: nil,
		},
		{
			name:          "stream with error",
			stream:        NewMemory([]mockJsonMarshaler{}, errors.New("stream error")),
			expectedJSON:  ``,
			expectedError: errors.New("stream error"),
		},
		{
			name:          "stream with error io.EOF",
			stream:        NewMemory([]mockJsonMarshaler{}, io.EOF),
			expectedJSON:  `[]`,
			expectedError: nil,
		},
		{
			name: "marshaler with error",
			stream: NewMemory([]mockJsonMarshaler{
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
	stream := NewMemory([]mockJsonMarshaler{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}, nil)

	_, err := JSONTransform(stream).WriteTo(os.Stdout)
	assert.NoError(t, err)
}
