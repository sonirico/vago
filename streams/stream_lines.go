package streams

import (
	"bufio"
	"io"
)

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
		if err == io.EOF {
			if len(line) > 0 {
				// Handle last line without newline
				r.current = line
				// Remove trailing newline if present
				if len(line) > 0 && line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
					// Also remove \r if present (Windows line endings)
					if len(line) > 0 && line[len(line)-1] == '\r' {
						line = line[:len(line)-1]
					}
				}
				r.current = line
				// Mark EOF as processed but don't set it as error yet
				r.err = err
				return true
			}
			// EOF with no data, normal end
			return false
		}
		// Other errors are real errors
		r.err = err
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
	// Don't report EOF as an error - it's the normal end of stream
	if r.err == io.EOF {
		return nil
	}
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

// Lines creates a new ReadStream that reads lines as strings from an io.Reader
func Lines(reader io.Reader) ReadStream[string] {
	return &LineReaderStream{
		original: reader,
		reader:   bufio.NewReader(reader),
	}
}
