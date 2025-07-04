package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/sonirico/vago/zero"
)

// NullJSON represents a map[K]V that may be null.
// NullJSON implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullJSON[K comparable, V any] struct {
	JSON  map[K]V
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullJSON[K, V]) Scan(value any) error {
	if value == nil {
		n.JSON, n.Valid = nil, false

		return nil
	}

	n.Valid = true

	switch val := value.(type) {
	case []byte:
		return json.Unmarshal(val, &n.JSON)
	case string:
		return json.Unmarshal([]byte(val), &n.JSON)
	default:
		return fmt.Errorf("cannot convert %v of type %T to NullJSON", value, value)
	}
}

// Value implements the driver Valuer interface.
func (n NullJSON[K, V]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return json.Marshal(n.JSON)
}

// NullJSONArray represents a []T that may be null.
// NullJSONArray implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullJSONArray[T any] struct {
	Array []T
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullJSONArray[T]) Scan(value any) error {
	if value == nil {
		n.Array, n.Valid = nil, false

		return nil
	}

	n.Valid = true

	switch val := value.(type) {
	case []byte:
		return json.Unmarshal(val, &n.Array)
	case string:
		return json.Unmarshal(zero.S2B(val), &n.Array)
	}

	return fmt.Errorf("cannot convert %v of type %T to NullJSONArray", value, value)
}

// Value implements the driver Valuer interface.
func (n NullJSONArray[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return json.Marshal(n.Array)
}
