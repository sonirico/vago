package db

import (
	"fmt"
	"strings"

	"errors"

	"github.com/sonirico/vago/lol"
	"github.com/sonirico/vago/slices"
)

func buildSqlStmt(
	rows BulkableRanger,
	tableName string,
	updateOnConflict bool,
) (string, []any, error) {
	var (
		row             = rows.Get(0)
		pk              = row.PK()
		uniqueKeys      = row.UniqueKeys()
		includePK       = row.IncludePKOnUpsert()
		cols            = row.Cols()
		upsertColsCount = len(cols)
		insertCols      = make([]string, 0, upsertColsCount)

		valuePlaceholders = make([]string, 0, rows.Len())
		args              = make([]any, 0, rows.Len()*upsertColsCount)
		onUpdateCols      = make([]string, 0, rows.Len())
		onConflictStmt    = "DO NOTHING"

		cursor       = 0
		placeholders string
	)

	if len(uniqueKeys) == 0 && !includePK {
		// Return error wraping NoUniqueConstraintError
		return "", nil, fmt.Errorf("failed to create bulk upser for %s %w",
			tableName, NoUniqueConstraintError)
	}

	if len(uniqueKeys) == 0 {
		uniqueKeys = append(uniqueKeys, pk...)
	} else if !includePK {
		upsertColsCount -= len(pk)
	}

	tpl := createTemplate(upsertColsCount)

	rows.Range(func(bulkable Bulkable) {
		placeholders, cursor = interpolateTemplate(tpl, upsertColsCount, cursor)
		valuePlaceholders = append(valuePlaceholders, placeholders)
		args = append(args, bulkable.Row()...)
	})

	slices.Range(cols, func(col string, _ int) bool {
		if slices.Includes(pk, col) && !includePK {
			return false
		}
		if updateOnConflict {
			onUpdateCols = append(onUpdateCols, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
		}

		insertCols = append(insertCols, col)

		return true
	})

	if updateOnConflict {
		onConflictStmt = fmt.Sprintf(`
			DO UPDATE SET
			%s`,
			strings.Join(onUpdateCols, ","),
		)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO %s
		(%s)
		VALUES
		%s
		ON CONFLICT (%s) %s
		RETURNING %s
		`,
		tableName,
		strings.Join(insertCols, ","),
		strings.Join(valuePlaceholders, ","),
		strings.Join(uniqueKeys, ","),
		onConflictStmt,
		strings.Join(cols, ","),
	)
	return stmt, args, nil
}

func BulkUpsert(
	ctx Context, logger lol.Logger, rows BulkableRanger, tableName string,
	updateOnConflict bool,
) (Rows, error) {
	var (
		query string
		args  []any
		err   error
	)

	if rows.Len() < 1 {
		return nil, errors.New("empty rows")
	}

	if query, args, err = buildSqlStmt(rows, tableName, updateOnConflict); err != nil {
		return nil, err
	}
	logger.Debugln(query, args)

	return ctx.Querier().QueryContext(ctx, query, args...)
}
