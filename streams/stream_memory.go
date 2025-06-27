package streams

type MemoryStream[T any] struct {
	items  []T
	cursor int
	error  error
}

func (s *MemoryStream[T]) Next() bool {
	s.cursor++

	return s.cursor < len(s.items)
}

func (s *MemoryStream[T]) Data() T {
	return s.items[s.cursor]
}

func (s *MemoryStream[T]) Err() error {
	return s.error
}

func (s *MemoryStream[T]) Close() error {
	return nil
}

// Mem creates a new ReadStream that reads from a slice in memory.
// This is useful for testing, converting slices to streams, or creating
// simple data sources for streaming pipelines.
//
// The error parameter allows you to simulate error conditions during streaming.
// If err is not nil, the stream will return this error when Err() is called.
func MemReader[T any](items []T, err error) *MemoryStream[T] {
	return &MemoryStream[T]{
		items:  items,
		cursor: -1,
		error:  err,
	}
}

var _ ReadStream[any] = new(MapperStreamErr[any, any])
