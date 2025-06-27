package streams

import (
	"encoding/csv"
	"fmt"
	"io"

	"errors"
)

const (
	CSVSeparatorComma    rune   = ','
	CSVSeparatorCommaStr string = ","
	CSVSeparatorTab      rune   = '\t'
	CSVSeparatorTabStr   string = "\t"
)

// csvMarshaler defines an interface for types that can marshal their data into CSV format.
type csvMarshaler interface {
	MarshalCSV() ([]string, []string, error) // Returns header, record, and error
}

// TransformCSV transforms a stream of csvMarshaler into a CSV-formatted stream.
type TransformCSV[T csvMarshaler] struct {
	stream    ReadStream[T]
	separator rune
}

// WriteTo writes the transformed stream into an io.Writer in CSV format.
// It writes the header only once and streams the records as they are marshaled.
// Returned written result means the actual number of rows written, and not the
// total bytes. Inner csv writer does not provide that information and it's
// inefficient to do it here.
func (p *TransformCSV[T]) WriteTo(w io.Writer) (written int64, err error) {
	var (
		headerWritten bool
	)

	writer := csv.NewWriter(w)
	writer.Comma = p.separator
	defer writer.Flush()

	for p.stream.Next() {
		if err = p.stream.Err(); err != nil {
			if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("stream err: %w", err)
			}
			return
		}

		var (
			header, record []string
			data           = p.stream.Data()
		)
		header, record, err = data.MarshalCSV()
		if err != nil {
			return written, fmt.Errorf("csv marshaling error: %w", err)
		}

		// Write header only once
		if !headerWritten {
			headerWritten = true
			if err = writer.Write(header); err != nil {
				return written, fmt.Errorf("csv header write error: %w", err)
			}
			written++
		}

		// Write record
		if err = writer.Write(record); err != nil {
			return written, fmt.Errorf("csv record write error: %w", err)
		}
		written++
	}

	// Check for any remaining errors from the stream
	if err = p.stream.Err(); err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		} else {
			err = fmt.Errorf("final stream error: %w", err)
		}
	}

	return
}

// CSVTransform creates a new PipeCSVTransform for a given stream.
func CSVTransform[T csvMarshaler](stream ReadStream[T], separator rune) Transform[T] {
	return &TransformCSV[T]{
		stream:    stream,
		separator: separator,
	}
}
