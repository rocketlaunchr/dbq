// Copyright 2019-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package dbq_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/rocketlaunchr/dbq/v2"
)

var (
	db  *sql.DB
	ctx = context.Background()
	_   = spew.UnsafeDisabled
)

var (
	user   string = ""
	pword  string = ""
	host   string = ""
	port   string = ""
	dbname string = ""
)

// https://blog.golang.org/subtests
// https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html#benchmarking

func init() {
	db, _ = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pword, host, port, dbname))
	db.SetMaxOpenConns(1)
	err := db.Ping()
	if err != nil {
		panic(err)
	}
}

type modeldbq struct {
	ID    int    `dbq:"id"`
	Name  string `dbq:"name"`
	Email string `dbq:"email"`
}

func (m *modeldbq) ScanFast() []interface{} {
	return []interface{}{&m.ID, &m.Name, &m.Email}
}

type modelgorm struct {
	ID    int    `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

func (modelgorm) TableName() string {
	return "tests"
}

type modelsqlx struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func Benchmark(b *testing.B) {
	setup()
	defer cleanup()

	limits := []int{
		5,
		50,
		500,
		10000,
	}

	// Benchmark dbq
	for _, lim := range limits {
		lim := lim
		b.Run(fmt.Sprintf("dbq limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				res, err := dbq.Qs(ctx, db, q("tests", lim), modeldbq{}, nil)
				if err != nil {
					b.Fatal(err)
				}
				if len(res.([]*modeldbq)) != lim {
					panic("something is wrong")
				}
				// spew.Dump(res)
			}
		})
	}

	// Benchmark sqlx
	for _, lim := range limits {
		lim := lim
		b.Run(fmt.Sprintf("sqlx limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db := sqlx.NewDb(db, "mysql")

				res := []modelsqlx{}
				err := db.Select(&res, q("tests", lim))
				if err != nil {
					panic(err)
				}
				if len(res) != lim {
					panic("something is wrong")
				}
				// spew.Dump(res)
			}
		})
	}

	// Benchmark gorm
	for _, lim := range limits {
		lim := lim
		b.Run(fmt.Sprintf("gorm limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				g, err := gorm.Open("mysql", db)
				if err != nil {
					panic(err)
				}

				var res = []modelgorm{}

				err = g.Order("id").Limit(lim).Find(&res).Error
				if err != nil {
					panic(err)
				}
				if len(res) != lim {
					panic("something is wrong")
				}
				// spew.Dump(res)
			}
		})
	}
}

func q(table string, limit int) string {
	return fmt.Sprintf("SELECT id, name, email FROM %s ORDER BY id LIMIT %d", table, limit)
}

func setup() {
	// Create table
	createQ := `
	CREATE TABLE tests (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(50) NOT NULL DEFAULT '',
		email varchar(150) NOT NULL DEFAULT '',
		PRIMARY KEY (id)
	)`

	_, err := db.Exec(createQ)
	if err != nil {
		panic(err)
	}

	// Add 10,000 fake entries
	entries := []interface{}{}
	for i := 0; i < 10000; i++ {
		entry := []interface{}{
			i + 1,
			gofakeit.Name(),
			gofakeit.Email(),
		}
		entries = append(entries, entry)
	}
	stmt := dbq.INSERTStmt("tests", []string{"id", "name", "email"}, len(entries))
	_, err = dbq.E(ctx, db, stmt, nil, entries)
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	_, err := db.Exec(`DROP TABLE tests`)
	if err != nil {
		panic(err)
	}
}
