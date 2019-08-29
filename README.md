# dbq - Barbeque the boilerplate code [![GoDoc](http://godoc.org/github.com/rocketlaunchr/dbq?status.svg)](http://godoc.org/github.com/rocketlaunchr/dbq) [![Go Report Card](https://goreportcard.com/badge/github.com/rocketlaunchr/dbq)](https://goreportcard.com/report/github.com/rocketlaunchr/dbq)

<p align="center">
<img src="https://github.com/rocketlaunchr/dbq/raw/master/logo.png" alt="dbq" />
</p>

(Now compatible with MySQL and PostgreSQL!)

Everyone knows that performing simple **DATABASE queries** in Go takes numerous lines of code that is often repetitive. If you want to avoid the cruft, you have two options: A heavy-duty ORM that is not up to the standard of Laraval or Django. Or DBQ!

**WARNING: You will seriously reduce your database code to a few lines**

## What is included

- Supports ANY type of query
- **MySQL** and **PostgreSQL** compatible
- **Convenient** and **Developer Friendly**
- Accepts any type of slice for query args
- Flattens query arg slices to individual values
- Bulk Insert seamlessly
- Automatically unmarshal query results directly to a struct using [mapstructure](https://github.com/mitchellh/mapstructure) package
- Lightweight
- Compatible with [mysql-go](https://github.com/rocketlaunchr/mysql-go) for proper MySQL query cancelation

## Dependencies

- [MySQL driver](https://github.com/go-sql-driver/mysql) OR
- [PostgreSQL driver](https://github.com/lib/pq)

## Installation

```
go get -u github.com/rocketlaunchr/dbq
```

## Examples

Let's assume a table called `users`:

| id  | name  | age | created_at |
| --- | ----- | --- | ---------- |
| 1   | Sally | 12  | 2019-03-01 |
| 2   | Peter | 15  | 2019-02-01 |
| 3   | Tom   | 18  | 2019-01-01 |

### Bulk Insert

You can insert multiple rows at once.

```go

db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db")

users := []interface{}{
  []interface{}{"Brad", 45, time.Now()},
  []interface{}{"Ange", 36, time.Now()},
  []interface{}{"Emily", 22, time.Now()},
}

stmt := dbq.INSERT("users", []string{"name", "age", "created_at"}, len(users))

dbq.E(ctx, db, stmt, nil, users)

```

### Query

`dbq.Q` ordinarily returns `[]map[string]interface{}` results, but you can automatically
unmarshal to a struct. You will need to type assert the results.

```go

type user struct {
  ID        int       `dbq:"id"`
  Name      string    `dbq:"name"`
  Age       int       `dbq:"age"`
  CreatedAt time.Time `dbq:"created_at"`
}

opts := &dbq.Options{ConcreteStruct: user{}, DecoderConfig:x}

results, err := dbq.Q(ctx, db, "SELECT * FROM users", opts)

```

Results:

```groovy
([]interface {}) (len=6 cap=8) {
 (*main.user)(0xc00009e1c0)({
  ID: (int) 1,
  Name: (string) (len=5) "Sally",
  Age: (int) 12,
  CreatedAt: (time.Time) 2019-03-01 00:00:00 +0000 UTC
 }),
 (*main.user)(0xc00009e300)({
  ID: (int) 2,
  Name: (string) (len=5) "Peter",
  Age: (int) 15,
  CreatedAt: (time.Time) 2019-02-01 00:00:00 +0000 UTC
 }),
 (*main.user)(0xc00009e440)({
  ID: (int) 3,
  Name: (string) (len=3) "Tom",
  Age: (int) 18,
  CreatedAt: (time.Time) 2019-01-01 00:00:00 +0000 UTC
 }),
 (*main.user)(0xc00009e5c0)({
  ID: (int) 4,
  Name: (string) (len=4) "Brad",
  Age: (int) 45,
  CreatedAt: (time.Time) 2019-07-24 14:36:58 +0000 UTC
 }),
 (*main.user)(0xc00009e700)({
  ID: (int) 5,
  Name: (string) (len=4) "Ange",
  Age: (int) 36,
  CreatedAt: (time.Time) 2019-07-24 14:36:58 +0000 UTC
 }),
 (*main.user)(0xc00009e840)({
  ID: (int) 6,
  Name: (string) (len=5) "Emily",
  Age: (int) 22,
  CreatedAt: (time.Time) 2019-07-24 14:36:58 +0000 UTC
 })
}
```

### Query Single Row

If you know that the query will return at maximum 1 row:

```go
result := dbq.MustQ(ctx, db, "SELECT * FROM users LIMIT 1", dbq.SingleResult)
if result == nil {
  // no result
} else {
  result.(map[string]interface{})
}

```

### MySQL cancelation

To properly cancel a MySQL query, you need to use the [mysql-go](https://github.com/rocketlaunchr/mysql-go) package. `dbq` plays nicely with it.

```go
import (
   stdSql "database/sql"
   sql "github.com/rocketlaunchr/mysql-go"
)

p, _ := stdSql.Open("mysql", "user:password@tcp(localhost:3306)/db")
kP, _ := stdSql.Open("mysql", "user:password@tcp(localhost:3306)/db")
kP.SetMaxOpenConns(1)

pool := &sql.DB{p, kP}

conn, err := pool.Conn(ctx)
defer conn.Close()

result := dbq.MustQ(ctx, conn, "SELECT * FROM users LIMIT 1", dbq.SingleResult)
if result == nil {
  // no result
} else {
  result.(map[string]interface{})
}
```

## Other useful packages

- [remember-go](https://github.com/rocketlaunchr/remember-go) - Cache slow database queries
- [mysql-go](https://github.com/rocketlaunchr/mysql-go) - Properly cancel slow MySQL queries
- [react](https://github.com/rocketlaunchr/react) - Build front end applications using Go
- [igo](https://github.com/rocketlaunchr/igo) - A Go transpiler with cool new syntax such as fordefer (defer for for-loops)
- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go) - Statistics and data manipulation.

#

### Legal Information

The license is a modified MIT license. Refer to the `LICENSE` file for more details.

**© 2019 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests. Note that the project is written in [igo](https://github.com/rocketlaunchr/igo) and transpiled into Go.

**Star** the project to show your appreciation.
