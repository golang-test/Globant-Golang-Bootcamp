package service

import (
	"database/sql"

	"github.com/DmitriyZhevnov/library/src/entities"
)

type BookService interface {
	Validate(book *entities.Book) error
	FindAll(db *sql.DB) (book []entities.Book, err error)
	FindById(db *sql.DB, id string) (book []entities.Book, err error)
	FindByName(db *sql.DB, name string) (book []entities.Book, err error)
	FilterByGenre(db *sql.DB, id string) (book []entities.Book, err error)
	FilterByPrices(db *sql.DB, min, max string) (book []entities.Book, err error)
	Create(db *sql.DB, book *entities.Book) error
	Update(db *sql.DB, id string, book *entities.Book) error
	Delete(db *sql.DB, id string) (int64, error)
}
