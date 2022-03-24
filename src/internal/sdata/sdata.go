package sdata

import (
	"context"
	"database/sql"

	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/pachyderm/pachyderm/v2/src/internal/pachsql"
)

// MaterializationResult is returned by MaterializeSQL
type MaterializationResult struct {
	ColumnNames   []string
	ColumnDBTypes []string
	RowCount      uint64
}

// MaterializeSQL reads all the rows from a *sql.Rows, and writes them to tw.
// It flushes tw and returns a MaterializationResult
func MaterializeSQL(tw TupleWriter, rows *sql.Rows) (*MaterializationResult, error) {
	colNames, err := rows.Columns()
	if err != nil {
		return nil, errors.EnsureStack(err)
	}
	cTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, errors.EnsureStack(err)
	}
	var dbTypes []string
	for _, cType := range cTypes {
		dbTypes = append(dbTypes, cType.DatabaseTypeName())
	}
	var count uint64
	row := newTupleFromSQL(cTypes)
	for rows.Next() {
		if err := rows.Scan(row...); err != nil {
			return nil, errors.EnsureStack(err)
		}
		if err := tw.WriteTuple(row); err != nil {
			return nil, errors.EnsureStack(err)
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return nil, errors.EnsureStack(err)
	}
	if err := tw.Flush(); err != nil {
		return nil, errors.EnsureStack(err)
	}
	return &MaterializationResult{
		ColumnNames:   colNames,
		ColumnDBTypes: dbTypes,
		RowCount:      count,
	}, nil
}

func newTupleFromSQL(colTypes []*sql.ColumnType) Tuple {
	row := make([]interface{}, len(colTypes))
	for i, cType := range colTypes {
		var err error
		isNull, nullOk := cType.Nullable()
		row[i], err = makeTupleElement(cType.DatabaseTypeName(), isNull || nullOk)
		if err != nil {
			panic(err)
		}
	}
	return row
}

func NewTupleFromTable(ctx context.Context, db *pachsql.DB, tableName string) (Tuple, error) {
	colInfos, err := pachsql.GetTableColumns(ctx, db, tableName)
	if err != nil {
		return nil, err
	}
	tuple := make(Tuple, len(colInfos))
	for i := range colInfos {
		var err error
		tuple[i], err = makeTupleElement(colInfos[i].DataType, colInfos[i].IsNullable)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func makeTupleElement(dbType string, nullable bool) (interface{}, error) {
	switch dbType {
	case "BOOL":
		if nullable {
			return &sql.NullBool{}, nil
		} else {
			ret := false
			return &ret, nil
		}
	case "INT":
		if nullable {
			return &sql.NullInt32{}, nil
		} else {
			ret := false
			return &ret, nil
		}
	case "BIGINT":
		if nullable {
			return &sql.NullInt64{}, nil
		} else {
			ret := false
			return &ret, nil
		}
	case "FLOAT":
		if nullable {
			return &sql.NullFloat64{}, nil
		} else {
			ret := false
			return &ret, nil
		}
	case "VARCHAR", "TEXT":
		if nullable {
			return &sql.NullString{}, nil
		} else {
			ret := ""
			return &ret, nil
		}
	default:
		return nil, errors.Errorf("unrecognized type: %v", dbType)
	}
}
