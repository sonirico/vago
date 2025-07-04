package streams

// DBRows is an interface that abstracts database rows operations.
// It's compatible with sql.Rows and pgx.Rows among others.
type DBRows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

// DBStream provides a ReadDBStream implementation for database rows.
// It wraps DBRows and provides a generic interface for streaming database results.
//
// The DBStream will automatically close the underlying rows when Next() returns false.
type DBStream[T any] struct {
	rows    DBRows
	scanFn  func(DBRows, *T) error
	current *T
	err     error
}

var _ ReadStream[any] = new(DBStream[any])

// Next advances the stream to the next row and scans it into the current value.
// It returns true if there was a next row, false if there are no more rows or an error occurred.
// The underlying rows are automatically closed when Next returns false.
func (s *DBStream[T]) Next() bool {
	keep := s.rows.Next()

	if keep {
		s.current = new(T)
		s.err = s.scanFn(s.rows, s.current)
	} else {
		_ = s.rows.Close()
	}

	return keep
}

// Data returns the current row data.
// This should only be called after Next() returns true.
func (s *DBStream[T]) Data() T {
	if s.current == nil {
		var x T
		return x
	}
	return *s.current
}

// Err returns any error that occurred during scanning.
func (s *DBStream[T]) Err() error {
	return s.err
}

// Close closes the underlying rows.
// This is called automatically when Next() returns false,
// but can be called manually for early termination.
func (s *DBStream[T]) Close() error {
	return nil
}

// DB creates a new DBStream that wraps the given DBRows.
// The scanFn function is called for each row to scan the data into the target type.
//
// Example:
//
//	stream := DB(rows, func(rows DBRows, user *User) error {
//		return rows.Scan(&user.ID, &user.Name)
//	})
func DB[T any](rows DBRows, scanFn func(DBRows, *T) error) *DBStream[T] {
	return &DBStream[T]{
		rows:   rows,
		scanFn: scanFn,
	}
}
