package streams

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of csvMarshaler for testing
type mockCsvMarshaler struct {
	ID    int
	Name  string
	Email string
	Err   error
}

func (m mockCsvMarshaler) MarshalCSV() ([]string, []string, error) {
	if m.Err != nil {
		return nil, nil, m.Err
	}
	header := []string{"ID", "Name", "Email"}
	record := []string{strconv.FormatInt(int64(m.ID), 10), m.Name, m.Email}
	return header, record, nil
}

func TestPipeCSVTransform_WriteTo(t *testing.T) {
	tests := []struct {
		name          string
		separator     rune
		stream        ReadStream[mockCsvMarshaler]
		expectedCSV   string
		expectedError error
	}{
		{
			name:      "valid stream with multiple records comma separated",
			separator: CSVSeparatorComma,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedCSV:   "ID,Name,Email\n1,Alice,alice@example.com\n2,Bob,bob@example.com\n",
			expectedError: nil,
		},
		{
			name:      "valid stream with multiple records tab separated",
			separator: CSVSeparatorTab,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil),
			expectedCSV:   "ID\tName\tEmail\n1\tAlice\talice@example.com\n2\tBob\tbob@example.com\n",
			expectedError: nil,
		},
		{
			name:          "empty stream",
			separator:     CSVSeparatorComma,
			stream:        MemReader([]mockCsvMarshaler{}, nil),
			expectedCSV:   "",
			expectedError: nil,
		},
		{
			name:      "stream with single record",
			separator: CSVSeparatorComma,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
			}, nil),
			expectedCSV:   "ID,Name,Email\n1,Alice,alice@example.com\n",
			expectedError: nil,
		},
		{
			name:          "stream with error",
			separator:     CSVSeparatorComma,
			stream:        MemReader([]mockCsvMarshaler{}, errors.New("stream error")),
			expectedCSV:   "",
			expectedError: errors.New("stream error"),
		},
		{
			name:          "stream with error io.EOF",
			separator:     CSVSeparatorComma,
			stream:        MemReader([]mockCsvMarshaler{}, io.EOF),
			expectedCSV:   "",
			expectedError: nil,
		},
		{
			name:      "marshaler with error",
			separator: CSVSeparatorComma,
			stream: MemReader([]mockCsvMarshaler{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{
					ID:    2,
					Name:  "Error",
					Email: "error@example.com",
					Err:   errors.New("marshal error"),
				},
			}, nil),
			expectedCSV:   "",
			expectedError: errors.New("csv marshaling error: marshal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			transform := CSVTransform(tt.stream, tt.separator)

			written, err := transform.WriteTo(buffer)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(len(strings.Split(tt.expectedCSV, "\n")))-1, written)
				assert.Equal(t, tt.expectedCSV, buffer.String())
			}
		})
	}
}

func TestPipeCSVTransform_WriteToFile(t *testing.T) {
	stream := MemReader([]mockCsvMarshaler{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}, nil)

	_, err := CSVTransform(stream, CSVSeparatorTab).WriteTo(os.Stdout)
	assert.NoError(t, err)
}
