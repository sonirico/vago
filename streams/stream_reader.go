package streams

import (
	"bufio"
	"io"
)

// ReaderStream is a ReadStream implementation that reads from an io.Reader
// This is useful for reading streams of bytes from files, network connections, etc.
type ReaderStream struct {
	original io.Reader
	reader   *bufio.Reader
	current  []byte
	err      error
	closed   bool
}

// Next reads the next chunk of data from the reader
func (r *ReaderStream) Next() bool {
	if r.closed || r.err != nil {
		return false
	}

	// Read a line or chunk of data
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		r.err = err
		if err == io.EOF && len(line) > 0 {
			// Handle last line without newline
			r.current = line
			return true
		}
		return false
	}

	r.current = line
	return true
}

// Data returns the current chunk of data
func (r *ReaderStream) Data() []byte {
	return r.current
}

// Err returns the current error state
func (r *ReaderStream) Err() error {
	return r.err
}

// Close closes the reader if it implements io.Closer
func (r *ReaderStream) Close() error {
	if r.closed {
		return nil
	}

	r.closed = true

	if closer, ok := r.original.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Reader creates a new ReadStream that reads from an io.Reader
func Reader(reader io.Reader) ReadStream[[]byte] {
	return &ReaderStream{
		original: reader,
		reader:   bufio.NewReader(reader),
	}
}
