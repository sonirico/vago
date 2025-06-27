package streams

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryWriteStream(t *testing.T) {
	stream := MemWriter[string]()

	// Test writing items
	items := []string{"hello", "world", "test"}
	var totalBytes int64

	for _, item := range items {
		n, err := stream.Write(item)
		require.NoError(t, err, "Should write item %s without error", item)
		assert.Positive(t, n, "Should write positive bytes for %s", item)
		totalBytes += n
	}

	// Test retrieving items
	result := stream.Items()
	require.Len(t, result, len(items), "Should have correct number of items")

	for i, expected := range items {
		assert.Equal(t, expected, result[i], "Item %d should match", i)
	}

	// Test Flush and Close
	assert.NoError(t, stream.Flush(), "Flush should not return error")
	assert.NoError(t, stream.Close(), "Close should not return error")

	// Test error state
	assert.NoError(t, stream.Err(), "Err should be nil")
}

func TestMemoryWriteStreamWithError(t *testing.T) {
	stream := MemWriter[int]()

	// Write some items first
	_, err := stream.Write(1)
	require.NoError(t, err, "Should write first item without error")

	// Set an error
	testErr := errors.New("test error")
	stream.SetError(testErr)

	// Try to write after error - should return the error
	_, err = stream.Write(2)
	assert.Error(t, err, "Should return error when writing after SetError")

	// Check error state
	assert.Error(t, stream.Err(), "Should have error state set")
}

func TestWriterStream(t *testing.T) {
	var buf bytes.Buffer
	stream := Writer(&buf)

	// Test writing data
	testData := [][]byte{
		[]byte("hello "),
		[]byte("world"),
		[]byte("!"),
	}

	var totalBytes int64
	for _, data := range testData {
		n, err := stream.Write(data)
		require.NoError(t, err, "Should write data without error")
		assert.Equal(t, int64(len(data)), n, "Should write correct number of bytes")
		totalBytes += n
	}

	// Check the buffer content
	expected := "hello world!"
	assert.Equal(t, expected, buf.String(), "Buffer should contain expected content")

	// Test Flush and Close
	assert.NoError(t, stream.Flush(), "Flush should not return error for bytes.Buffer")
	assert.NoError(t, stream.Close(), "Close should not return error for bytes.Buffer")

	// Test error state
	assert.NoError(t, stream.Err(), "Err should be nil")

	// Verify the total bytes written matches buffer size
	expectedBytes := int64(len("hello world!"))
	assert.Equal(t, expectedBytes, totalBytes, "Should write expected total bytes")

	// Verify the actual buffer size matches
	assert.Equal(t, expectedBytes, int64(buf.Len()), "Buffer size should match expected bytes")
}

func TestWriterStreamInterface(t *testing.T) {
	// Test that MemoryWriteStream implements WriteStream interface
	var _ WriteStream[string] = &MemoryWriteStream[string]{}

	// Test that WriterStream implements WriteStream interface for []byte
	var _ WriteStream[[]byte] = &WriterStream{}
}

func TestWriteAll(t *testing.T) {
	stream := MemWriter[int]()
	items := []int{1, 2, 3, 4, 5}

	bytesWritten, err := WriteAll(stream, items)
	require.NoError(t, err, "WriteAll should succeed")
	assert.Positive(t, bytesWritten, "Should write positive bytes")

	result := stream.Items()
	require.Len(t, result, len(items), "Should have all items")

	for i, expected := range items {
		assert.Equal(t, expected, result[i], "Item %d should match", i)
	}
}

func TestPipeStream(t *testing.T) {
	// Create source data
	sourceItems := []string{"hello", "world", "test"}
	src := MemReader(sourceItems, nil)

	// Create destination stream
	dst := MemWriter[string]()

	// Copy from source to destination
	bytesWritten, err := Pipe(src, dst)
	require.NoError(t, err, "Pipe should succeed")
	assert.Positive(t, bytesWritten, "Should write positive bytes")

	// Verify the copy
	result := dst.Items()
	require.Len(t, result, len(sourceItems), "Should have all source items")

	for i, expected := range sourceItems {
		assert.Equal(t, expected, result[i], "Item %d should match", i)
	}
}
