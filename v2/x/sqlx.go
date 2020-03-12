package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rocketlaunchr/dbq/v2"
)

type res struct{}

func (*res) LastInsertId() (int64, error) {
	return 0, nil
}

func (*res) RowsAffected() (int64, error) {
	return 0, nil
}

// BulkUpdateOptions is used to Set all required Options for BulkUpdatex Method
type BulkUpdateOptions struct {
	Table      string
	Columns    []string
	PrimaryKey string
	StmtSuffix string
	// if set as nil, default is dbq.MySQL
	DBType *dbq.Database
}

// BulkUpdate is used to bulk update multiple columns in a table.
// updateData's key must be the primary key in the table, and it's value is a slice of update values. The update value can be null.
// WARNING: updateData's value (which is a slice) must have the same length as columns.
//
// See: http://blog.bubble.ro/how-to-make-multiple-updates-using-a-single-query-in-mysql/
func BulkUpdate(ctx context.Context, db dbq.ExecContexter, updateData map[interface{}][]interface{}, opts BulkUpdateOptions) (sql.Result, error) {

	var (
		tableName        string
		columns          []string
		dbtype           dbq.Database
		primaryKeyColumn string
		extra            string
	)

	if db == nil || opts.Table == "" || len(opts.Columns) == 0 {
		return nil, errors.New("table is empty or no columns specified")
	}
	tableName = opts.Table
	columns = opts.Columns

	if opts.DBType != nil {
		dbtype = *opts.DBType
	}

	if len(updateData) == 0 {
		return &res{}, nil
	}

	if opts.PrimaryKey == "" {
		return nil, errors.New("primary key column in database table needs to be specified")
	}
	primaryKeyColumn = opts.PrimaryKey

	extra = opts.StmtSuffix

	queryArgs := []interface{}{}

	sqlUpdate := fmt.Sprintf("UPDATE %s\n SET\n", tableName)
	sqlUpdateBack := "\nWHERE " + primaryKeyColumn + " IN (%s)"

	//Generate query
	var primaryKeys []interface{}

	var phIdx int

	for j, field := range columns {

		eachSet := fmt.Sprintf("%s = CASE\n", field)

		for primaryKey, val := range updateData {
			if j == 0 {
				primaryKeys = append(primaryKeys, primaryKey)
			}

			if val[j] == nil {
				if dbtype == dbq.PostgreSQL {

					eachSet = eachSet + fmt.Sprintf("WHEN %v = $%d THEN NULL\n", primaryKeyColumn, phIdx+1)
					phIdx++
				} else {
					eachSet = eachSet + fmt.Sprintf("WHEN %v = ? THEN NULL\n", primaryKeyColumn)
				}

				queryArgs = append(queryArgs, primaryKey)
			} else {
				var v interface{}

				if reflect.ValueOf(val[j]).Kind() == reflect.Ptr {
					if reflect.ValueOf(val[j]).IsNil() {
						v = nil
					} else {
						v = reflect.ValueOf(val[j]).Elem().Interface()
					}
				} else {
					v = val[j]
				}

				if dbtype == dbq.PostgreSQL {
					eachSet = eachSet + fmt.Sprintf("WHEN %v = $%d THEN $%d\n", primaryKeyColumn, phIdx+1, phIdx+2)
					phIdx += 2
				} else {
					eachSet = eachSet + fmt.Sprintf("WHEN %v = ? THEN ?\n", primaryKeyColumn)
				}
				queryArgs = append(queryArgs, primaryKey, v)
			}
		}

		eachSet = eachSet + "END,\n"

		sqlUpdate = fmt.Sprintf("%s %s", sqlUpdate, eachSet)
	}
	sqlUpdate = strings.TrimSuffix(sqlUpdate, ",\n")

	stmt := sqlUpdate + fmt.Sprintf(sqlUpdateBack, dbq.Ph(1, len(primaryKeys), phIdx, dbtype))
	if extra != "" {
		stmt = stmt + " " + extra
	}

	fmt.Println(stmt)

	queryArgs = append(queryArgs, primaryKeys...)

	return db.ExecContext(ctx, stmt, queryArgs...)
}
