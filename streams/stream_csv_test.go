package streams

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type csvTest struct {
	Nombre   string
	Apellido string
}

func (t *csvTest) UnmarshalCSV(data []string) error {
	if len(data) < 2 {
		return fmt.Errorf("want at least 2 cols, have %d", len(data))
	}
	t.Nombre = data[0]
	t.Apellido = data[1]
	return nil
}

func TestNewCSVStream(t *testing.T) {
	buf := io.NopCloser(strings.NewReader(strings.TrimSpace(`
nombre,apellido
nombre2,apellido2,
nombre3,apellido3, , 


`)))

	s := NewStreamCSV[csvTest](buf, ",")

	expected := []csvTest{
		{"nombre", "apellido"},
		{"nombre2", "apellido2"},
		{"nombre3", "apellido3"},
	}
	i := 0

	for s.Next() {
		if err := s.Err(); err != nil {
			t.Fatalf("unexpected err %v", err)
		}

		actual := s.Data()

		assert.Equal(t, actual, expected[i])
		i++
	}
}

func TestCSVStreamClose(t *testing.T) {
	buf := io.NopCloser(strings.NewReader("nombre,apellido"))
	s := NewStreamCSV[csvTest](buf, ",")

	err := s.Close()
	assert.NoError(t, err)
}

func TestCSVStreamWithStringSlice(t *testing.T) {
	buf := io.NopCloser(strings.NewReader(strings.TrimSpace(`
field1,field2,field3
value1,value2,value3
a,b,c
`)))

	s := NewStreamCSV[[]string](buf, ",")

	expected := [][]string{
		{"field1", "field2", "field3"},
		{"value1", "value2", "value3"},
		{"a", "b", "c"},
	}
	i := 0

	for s.Next() {
		if err := s.Err(); err != nil {
			t.Fatalf("unexpected err %v", err)
		}

		actual := s.Data()
		assert.Equal(t, expected[i], actual)
		i++
	}
}
