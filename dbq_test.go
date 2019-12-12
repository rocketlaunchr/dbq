// Copyright 2019 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package dbq

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/mapstructure"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type store struct {
	ID        int64     `dbq:"id"`
	Product   string    `dbq:"product"`
	Price     float64   `dbq:"price"`
	Quantity  int64     `dbq:"quantity"`
	Available int64     `dbq:"available"`
	DateAdded time.Time `dbq:"date_added"`
}

func (s *store) PostUnmarshal(ctx context.Context, row, count int) error {
	// This postUnmarshall method changes the timezone on DateAdded in Store struct
	// From UTC to CEST (Europe/Budapest)

	loc, err := time.LoadLocation("Europe/Budapest")
	if err != nil {
		return err
	}
	newTimeZone := s.DateAdded.In(loc)
	s.DateAdded = newTimeZone

	return nil
}

func TestMustQ(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// tRef := "2006-01-02 15:04:05"
	tRef := time.Now()

	rows := sqlmock.NewRows([]string{"id", "product", "price", "quantity", "available", "date_added"}).
		AddRow(int64(1), "wrist watch", float64(45000.98), int64(6), int64(1), tRef).
		AddRow(int64(2), "bags", float64(25089.55), int64(10), int64(0), tRef).
		AddRow(int64(3), "car", float64(598000999.99), int64(3), int64(1), tRef)

	expected := []interface{}{
		&store{
			ID:        1,
			Product:   "wrist watch",
			Price:     float64(45000.98),
			Quantity:  int64(6),
			Available: int64(1),
			DateAdded: tRef,
		},
		&store{
			ID:        2,
			Product:   "bags",
			Price:     float64(25089.55),
			Quantity:  int64(10),
			Available: int64(0),
			DateAdded: tRef,
		},
		&store{
			ID:        3,
			Product:   "car",
			Price:     float64(598000999.99),
			Quantity:  int64(3),
			Available: int64(1),
			DateAdded: tRef,
		},
	}

	row := sqlmock.NewRows([]string{"id", "product", "price", "quantity", "available", "date_added"}).
		AddRow(int64(1), "wrist watch", float64(45000.98), int64(6), int64(1), tRef)

	mock.ExpectQuery("^SELECT (.+) FROM store$").WillReturnRows(rows) // Multiple result select query

	mock.ExpectQuery("^SELECT (.+) FROM store.*$").WillReturnRows(sqlmock.NewRows(nil)) // zero result

	mock.ExpectQuery("^SELECT (.+) FROM store LIMIT 1$").WillReturnRows(row) // single result

	ctx := context.Background()

	// Testing Multiple Data select with MustQ
	opts := &Options{ConcreteStruct: store{}, DecoderConfig: &StructorConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true}}

	actual := MustQ(ctx, db, "SELECT * FROM store", opts)

	if !cmp.Equal(expected, actual) {
		t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
	}

	// Test zero data select query with MustQ

	_, err = Q(ctx, db, "SELECT * FROM store WHERE id = 20", opts)
	if err != nil {
		t.Errorf("There was an error while executing statement: %s", err)
	}

	// Testing Single Data Select
	opts2 := &Options{ConcreteStruct: store{}, DecoderConfig: &StructorConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true}, SingleResult: true}

	// Test Select return at most 1 result
	_ = MustQ(ctx, db, "SELECT * FROM store LIMIT 1", opts2)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestMustE(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// tRef := "2006-01-02 15:04:05"
	tRef := time.Now()

	mock.ExpectQuery("^SELECT (.+) FROM store$").
		WillReturnError(fmt.Errorf("There was error while executing statement"))

	mock.ExpectExec("INSERT INTO store").
		WithArgs(4, "mobile phone", 456787.45, 8, 1, tRef).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// This is for batch Insert with MySQL
	mock.ExpectExec("INSERT INTO store").
		WithArgs(
			int64(6), "Dish Washer", float64(45534.34), int64(34), int64(1), AnyTime{},
			int64(7), "Sewing Machine", float64(9843.35), int64(8), int64(0), AnyTime{},
			int64(8), "Private Jet", float64(98748594.34), int64(2), int64(1), AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// This is for batch insert with PostgreSQL
	mock.ExpectExec("INSERT INTO store").
		WithArgs(
			int64(6), "Dish Washer", float64(45534.34), int64(34), int64(1), AnyTime{},
			int64(7), "Sewing Machine", float64(9843.35), int64(8), int64(0), AnyTime{},
			int64(8), "Private Jet", float64(98748594.34), int64(2), int64(1), AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE store SET product").
		WithArgs("buckets", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM store").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	// Error testing on Select Query
	opts := &Options{ConcreteStruct: store{}, DecoderConfig: &StructorConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
		WeaklyTypedInput: true}}

	_, err = E(ctx, db, "SELECT * FROM store", opts)
	if err == nil {
		t.Errorf("was expecting an error, but there was none.")
	}

	// Testing Single Insert
	insertArgs := []interface{}{4, "mobile phone", 456787.45, 8, 1, tRef}

	_ = MustE(ctx, db, "INSERT INTO store(id, product, price, quantity, available, date_added) VALUES (?, ?, ?, ?, ?, ?)", nil, insertArgs)

	// Testing Batch Insert

	storeProducts := []interface{}{
		[]interface{}{int64(6), "Dish Washer", float64(45534.34), int64(34), int64(1), tRef},
		[]interface{}{int64(7), "Sewing Machine", float64(9843.35), int64(8), int64(0), tRef},
		[]interface{}{int64(8), "Private Jet", float64(98748594.34), int64(2), int64(1), tRef},
	}

	// batch insert statement on MySQL
	stmt := INSERT("store", []string{"id", "product", "price", "quantity", "available", "date_added"}, len(storeProducts), MySQL)

	_ = MustE(ctx, db, stmt, opts, storeProducts)

	// batch insert statement on PostgreSQL
	stmt2 := INSERT("store", []string{"id", "product", "price", "quantity", "available", "date_added"}, len(storeProducts), PostgreSQL)

	_ = MustE(ctx, db, stmt2, opts, storeProducts)

	// Testing Data update with MustE

	updateArgs := []interface{}{"buckets", 2}
	_ = MustE(ctx, db, "UPDATE store SET product = ? WHERE id = ?", nil, updateArgs)

	// Testing Delete from table store
	_ = MustE(ctx, db, "DELETE FROM store WHERE ID = ?", nil, []interface{}{int64(1)})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostUnmarshalConcurrent(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// tRef := "2006-01-02 15:04:05"
	tRef := time.Now()

	// convert tRef to new timezone
	loc, err := time.LoadLocation("Europe/Budapest")
	if err != nil {
		t.Errorf("an unexpected error occurred %s", err)
	}

	newTref := tRef.In(loc)

	rows := sqlmock.NewRows([]string{"id", "product", "price", "quantity", "available", "date_added"}).
		AddRow(int64(1), "wrist watch", float64(45000.98), int64(6), int64(1), tRef).
		AddRow(int64(2), "bags", float64(25089.55), int64(10), int64(0), tRef).
		AddRow(int64(3), "car", float64(598000999.99), int64(3), int64(1), tRef)

	expected := []interface{}{
		&store{
			ID:        1,
			Product:   "wrist watch",
			Price:     float64(45000.98),
			Quantity:  int64(6),
			Available: int64(1),
			DateAdded: newTref,
		},
		&store{
			ID:        2,
			Product:   "bags",
			Price:     float64(25089.55),
			Quantity:  int64(10),
			Available: int64(0),
			DateAdded: newTref,
		},
		&store{
			ID:        3,
			Product:   "car",
			Price:     float64(598000999.99),
			Quantity:  int64(3),
			Available: int64(1),
			DateAdded: newTref,
		},
	}

	mock.ExpectQuery("^SELECT (.+) FROM store$").WillReturnRows(rows) // Multiple result select query

	ctx := context.Background()

	// Testing Multiple Data select with MustQ
	opts := &Options{ConcreteStruct: store{}, DecoderConfig: &StructorConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true},
		ConcurrentPostUnmarshal: true}

	actual := MustQ(ctx, db, "SELECT * FROM store", opts)

	if !cmp.Equal(expected, actual) {
		t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
	}

}

func TestPostUnmarshalSequential(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// tRef := "2006-01-02 15:04:05"
	tRef := time.Now()

	// convert tRef to newtimezone
	loc, err := time.LoadLocation("Europe/Budapest")
	if err != nil {
		t.Errorf("an unexpected error occurred %s", err)
	}
	newTref := tRef.In(loc)

	rows := sqlmock.NewRows([]string{"id", "product", "price", "quantity", "available", "date_added"}).
		AddRow(int64(1), "wrist watch", float64(45000.98), int64(6), int64(1), tRef).
		AddRow(int64(2), "bags", float64(25089.55), int64(10), int64(0), tRef).
		AddRow(int64(3), "car", float64(598000999.99), int64(3), int64(1), tRef)

	expected := []interface{}{
		&store{
			ID:        1,
			Product:   "wrist watch",
			Price:     float64(45000.98),
			Quantity:  int64(6),
			Available: int64(1),
			DateAdded: newTref,
		},
		&store{
			ID:        2,
			Product:   "bags",
			Price:     float64(25089.55),
			Quantity:  int64(10),
			Available: int64(0),
			DateAdded: newTref,
		},
		&store{
			ID:        3,
			Product:   "car",
			Price:     float64(598000999.99),
			Quantity:  int64(3),
			Available: int64(1),
			DateAdded: newTref,
		},
	}

	mock.ExpectQuery("^SELECT (.+) FROM store$").WillReturnRows(rows) // Multiple result select query

	ctx := context.Background()

	// Testing Multiple Data select with MustQ
	opts := &Options{ConcreteStruct: store{}, DecoderConfig: &StructorConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true},
		ConcurrentPostUnmarshal: false}

	actual := MustQ(ctx, db, "SELECT * FROM store", opts)

	if !cmp.Equal(expected, actual) {
		t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
	}

}

func TestQueryRawResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// tRef := "2006-01-02 15:04:05"
	tRef := time.Now()

	rows := sqlmock.NewRows([]string{"id", "product", "price", "quantity", "available", "date_added"}).
		AddRow([]byte("1"), []byte("wrist watch"), []byte("45000.98"), []byte("6"), []byte("1"), tRef).
		AddRow([]byte("2"), []byte("bags"), []byte("25089.55"), []byte("10"), []byte("0"), tRef).
		AddRow([]byte("3"), []byte("car"), []byte("598000999.99"), []byte("3"), []byte("1"), tRef)

	// expected := []interface{}{
	// 	&store{
	// 		ID:        1,
	// 		Product:   "wrist watch",
	// 		Price:     float64(45000.98),
	// 		Quantity:  int64(6),
	// 		Available: int64(1),
	// 		DateAdded: tRef,
	// 	},
	// 	&store{
	// 		ID:        2,
	// 		Product:   "bags",
	// 		Price:     float64(25089.55),
	// 		Quantity:  int64(10),
	// 		Available: int64(0),
	// 		DateAdded: tRef,
	// 	},
	// 	&store{
	// 		ID:        3,
	// 		Product:   "car",
	// 		Price:     float64(598000999.99),
	// 		Quantity:  int64(3),
	// 		Available: int64(1),
	// 		DateAdded: tRef,
	// 	},
	// }

	mock.ExpectQuery("^SELECT (.+) FROM store$").WillReturnRows(rows) // Multiple result select query

	ctx := context.Background()

	// Testing Multiple Data select with MustQ
	opts := &Options{ConcreteStruct: store{}, DecoderConfig: &StructorConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		WeaklyTypedInput: true},
		RawResults: true}

	actual := MustQ(ctx, db, "SELECT * FROM store", opts)
	// spew.Dump(actual)
	_ = actual

	// if !cmp.Equal(expected, actual) {
	// 	t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
	// }
}
