package db

import (
	"fmt"
	"strings"
)

type (
	Bulkable interface {
		PK() []string
		UniqueKeys() []string
		IncludePKOnUpsert() bool
		Cols() []string
		Row() []any
	}

	BulkableRanger interface {
		Get(int) Bulkable
		Len() int
		Range(func(bulkable Bulkable))
	}

	BulkRanger[T Bulkable] []T
)

func (d BulkRanger[T]) Len() int {
	return len(d)
}

func (d BulkRanger[T]) Get(i int) Bulkable {
	return d[i]
}

func (d BulkRanger[T]) Range(fn func(item Bulkable)) {
	for _, m := range d {
		fn(m)
	}
}

func createTemplate(count int) string {
	arr := make([]string, count)
	for i := 0; i < count; i++ {
		arr[i] = "$%d"
	}
	return fmt.Sprintf("(%s)", strings.Join(arr, ","))
}

func interpolateTemplate(tpl string, colCount, cursor int) (string, int) {
	limit := cursor + colCount
	arr := make([]any, colCount)
	i := 0

	for cursor < limit {
		cursor++
		arr[i] = cursor
		i++
	}

	return fmt.Sprintf(tpl, arr...), cursor
}

func createInterpolatedTemplate(cursor int, cols [][2]string) (int, string) {
	count := len(cols)
	arr := make([]string, count)
	for i := 0; i < count; i++ {
		arr[i] = fmt.Sprintf("$%d::%s", cursor+i, cols[i][0])
	}
	return cursor + count, fmt.Sprintf("(%s)", strings.Join(arr, ", "))
}
