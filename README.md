dbq - Barbeque the boiler plate code [![GoDoc](http://godoc.org/github.com/rocketlaunchr/dbq?status.svg)](http://godoc.org/github.com/rocketlaunchr/dbq) [![cover.run](https://cover.run/go/github.com/rocketlaunchr/dbq.svg?style=flat&tag=golang-1.12)](https://cover.run/go?tag=golang-1.12&repo=github.com%2Frocketlaunchr%2Fdbq) [![Go Report Card](https://goreportcard.com/badge/github.com/rocketlaunchr/dbq)](https://goreportcard.com/report/github.com/rocketlaunchr/dbq)

(Now compatible with MySQL and PostgreSQL!)
===============

Everyone knows that performing simple **DATABASE queries** in Go takes numerous lines of code that is often repetitive. If you want to avoid the clutter, you have two options: A heavy-duty ORM that is not up to the standard of Laraval or Django. Or DBQ!


**WARNING: You will seriously reduce your database code to a few lines**


## What is included

* Supports ANY type of query
* **MySQL** and **PostgreSQL** compatible
* **Convenient** and **Developer Friendly**
* Bulk Insert seamlessly
* Unmarshal query results directly to struct using [mapstructure](https://github.com/mitchellh/mapstructure) package
* Super lightweight
* Compatible with [mysql-go](https://github.com/rocketlaunchr/mysql-go) for proper MySQL cancelation

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

### Fordefer

See [Blog post](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01) on why this is an improvement. It can be especially helpful in unit tests.

```go

for {
	row, err := db.Query("SELECT ...")
	if err != nil {
		panic(err)
	}

	fordefer row.Close()
}

```


### Defer go

This feature makes Go's language syntax more internally consistent. There is no reason why `defer` and `go` should not work together.

```go

mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	// Transmit how long the request took to serve without delaying response to client.
	defer go transmitRequestStats(start)

	fmt.Fprintf(w, "Welcome to the home page!")
})

```


#

### Legal Information

The license is a modified MIT license. Refer to the `LICENSE` file for more details.

**© 2019 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests. Note that the project is written in [igo](https://github.com/rocketlaunchr/igo) and transpiled into Go.

**Star** the project to show your appreciation.
