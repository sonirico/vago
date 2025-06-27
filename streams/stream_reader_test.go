package streams

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderStream(t *testing.T) {
	// Test reading bytes from a string
	testData := "hello\nworld\ntest\n"
	reader := strings.NewReader(testData)
	stream := Reader(reader)

	var result [][]byte
	for stream.Next() {
		result = append(result, stream.Data())
	}

	// Check for errors
	err := stream.Err()
	if err != nil {
		assert.Equal(t, io.EOF, err, "Should only have EOF error")
	}

	// Should have 3 lines
	expected := []string{"hello\n", "world\n", "test\n"}
	require.Len(t, result, len(expected), "Should have correct number of lines")

	for i, expectedLine := range expected {
		assert.Equal(t, expectedLine, string(result[i]), "Line %d should match", i)
	}
}

func TestReaderStreamNoTrailingNewline(t *testing.T) {
	// Test reading data without trailing newline
	testData := "hello\nworld\ntest"
	reader := strings.NewReader(testData)
	stream := Reader(reader)

	result, err := Consume(stream)
	require.NoError(t, err, "Should consume stream without error")

	// Should have 3 chunks
	expected := []string{"hello\n", "world\n", "test"}
	require.Len(t, result, len(expected), "Should have correct number of chunks")

	for i, expectedChunk := range expected {
		assert.Equal(t, expectedChunk, string(result[i]), "Chunk %d should match", i)
	}
}

func TestReaderStreamEmpty(t *testing.T) {
	// Test reading from empty reader
	reader := strings.NewReader("")
	stream := Reader(reader)

	result, err := Consume(stream)
	require.NoError(t, err, "Should consume empty stream without error")
	assert.Empty(t, result, "Should have empty result")
}

func TestReaderStreamClose(t *testing.T) {
	// Test closing the stream
	testData := "hello\nworld\n"
	reader := strings.NewReader(testData)
	stream := Reader(reader)

	// Read first line
	assert.True(t, stream.Next(), "Should read first line")

	// Close the stream
	require.NoError(t, stream.Close(), "Should close stream without error")

	// Try to read more - should not work
	assert.False(t, stream.Next(), "Should not read after closing")
}

func TestLineReaderStream(t *testing.T) {
	// Test reading lines as strings
	testData := "hello\nworld\ntest line\n"
	reader := strings.NewReader(testData)
	stream := Lines(reader)

	var result []string
	for stream.Next() {
		result = append(result, stream.Data())
	}

	// Check for errors
	err := stream.Err()
	if err != nil {
		assert.Equal(t, io.EOF, err, "Should only have EOF error")
	}

	// Should have 3 lines without newlines
	expected := []string{"hello", "world", "test line"}
	require.Len(t, result, len(expected), "Should have correct number of lines")

	for i, expectedLine := range expected {
		assert.Equal(t, expectedLine, result[i], "Line %d should match", i)
	}
}

func TestLineReaderStreamWindowsLineEndings(t *testing.T) {
	// Test reading Windows line endings (\r\n)
	testData := "hello\r\nworld\r\ntest\r\n"
	reader := strings.NewReader(testData)
	stream := Lines(reader)

	result, err := Consume(stream)
	require.NoError(t, err, "Should consume stream without error")

	expected := []string{"hello", "world", "test"}
	require.Len(t, result, len(expected), "Should have correct number of lines")

	for i, expectedLine := range expected {
		assert.Equal(t, expectedLine, result[i], "Line %d should match", i)
	}
}

func TestLineReaderStreamNoTrailingNewline(t *testing.T) {
	// Test reading lines without trailing newline
	testData := "hello\nworld\ntest"
	reader := strings.NewReader(testData)
	stream := Lines(reader)

	result, err := Consume(stream)
	require.NoError(t, err, "Should consume stream without error")

	expected := []string{"hello", "world", "test"}
	require.Len(t, result, len(expected), "Should have correct number of lines")

	for i, expectedLine := range expected {
		assert.Equal(t, expectedLine, result[i], "Line %d should match", i)
	}
}

func TestLineReaderStreamEmpty(t *testing.T) {
	// Test reading from empty reader
	reader := strings.NewReader("")
	stream := Lines(reader)

	result, err := Consume(stream)
	require.NoError(t, err, "Should consume empty stream without error")
	assert.Empty(t, result, "Should have empty result")
}

func TestReaderStreamInterfaces(t *testing.T) {
	// Test that streams implement ReadStream interface
	var _ ReadStream[[]byte] = &ReaderStream{}
	var _ ReadStream[string] = &LineReaderStream{}
}

func TestReaderStreamWithWriteStream(t *testing.T) {
	// Test integration: read from one stream and write to another
	testData := "line1\nline2\nline3\n"
	reader := strings.NewReader(testData)
	readStream := Reader(reader)

	// Write to memory stream
	writeStream := MemWriter[[]byte]()

	// Connect them
	bytesWritten, err := Pipe(readStream, writeStream)
	require.NoError(t, err, "Should connect streams without error")
	assert.Equal(t, int64(3), bytesWritten, "Should write 3 items")

	// Verify the data
	result := writeStream.Items()
	require.Len(t, result, 3, "Should have 3 items")

	expected := []string{"line1\n", "line2\n", "line3\n"}
	for i, expectedLine := range expected {
		assert.Equal(t, expectedLine, string(result[i]), "Item %d should match", i)
	}
}

func TestLineReaderStreamWithFilter(t *testing.T) {
	// Test integration with FilterStream
	testData := "hello\nworld\ntest\nfilter\nstream\n"
	reader := strings.NewReader(testData)
	readStream := Lines(reader)

	// Filter lines with length > 4
	filtered := Filter(readStream, func(line string) bool {
		return len(line) > 4
	})

	result, err := Consume(filtered)
	require.NoError(t, err, "Should consume filtered stream without error")

	expected := []string{"hello", "world", "filter", "stream"}
	require.Len(t, result, len(expected), "Should have correct number of filtered lines")

	for i, expectedLine := range expected {
		assert.Equal(t, expectedLine, result[i], "Line %d should match", i)
	}
}

func TestReaderStreamLargeData(t *testing.T) {
	// Test with larger data to ensure buffering works
	var buf bytes.Buffer
	for i := 0; i < 1000; i++ {
		buf.WriteString("line ")
		buf.WriteString(string(rune('0' + i%10)))
		buf.WriteString("\n")
	}

	stream := Lines(&buf)
	result, err := Consume(stream)
	require.NoError(t, err, "Should consume large stream without error")
	assert.Len(t, result, 1000, "Should have 1000 lines")

	// Check a few sample lines
	assert.Equal(t, "line 0", result[0], "First line should match")
	assert.Equal(t, "line 9", result[999], "Last line should match")
}

func TestReaderStreamErrorHandling(t *testing.T) {
	// Test error handling during read
	reader := &errorReader{failAfter: 2}
	stream := Reader(reader)

	// Should read some data before failing
	assert.True(t, stream.Next(), "Should read first chunk")
	assert.True(t, stream.Next(), "Should read second chunk")

	// Third read should fail
	assert.False(t, stream.Next(), "Should fail on third read")
	assert.Error(t, stream.Err(), "Should have error after failure")
}

func TestLineReaderStreamErrorHandling(t *testing.T) {
	// Test error handling for line reader
	reader := &errorReader{failAfter: 1}
	stream := Lines(reader)

	// First read should work
	assert.True(t, stream.Next(), "Should read first line")

	// Second read should fail
	assert.False(t, stream.Next(), "Should fail on second read")
	assert.Error(t, stream.Err(), "Should have error after failure")
}

func TestReaderStreamMultipleClose(t *testing.T) {
	// Test that multiple closes don't cause issues
	reader := strings.NewReader("test\n")
	stream := Reader(reader)

	// Close multiple times should not error
	assert.NoError(t, stream.Close(), "First close should succeed")
	assert.NoError(t, stream.Close(), "Second close should succeed")
	assert.NoError(t, stream.Close(), "Third close should succeed")
}

func TestReaderStreamChaining(t *testing.T) {
	// Test complex chaining: Reader -> Filter -> Map -> Write
	testData := "apple\nbanana\ncherry\ndate\nelderberry\nfig\n"
	reader := strings.NewReader(testData)
	readStream := Lines(reader)

	// Filter fruits with length > 4
	filtered := Filter(readStream, func(fruit string) bool {
		return len(fruit) > 4
	})

	// Map to uppercase
	mapped := Map(filtered, func(fruit string) string {
		return strings.ToUpper(fruit)
	})

	// Collect results
	result, err := Consume(mapped)
	require.NoError(t, err, "Should process chain without error")

	expected := []string{"APPLE", "BANANA", "CHERRY", "ELDERBERRY"}
	assert.Equal(t, expected, result, "Should have correct filtered and mapped results")
}

func TestReaderStreamWithMulticast(t *testing.T) {
	// Test Multicasting reader stream to multiple writers
	testData := "line1\nline2\nline3\n"
	reader := strings.NewReader(testData)
	readStream := Reader(reader)

	// Create multiple write streams
	writer1 := MemWriter[[]byte]()
	writer2 := MemWriter[[]byte]()
	writer3 := MemWriter[[]byte]()

	// Multicast to all writers
	bytesWritten, err := Multicast(readStream, writer1, writer2, writer3)
	require.NoError(t, err, "Should Multicast without error")

	// All should have written the same amount
	expectedBytes := []int64{3, 3, 3} // 3 items each
	assert.Equal(t, expectedBytes, bytesWritten, "Should write same amount to each")

	// All should have the same data
	expected := []string{"line1\n", "line2\n", "line3\n"}
	writers := []*MemoryWriteStream[[]byte]{writer1, writer2, writer3}

	for i, writer := range writers {
		items := writer.Items()
		require.Len(t, items, 3, "Writer %d should have 3 items", i)

		for j, expectedLine := range expected {
			assert.Equal(t, expectedLine, string(items[j]), "Writer %d, item %d should match", i, j)
		}
	}
}

// errorReader is a test helper that fails after a certain number of reads
type errorReader struct {
	count     int
	failAfter int
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	r.count++
	if r.count > r.failAfter {
		return 0, assert.AnError
	}

	// Return some test data
	data := "test data\n"
	copy(p, data)
	return len(data), nil
}
