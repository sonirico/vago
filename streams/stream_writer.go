package streams

import (
	"io"
)

// WriterStream is a WriteStream implementation that writes to an io.Writer
// This is useful for writing streams of bytes or strings to files, network connections, etc.
type WriterStream struct {
	writer io.Writer
	err    error
}

// Write writes data to the underlying io.Writer and returns bytes written
func (w *WriterStream) Write(data []byte) (int64, error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err := w.writer.Write(data)
	if err != nil {
		w.err = err
	}
	return int64(n), err
}

// Flush attempts to flush the writer if it implements io.Flusher
func (w *WriterStream) Flush() error {
	if w.err != nil {
		return w.err
	}

	if flusher, ok := w.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			w.err = err
		}
		return err
	}
	return nil
}

// Err returns the current error state
func (w *WriterStream) Err() error {
	return w.err
}

// Close closes the writer if it implements io.Closer
func (w *WriterStream) Close() error {
	if w.err != nil {
		return w.err
	}
	if closer, ok := w.writer.(io.Closer); ok {
		w.err = closer.Close()
	}
	return w.err
}

// Writer creates a new WriteStream that writes to an io.Writer
func Writer(writer io.Writer) *WriterStream {
	return &WriterStream{
		writer: writer,
	}
}
