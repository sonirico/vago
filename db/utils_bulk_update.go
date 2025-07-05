package db

import (
	"fmt"
	"strings"

	"errors"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/slices"
)

type (
	BulkUpdatable interface {
		PK() [2]string
		BulkUpdateCols() [][2]string
		BulkUpdateValues() []any
	}

	bulkUpdate interface {
		Len() int
		Get(i int) BulkUpdatable
		Range(func(x BulkUpdatable))
	}
)

// TODO: Make BulkUpdateReturning

func BulkUpdate(
	ctx Context, logger lol.Logger, rows bulkUpdate, tableName string,
) (Result, error) {
	if rows.Len() < 1 {
		return nil, errors.New("empty rows")
	}

	var (
		row  = rows.Get(0)
		cols = row.BulkUpdateCols()

		valuePlaceholders = make([]string, 0, rows.Len()*(len(cols)))
		args              = make([]any, 0, rows.Len()*len(cols))
	)

	pk := row.PK()
	pkType, pkName := pk[0], pk[1]

	var (
		setStmt = strings.Join(slices.Map(cols, func(t [2]string) string {
			colType, colName := t[0], t[1]
			return fmt.Sprintf("%s = bulk_update_tmp.%s::%s", colName, colName, colType)
		}), ",")

		whereStmt = fmt.Sprintf(
			"%s.%s::%s = bulk_update_tmp.%s::%s",
			tableName,
			pkName,
			pkType,
			pkName,
			pkType,
		)

		tupleDef = fmt.Sprintf(
			"%s, %s",
			pkName,
			strings.Join(slices.Map(cols, func(colDef [2]string) string {
				return colDef[1]
			}), ", "),
		)

		cursor = 1
		tpl    string
	)

	colsWithPk := append([][2]string{pk}, cols...)

	rows.Range(func(item BulkUpdatable) {
		cursor, tpl = createInterpolatedTemplate(cursor, colsWithPk) // len(cols) + pk
		valuePlaceholders = append(valuePlaceholders, tpl)
		args = append(args, item.BulkUpdateValues()...)
	})

	valuesStmt := strings.Join(valuePlaceholders, ", ")

	stmt := fmt.Sprintf(`
		UPDATE %s
		SET %s
		FROM (VALUES %s) as bulk_update_tmp(%s)
		WHERE %s
		`,
		tableName,
		setStmt,
		valuesStmt,
		tupleDef,
		whereStmt,
	)

	logger.Debugln(stmt, args)

	return ctx.Querier().ExecContext(ctx, stmt, args...)
}
