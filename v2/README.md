<p align="right">
  <a href="http://godoc.org/github.com/rocketlaunchr/dbq/v2"><img src="http://godoc.org/github.com/rocketlaunchr/dbq?status.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/rocketlaunchr/dbq"><img src="https://goreportcard.com/badge/github.com/rocketlaunchr/dbq" /></a>
</p>

<p align="center">
<img src="https://github.com/rocketlaunchr/dbq/raw/master/logo.png" alt="dbq" />
</p>

(Now compatible with MySQL and PostgreSQL!)

Everyone knows that performing simple **DATABASE queries** in Go takes numerous lines of code that is often repetitive. If you want to avoid the cruft, you have two options: A heavy-duty ORM that is not up to the standard of Laraval or Django. Or DBQ!

⚠️ **WARNING: You will seriously reduce your database code to a few lines**

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
- Automatically retry query with exponential backoff if operation fails
- Transaction management (automatic rollback)

## Dependencies

- [MySQL driver](https://github.com/go-sql-driver/mysql) OR
- [PostgreSQL driver](https://github.com/lib/pq)

**NOTE:** For mysql driver, [`parseTime=true`](https://github.com/go-sql-driver/mysql#parsetime) setting can interfere with unmarshaling to [`civil.*`](https://pkg.go.dev/cloud.google.com/go/civil?tab=doc) types.

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

### Query

[`Q`](https://godoc.org/github.com/rocketlaunchr/dbq/v2#Q) ordinarily returns `[]map[string]interface{}` results, but you can automatically
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
results, err := dbq.Qs(ctx, db, "SELECT * FROM users", user{}, nil)

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

stmt := dbq.INSERTStmt("users", []string{"name", "age", "created_at"}, len(users))

dbq.E(ctx, db, stmt, nil, users)

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

**NOTE:** [FlattenArgs](https://godoc.org/github.com/rocketlaunchr/dbq/v2#FlattenArgs) function can be used more generally.

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

After fetching the results, you can further modify the results by implementing the [`PostUnmarshaler`](https://godoc.org/github.com/rocketlaunchr/dbq/v2#PostUnmarshaler) interface. The `PostUnmarshal` function must be attached to the pointer of the struct.

```go
type user struct {
  ID        int       `dbq:"id"`
  Name      string    `dbq:"name"`
  Age       int       `dbq:"age"`
  CreatedAt time.Time `dbq:"created_at"`
  HashedID  string    `dbq:"-"`          // Obfuscate ID
}

func (u *user) PostUnmarshal(ctx context.Context, row, total int) error {
  u.HashedID = obfuscate(u.ID)
  return nil
}
```

### ScanFaster

The [`ScanFaster`](https://godoc.org/github.com/rocketlaunchr/dbq/v2#ScanFaster) interface eradicates the use of the reflect package when unmarshaling. If you don't need to perform fancy time conversions or interpret weakly typed data, then it is more performant.

```go
type user struct {
  ID       int    `dbq:"id"`
  Name     string `dbq:"name"`
}

func (u *user) ScanFast() []interface{} {
  return []interface{}{&u.ID, &u.Name}
}
```

### Retry with Exponential Backoff

If the database operation fails, you can automatically retry with exponentially increasing intervals between each retry attempt. You can also set the maximum number of retries.

```go
opts := &dbq.Options{
  RetryPolicy:  dbq.ExponentialRetryPolicy(60 * time.Second, 3),
}
```

### Transaction Management

You can conveniently perform numerous complex database operations within a transaction without having to worry about rolling back. Unless you explicitly commit, it will automatically rollback.

You have access to the `Q` and `E` function as well as the underlying `tx` for performance purposes.

```go
ctx := context.Background()
pool, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db")

dbq.Tx(ctx, pool, func(tx interface{}, Q dbq.QFn, E dbq.EFn, txCommit dbq.TxCommit) {
  
  stmt := dbq.INSERTStmt("table", []string{"name", "age", "created_at"}, 1)
  res, err := E(ctx, stmt, nil, "test name", 34, time.Now())
  if err != nil {
    return // Automatic rollback
  }
  txCommit() // Commit
})
```

## Custom Queries

The `v2/x` subpackage will house functions to perform custom SQL queries. If they are general to both MySQL and PostgreSQL, they are inside the `x` subpackage. If they are specific to MySQL xor PostgreSQL, they are in the `x/mysql` xor `x/pg` subpackage respectively.


### This is your package too!

If you want your own custom functions included, just submit a PR and place it in your **own directory** inside `v2/x`. As long as it compiles and is well documented it is welcome.

### Bulk Update

As a warmup, I have included a [Bulk Update](https://godoc.org/github.com/rocketlaunchr/dbq/v2/x#BulkUpdate) function that works with MySQL and PostgreSQL. It allows you to update thousands of rows in 1 query without a transaction!

## Other useful packages

- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go) - Statistics and data manipulation
- [electron-alert](https://github.com/rocketlaunchr/electron-alert) - SweetAlert2 for Electron Applications
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
