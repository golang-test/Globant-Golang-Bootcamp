package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	userName = "mysql"
	password = "mysql"
	ip       = "192.168.99.100"
	port     = "3306"
	dbName   = "bookstore_db"
)

func InitDB() *sql.DB {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	db, err := sql.Open("mysql", path)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
