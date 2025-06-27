package streams

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
