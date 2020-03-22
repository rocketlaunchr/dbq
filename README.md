# dbq - Barbeque the boilerplate code [![GoDoc](http://godoc.org/github.com/rocketlaunchr/dbq?status.svg)](https://godoc.org/github.com/rocketlaunchr/dbq/v2) [![Go Report Card](https://goreportcard.com/badge/github.com/rocketlaunchr/dbq)](https://goreportcard.com/report/github.com/rocketlaunchr/dbq)

<p align="center">
<img src="https://github.com/rocketlaunchr/dbq/raw/master/logo.png" alt="dbq" />
</p>

(Now compatible with MySQL and PostgreSQL!)

Everyone knows that performing simple **DATABASE queries** in Go takes numerous lines of code that is often repetitive. If you want to avoid the cruft, you have two options: A heavy-duty ORM that is not up to the standard of Laraval or Django. Or DBQ!

**WARNING: You will seriously reduce your database code to a few lines**

⭐ **the project to show your appreciation.**

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

**NOTE:** For mysql driver, `parseTime=true` setting can interfere with unmarshaling to `civil.*` types.

## Installation

```
go get -u github.com/rocketlaunchr/dbq/v2
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

type Row struct {
  Name      string
  Age       int
  CreatedAt time.Time
}

users := []interface{}{
  dbq.Struct(Row{"Brad", 45, time.Now()}),
  dbq.Struct(Row{"Ange", 36, time.Now()}),
  dbq.Struct(Row{"Emily", 22, time.Now()}),
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
([]*main.user) (len=6 cap=8) {
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

### Flatten Query Args

All slices are flattened automatically.

```go
args1 := []string{"A", "B", "C"}
args2 := []interface{}{2, "D"}
args3 := dbq.Struct(Row{"Brad Pitt", 45, time.Now()})

results := dbq.MustQ(ctx, db, stmt, args1, args2, args3)

// Placeholder arguments will get flattened to:
results := dbq.MustQ(ctx, db, stmt, "A", "B", "C", 2, "D", "Brad Pitt", 45, time.Now())

```

### MySQL cancelation

To properly cancel a MySQL query, you need to use the [mysql-go](https://github.com/rocketlaunchr/mysql-go) package. `dbq` plays nicely with it.

```go
import sql "github.com/rocketlaunchr/mysql-go"

pool, _ := sql.Open("user:password@tcp(localhost:3306)/db")

conn, err := pool.Conn(ctx)

opts := &dbq.Options{
  SingleResult: true,
  PostFetch: func(ctx context.Context) error {
    return conn.Close()
  },
}

result := dbq.MustQ(ctx, conn, "SELECT * FROM users LIMIT 1", opts)
if result == nil {
  // no result
} else {
  result.(map[string]interface{})
}
```

### PostUnmarshaler

After fetching the results, you can further modify the results by implementing the `PostUnmarshaler` interface. The `PostUnmarshal` function must be attached to the pointer of the struct.

```go
type user struct {
  ID        int       `dbq:"id"`
  Name      string    `dbq:"name"`
  Age       int       `dbq:"age"`
  CreatedAt time.Time `dbq:"created_at"`
  HashedID  string    `dbq:"-"`          // Obfuscate ID
}

func (u *user) PostUnmarshal(ctx context.Context, row, count int) error {
  u.HashedID = obfuscate(u.ID)
  return nil
}
```

## Custom Queries

The `v2/x` subpackage will house functions to perform custom SQL queries. If you believe that a particular query is common or useful, submit a PR. They can be general to both MySQL and PostgreSQL or specific to either.

## Difference between v1 and v2

When a `ConcreteStruct` is provided, in `v1`, the `Q` and `MustQ` functions return `[]interface{}` while in `v2` they return `[]*struct`.

**NOTE:** `v1` is obsolete and will no longer receive updates.

## Other useful packages

- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go) - Statistics and data manipulation
- [igo](https://github.com/rocketlaunchr/igo) - A Go transpiler with cool new syntax such as fordefer (defer for for-loops)
- [mysql-go](https://github.com/rocketlaunchr/mysql-go) - Properly cancel slow MySQL queries
- [react](https://github.com/rocketlaunchr/react) - Build front end applications using Go
- [remember-go](https://github.com/rocketlaunchr/remember-go) - Cache slow database queries

#

### Legal Information

The license is a modified MIT license. Refer to the `LICENSE` file for more details.

**© 2019-20 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests. Note that the project is written in [igo](https://github.com/rocketlaunchr/igo) and transpiled into Go.
