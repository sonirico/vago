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

// NewReaderStream creates a new ReadStream that reads from an io.Reader
func NewReaderStream(reader io.Reader) ReadStream[[]byte] {
	return &ReaderStream{
		original: reader,
		reader:   bufio.NewReader(reader),
	}
}

// LineReaderStream reads lines as strings instead of bytes
type LineReaderStream struct {
	original io.Reader
	reader   *bufio.Reader
	current  string
	err      error
	closed   bool
}

// Next reads the next line from the reader
func (r *LineReaderStream) Next() bool {
	if r.closed || r.err != nil {
		return false
	}

	line, err := r.reader.ReadString('\n')
	if err != nil {
		r.err = err
		if err == io.EOF && len(line) > 0 {
			// Handle last line without newline
			r.current = line
			return true
		}
		return false
	}

	// Remove trailing newline
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
		// Also remove \r if present (Windows line endings)
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
	}

	r.current = line
	return true
}

// Data returns the current line
func (r *LineReaderStream) Data() string {
	return r.current
}

// Err returns the current error state
func (r *LineReaderStream) Err() error {
	return r.err
}

// Close closes the reader if it implements io.Closer
func (r *LineReaderStream) Close() error {
	if r.closed {
		return nil
	}

	r.closed = true

	if closer, ok := r.original.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// NewLineReaderStream creates a new ReadStream that reads lines as strings from an io.Reader
func NewLineReaderStream(reader io.Reader) ReadStream[string] {
	return &LineReaderStream{
		original: reader,
		reader:   bufio.NewReader(reader),
	}
}
