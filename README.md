dbq - Barbeque the boiler plate code [![GoDoc](http://godoc.org/github.com/rocketlaunchr/dbq?status.svg)](http://godoc.org/github.com/rocketlaunchr/dbq) [![cover.run](https://cover.run/go/github.com/rocketlaunchr/dbq.svg?style=flat&tag=golang-1.12)](https://cover.run/go?tag=golang-1.12&repo=github.com%2Frocketlaunchr%2Fdbq) [![Go Report Card](https://goreportcard.com/badge/github.com/rocketlaunchr/dbq)](https://goreportcard.com/report/github.com/rocketlaunchr/dbq)
===============

(Now compatible with MySQL and PostgreSQL!)

Everyone knows that performing simple **DATABASE queries** in Go takes numerous lines of code that is often repetitive. If you want to avoid the clutter, you have two options: A heavy-duty ORM that is not up to the standard of Laraval or Django. Or DBQ!


**WARNING: You will seriously reduce your database code to a few lines**


## What is included

* Supports ANY type of query
* **MySQL** and **PostgreSQL** compatible
* **Convenient** and **Developer Friendly**
* Bulk Insert seamlessly
* Automatically Unmarshal query results directly to struct using [mapstructure](https://github.com/mitchellh/mapstructure) package
* Lightweight
* Compatible with [mysql-go](https://github.com/rocketlaunchr/mysql-go) for proper MySQL query cancelation

## Dependencies

* [MySQL driver](https://github.com/go-sql-driver/mysql) OR
* [PostgreSQL driver](https://github.com/lib/pq)


## Installation

```
go get -u github.com/rocketlaunchr/dbq
```


## Examples

Let's assume a table called `users`:

| id | name  | age | created_at |
|----|-------|-----|------------|
| 1  | Sally | 12  | 2019-03-01 |
| 2  | Peter | 15  | 2019-02-01 |
| 3  | Tom   | 18  | 2019-01-01 |


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

`dqb.Q` ordinarily returns `[]map[string]interface{}` results but you can automatically
unmarshal to a struct. You will need to type assert the results.

```go

db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db")

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

If you know that the query will return at max 1 row:

```go

db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db")
result := dbq.MustQ(ctx, db, "SELECT * FROM users LIMIT 1", dbq.SingleResult)
if result == nil {
	// no result
} else {
	result.(map[string]interface{})
}

```





#

### Legal Information

The license is a modified MIT license. Refer to the `LICENSE` file for more details.

**© 2019 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests. Note that the project is written in [igo](https://github.com/rocketlaunchr/igo) and transpiled into Go.

**Star** the project to show your appreciation.
