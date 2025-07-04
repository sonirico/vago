package db

import (
	"fmt"
	"strings"

	"errors"

	"github.com/sonirico/vago/lol"
)

// BulkInsert saves to database the given rows by employing INSERT INTO statements
func BulkInsert(
	ctx Context, logger lol.Logger, rows BulkableRanger, tableName string,
) (Rows, error) {
	if rows.Len() < 1 {
		return nil, errors.New("empty rows")
	}

	var (
		cols     = rows.Get(0).Cols()
		colCount = len(cols)

		valuePlaceholders = make([]string, 0, rows.Len())
		args              = make([]any, 0, rows.Len()*colCount)

		cursor       = 0
		placeholders string
	)

	tpl := createTemplate(colCount)

	rows.Range(func(bulkable Bulkable) {
		placeholders, cursor = interpolateTemplate(tpl, colCount, cursor)
		valuePlaceholders = append(valuePlaceholders, placeholders)
		args = append(args, bulkable.Row()...)
	})

	stmt := fmt.Sprintf(`
		INSERT INTO %s
		(%s)
		VALUES
		%s`,
		tableName,
		strings.Join(cols, ","),
		strings.Join(valuePlaceholders, ","),
	)

	logger.Debugln(stmt, args)

	return ctx.Querier().QueryContext(ctx, stmt, args...)
}
