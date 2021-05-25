// +build !solution

package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	b "test_task/Database"
)

type BookService struct {
	L  *zap.Logger
	db *b.BookDatabase
}

func (s BookService) Initialize(username, password, database string) (b.BookDatabase, error) {
	db := b.BookDatabase{}
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		username, password, database)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	s.db = &db
	return db, nil
}
