package streams

// MemoryWriteStream is a WriteStream implementation that collects items in memory
type MemoryWriteStream[T any] struct {
	items []T
	err   error
}

// Write adds an item to the stream and returns 1 (one item written)
func (w *MemoryWriteStream[T]) Write(item T) (int64, error) {
	if w.err != nil {
		return 0, w.err
	}
	w.items = append(w.items, item)
	return 1, nil
}

// Flush is a no-op for memory streams since items are immediately available
func (w *MemoryWriteStream[T]) Flush() error {
	return w.err
}

// Err returns the current error state of the stream
func (w *MemoryWriteStream[T]) Err() error {
	return w.err
}

// Close finalizes the stream. For memory streams, this is a no-op
func (w *MemoryWriteStream[T]) Close() error {
	return w.err
}

// Items returns the collected items
func (w *MemoryWriteStream[T]) Items() []T {
	return w.items
}

// SetError sets an error state for the stream
func (w *MemoryWriteStream[T]) SetError(err error) {
	w.err = err
}

// MemWriter creates a new memory-based WriteStream
func MemWriter[T any]() *MemoryWriteStream[T] {
	return &MemoryWriteStream[T]{
		items: make([]T, 0),
	}
}
