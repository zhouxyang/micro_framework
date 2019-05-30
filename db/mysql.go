package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// MyDB implements log sql
type MyDB struct {
	*sqlx.DB
}

// Query implements log sql
func (my *MyDB) Query(log *logrus.Entry, query string, args ...interface{}) (*sql.Rows, error) {
	defer timeit(logSQL(log, "sql: %v, args: %v", query, args))
	return my.DB.Query(query, args...)
}

// QueryRow implements log sql
func (my *MyDB) QueryRow(log *logrus.Entry, query string, args ...interface{}) *sql.Row {
	defer timeit(logSQL(log, "sql: %v, args: %v", query, args))
	return my.DB.QueryRow(query, args...)
}

// Exec implements log sql
func (my *MyDB) Exec(log *logrus.Entry, query string, args ...interface{}) (sql.Result, error) {
	defer timeit(logSQL(log, "sql: %v, args: %v", query, args))
	return my.DB.Exec(query, args...)
}

// Get sqlx's Get(), implements log sql
func (my *MyDB) Get(log *logrus.Entry, dest interface{}, query string, args ...interface{}) error {
	defer timeit(logSQL(log, "sql: %v, args: %v", query, args))
	return my.DB.Get(dest, query, args...)
}

// Select sqlx's Select(), implements log sql
func (my *MyDB) Select(log *logrus.Entry, dest interface{}, query string, args ...interface{}) error {
	defer timeit(logSQL(log, "sql: %v, args: %v", query, args))
	return my.DB.Select(dest, query, args...)
}

// timeit log f's execute time
func timeit(f func(t time.Time)) {
	endTime := time.Now()
	f(endTime)
}

// logSQL log sql execute time
func logSQL(log *logrus.Entry, format string, args ...interface{}) func(t time.Time) {
	startTime := time.Now()
	logSQL := fmt.Sprintf(format, args...)

	return func(t time.Time) {
		log.Infof("%v, start: %v, end: %v, elasped: %v", logSQL, startTime, t, t.Sub(startTime))
	}
}

// AddSQLComment 增加sql注释
func AddSQLComment(comment string, query string) string {
	return fmt.Sprintf("/* grpc_framework:%s */ ", comment) + query
}
