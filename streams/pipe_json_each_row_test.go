package streams

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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
			stream: NewMemory([]mockJsonEachRowMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedJSON:  `{"id":1,"name":"Alice","email":"alice@example.com"}` + "\n" + `{"id":2,"name":"Bob","email":"bob@example.com"}` + "\n",
			expectedError: nil,
		},
		{
			name:          "empty stream",
			stream:        NewMemory([]mockJsonEachRowMarshaler{}, nil),
			expectedJSON:  ``,
			expectedError: nil,
		},
		{
			name: "stream with single record",
			stream: NewMemory([]mockJsonEachRowMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
			}, nil),
			expectedJSON:  `{"id":1,"name":"Alice","email":"alice@example.com"}` + "\n",
			expectedError: nil,
		},
		{
			name:          "stream with error",
			stream:        NewMemory([]mockJsonEachRowMarshaler{}, errors.New("stream error")),
			expectedJSON:  ``,
			expectedError: errors.New("stream error"),
		},
		{
			name:          "stream with error io.EOF",
			stream:        NewMemory([]mockJsonEachRowMarshaler{}, io.EOF),
			expectedJSON:  ``,
			expectedError: nil,
		},
		{
			name: "marshaler with error",
			stream: NewMemory([]mockJsonEachRowMarshaler{
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
	stream := NewMemory([]mockJsonEachRowMarshaler{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}, nil)

	_, err := PipeJSONEachRow(stream, io.Discard)
	assert.NoError(t, err)
}
