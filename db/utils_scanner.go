package db

type RowScanner interface {
	Scan(dest ...any) error
	Err() error
}
