package repository

import (
	"database/sql"

	"github.com/DmitriyZhevnov/library/src/entities"
	_ "github.com/lib/pq"
)

type BookRepository interface {
	FindAll(db *sql.DB) (book []entities.Book, err error)
	FindById(db *sql.DB, id int) (book []entities.Book, err error)
	FindByName(db *sql.DB, name string) (book []entities.Book, err error)
	FilterByGenre(db *sql.DB, id int) (book []entities.Book, err error)
	FilterByPrices(db *sql.DB, min, max float64) (book []entities.Book, err error)
	Create(db *sql.DB, book *entities.Book) error
	Update(db *sql.DB, id int64, book *entities.Book) error
	Delete(db *sql.DB, id int64) (int64, error)
}
