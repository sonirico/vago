package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPgxRow_Err(t *testing.T) {
	t.Parallel()

	// Arrange: pgx.Row has no Err method, so pgxRow.Err has nothing to wrap
	// and always reports nil - documented on the method itself.
	row := &pgxRow{}

	// Act
	err := row.Err()

	// Assert
	assert.NoError(t, err)
}
