package db

import (
	"fmt"
)

func In[T any](args []any, inArgs []T) (string, []any) {
	start := len(args) + 1
	end := start + len(inArgs)
	j := 0
	inStmt := ""

	for i := start; i < end; i++ {
		item := fmt.Sprintf("$%d", i)
		step := ","
		if j >= len(inArgs)-1 {
			step = ""
		}
		args = append(args, inArgs[j])
		inStmt = inStmt + item + step
		j++
	}

	inStmt = fmt.Sprintf("(%s)", inStmt)
	return inStmt, args
}
