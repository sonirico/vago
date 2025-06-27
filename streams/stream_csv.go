package streams

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
)

type csvUnmarshaler interface {
	UnmarshalCSV([]string) error
}

type CSVStream[T any] struct {
	reader    io.ReadCloser
	buf       *bufio.Scanner
	sep       string
	err       error
	curr      T
	parseFunc func([]string) (T, error)
}

// CSV creates a new CSVStream that reads from a given reader or file path.
func CSV[T any](opts ...CSVOpt) (*CSVStream[T], error) {
	optsDef := &csvOpts{
		flag:   os.O_RDONLY,
		perm:   0644,
		sep:    CSVSeparatorCommaStr,
		reader: os.Stdin,
		path:   "",
	}

	for _, opt := range opts {
		opt.apply(optsDef)
	}

	if optsDef.path == "" {
		file, err := os.OpenFile(optsDef.path, optsDef.flag, optsDef.perm)
		if err != nil {
			return nil, err
		}
		return newStreamCSV[T](file, optsDef.sep), nil
	}

	return newStreamCSV[T](optsDef.reader, optsDef.sep), nil
}

func newStreamCSV[T any](r io.ReadCloser, sep string) *CSVStream[T] {
	var zero T
	stream := &CSVStream[T]{
		sep:    sep,
		reader: r,
		buf:    bufio.NewScanner(r),
		curr:   zero,
	}

	// Determine the parsing strategy based on type T
	var _ []string = nil
	var nilSlice T

	// Check if T is []string
	if _, ok := any(nilSlice).([]string); ok {
		stream.parseFunc = func(data []string) (T, error) {
			return any(data).(T), nil
		}
	} else {
		// Check if T implements csvUnmarshaler
		var value T
		if _, ok := any(&value).(csvUnmarshaler); ok {
			stream.parseFunc = func(data []string) (T, error) {
				var value T
				if valuePtr, ok := any(&value).(csvUnmarshaler); ok {
					if err := valuePtr.UnmarshalCSV(data); err != nil {
						return zero, err
					}
					return value, nil
				}
				return zero, fmt.Errorf("type does not implement csvUnmarshaler")
			}
		} else {
			// If we get here, T is neither []string nor implements csvUnmarshaler
			stream.parseFunc = func(data []string) (T, error) {
				return zero, fmt.Errorf("type must be []string or implement csvUnmarshaler")
			}
		}
	}

	return stream
}

func (s *CSVStream[T]) Next() bool {
	if s.buf.Scan() {
		data := strings.Split(s.buf.Text(), s.sep)
		value, err := s.parseFunc(data)
		if err != nil {
			s.err = err
			return false
		}
		s.curr = value
		return true
	}
	return false
}

func (s *CSVStream[T]) Data() T {
	return s.curr
}

func (s *CSVStream[T]) Err() error {
	return s.err
}

func (s *CSVStream[T]) Close() error {
	return s.reader.Close()
}

func (s *CSVStream[T]) Iter() iter.Seq[T] {
	return Iter(s)
}

func (s *CSVStream[T]) Iter2() iter.Seq2[T, error] {
	return Iter2(s)
}
