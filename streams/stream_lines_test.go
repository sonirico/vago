package streams

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: nil, // Empty slice can be nil
		},
		{
			name:     "single line",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "single line with newline",
			input:    "hello\n",
			expected: []string{"hello"},
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "multiple lines with trailing newline",
			input:    "line1\nline2\nline3\n",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "windows line endings",
			input:    "line1\r\nline2\r\n",
			expected: []string{"line1", "line2"},
		},
		{
			name:     "empty lines",
			input:    "line1\n\nline3\n",
			expected: []string{"line1", "", "line3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(test.input)
			stream := Lines(reader)

			var result []string
			for stream.Next() {
				result = append(result, stream.Data())
			}

			assert.NoError(t, stream.Err(), "unexpected error")
			assert.Equal(t, test.expected, result, "unexpected result")
		})
	}
}

func TestLinesClose(t *testing.T) {
	reader := strings.NewReader("line1\nline2\n")
	stream := Lines(reader)

	// Read first line
	assert.True(t, stream.Next(), "expected to read first line")

	// Close the stream
	assert.NoError(t, stream.Close(), "unexpected error closing stream")

	// Try to read next line (should fail)
	assert.False(t, stream.Next(), "expected stream to be closed")
}

// ExampleLines demonstrates reading lines from a string.
func ExampleLines() {
	// Create a reader from a multiline string
	text := "line1\nline2\nline3\n"
	reader := strings.NewReader(text)

	// Create a lines stream
	lineStream := Lines(reader)

	// Collect the results
	result, _ := Consume(lineStream)
	fmt.Println(result)
	// Output: [line1 line2 line3]
}
