package streams

// Reduce applies a reduction function to a ReadStream, accumulating results in a map.
// The function takes the current map and an item from the stream, returning a new map.
// It returns the final accumulated map or an error if the stream encounters one.
func Reduce[T any, K comparable, V any](
	s ReadStream[T],
	fn func(map[K]V, T) map[K]V,
) (map[K]V, error) {
	res := make(map[K]V)
	for s.Next() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		res = fn(res, s.Data())
	}

	return res, s.Err()
}
