package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	dbURI  string
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
}
