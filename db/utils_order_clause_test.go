package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		order  OrderBy
		asc    bool
		desc   bool
		sign   string
		str    string
		clause string
	}{
		{"asc", OrderASC, true, false, ">", "ASC", "foo ASC"},
		{"desc", OrderDESC, false, true, "<", "DESC", "foo DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.asc, tt.order.IsAsc())
			assert.Equal(t, tt.desc, tt.order.IsDesc())
			assert.Equal(t, tt.sign, tt.order.Sign())
			assert.Equal(t, tt.str, tt.order.String())
			assert.Equal(t, " ORDER BY "+tt.clause, tt.order.FullClause("foo"))
			assert.Equal(t, tt.clause, tt.order.Clause("foo"))
		})
	}
}

// ExampleOrderBy_usage demonstrates how to use OrderBy types to generate SQL ORDER BY clauses.
func ExampleOrderBy_usage() {
	fmt.Println(OrderASC.FullClause("foo"))
	fmt.Println(OrderDESC.FullClause("foo"))
	// Output:
	//  ORDER BY foo ASC
	//  ORDER BY foo DESC
}
