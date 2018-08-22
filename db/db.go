package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db    *sql.DB
	dbURI string
	err   error
)

func init() {
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
