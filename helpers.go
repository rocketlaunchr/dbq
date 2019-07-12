// Copyright 2019 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package dbq

import (
	"fmt"
	"strings"
)

// Database is used to set the Database.
// Different databases have different syntax for placeholders etc.
type Database int

const (
	// MySQL database
	MySQL Database = 0
	// PostgreSQL database
	PostgreSQL Database = 1
)

// Ph generates the placeholders for SQL queries.
// For a bulk insert operation, rows is the number of rows you intend
// to insert, and fieldsN is the number of fields per row.
func Ph(fieldsN, rows int, dbtype ...Database) string {

	var typ Database
	if len(dbtype) > 0 {
		typ = dbtype[0]
	}

	if typ == MySQL {
		inner := "( " + strings.TrimSuffix(strings.Repeat("?,", fieldsN), ",") + " ),"
		return strings.TrimSuffix(strings.Repeat(inner, rows), ",")
	}

	var singleValuesStr string

	varCount := 1
	for i := 1; i <= rows; i++ {
		singleValuesStr = singleValuesStr + "("
		for j := 1; j <= fieldsN; j++ {
			singleValuesStr = singleValuesStr + fmt.Sprintf("$%d,", varCount)
			varCount++
		}
		singleValuesStr = strings.TrimSuffix(singleValuesStr, ",") + "),"
	}

	return strings.TrimSuffix(singleValuesStr, ",")
}
