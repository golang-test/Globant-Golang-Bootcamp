package repository

import (
	"database/sql"
	"fmt"

	"github.com/DmitriyZhevnov/library/src/entities"
)

type repo struct{}

func NewPostgresBookRepository() BookRepository {
	return &repo{}
}

func (*repo) FindAll(db *sql.DB) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book")
	return buildBooks(rows, err)
}

func (*repo) FindById(db *sql.DB, id int) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where id = '%d'", id)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func (*repo) FindByName(db *sql.DB, name string) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where name = '%s'", name)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func (*repo) FilterByGenre(db *sql.DB, id int) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where genre_id = '%d'", id)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func (*repo) FilterByPrices(db *sql.DB, min, max float64) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where price >= '%f' AND price <= '%f'", min, max)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func buildBooks(rows *sql.Rows, er error) (book []entities.Book, err error) {
	if er != nil {
		return nil, er
	} else {
		defer rows.Close()
		var books []entities.Book
		for rows.Next() {
			var b entities.Book
			if err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.GenreId, &b.Amount); err != nil {
				return nil, err
			}
			if b.Amount > 0 {
				books = append(books, b)
			}
		}
		return books, nil
	}
}

func (*repo) Create(db *sql.DB, book *entities.Book) error {
	id := 0
	err := db.QueryRow("insert into book(name, price, genre_id, amount) values ($1, $2, $3, $4) RETURNING book.id", book.Name, book.Price, book.GenreId, book.Amount).Scan(&id)
	if err != nil {
		return err
	} else {
		book.Id = id
		return nil
	}
}

func (*repo) Update(db *sql.DB, id int64, book *entities.Book) (err error) {
	idUpdatedBook := 0
	err = db.QueryRow("update book set name = $1, price = $2, genre_id = $3, amount = $4 where id = $5 RETURNING book.id",
		book.Name, book.Price, book.GenreId, book.Amount, id).Scan(&idUpdatedBook)
	return
}

func (*repo) Delete(db *sql.DB, id int64) (int64, error) {
	sqlRequest := fmt.Sprintf("delete from book where id = '%d'", id)
	result, err := db.Exec(sqlRequest)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
