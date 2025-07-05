package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testScanner struct {
	called bool
}

func (s *testScanner) Scan(dest ...any) error {
	s.called = true
	return nil
}
func (s *testScanner) Err() error { return nil }

func TestRowScanner(t *testing.T) {
	s := &testScanner{}
	var scanner RowScanner = s
	assert.NoError(t, scanner.Scan())
	assert.True(t, s.called)
	assert.NoError(t, scanner.Err())
}
