package Database

import (
	"context"
	"database/sql"
	"fmt"
	b "test_task/Server/database_interface"
	"test_task/Server/errors"
)

var (
	username = "user"
	password = "123"
	database = "postgres"
)

type BookDatabase struct {
	Conn *sql.DB
}

func (db BookDatabase) CreateBook(ctx context.Context, book b.Book) (int, error) {
	last, err := db.Conn.QueryContext(ctx, "SELECT id FROM BookStore ORDER BY id DESC LIMIT 1")
	if err != nil {
		return 0, errors.ErrDBRequest
	}
	defer last.Close()
	var last_id int
	last.Next()
	if err := last.Scan(&last_id); err != nil {
		return 0, errors.ErrDBRequest
	}
	fmt.Sprint(last_id)
	check, err := db.Conn.QueryContext(ctx, "SELECT * FROM BookStore WHERE name = $1", book.Name)
	if err != nil {
		return 0, errors.ErrDBRequest
	}
	if check.Next() {
		return 0, errors.ErrSameName
	}
	_, err = db.Conn.ExecContext(
		ctx,
		"INSERT INTO BookStore VALUES ($1, $2, $3, $4, $5)", last_id+1, book.Name, book.Price, book.Amount, book.Genre)
	if err != nil {
		return 0, errors.ErrDBRequest
	}
	return last_id + 1, err
}

func (db BookDatabase) UpdateBook(ctx context.Context, id int, book b.Book) error {
	_, err := db.Conn.ExecContext(
		ctx,
		"UPDATE BookStore SET name = $1, price = $2, genre = $3, amount = $4 WHERE id = $5",
		book.Name, book.Price, book.Amount, book.Genre, id)
	if err != nil {
		return errors.ErrWhileUpdate
	}
	return err
}

func (db BookDatabase) GetBook(ctx context.Context, id int) (b.BookWithGenre, error) {
	row, err := db.Conn.QueryContext(ctx, "SELECT id, name, price, genre, amount, genre_name FROM BookStore book "+
		"LEFT JOIN Genres genre	ON book.genre = genre.genre_id WHERE book.id = $1, amount > 0", id)
	if err != nil {
		return b.BookWithGenre{}, errors.ErrDBRequest
	}
	defer row.Close()
	var book b.BookWithGenre
	var i int
	row.Next()
	err = row.Scan(&i, &book.Name, &book.Price, &book.Genre, &book.Amount)
	if err != nil {
		return b.BookWithGenre{}, errors.ErrNotFound
	}
	return book, err
}

func (db BookDatabase) DeleteBook(ctx context.Context, id int) error {
	_, err := db.Conn.ExecContext(ctx, "DELETE FROM BookStore WHERE id = $1", id)
	if err != nil {
		return errors.ErrDBRequest
	}
	return err
}

func (db BookDatabase) GetAllBooks(ctx context.Context, request b.AllBooksRequest) ([]b.Book, error) {
	if request.Min_price < request.Max_price || request.Max_price < 0 || request.Min_price < 0 {
		return nil, errors.ErrDBRequest
	}
	row, err := db.Conn.QueryContext(ctx, "SELECT * FROM BookStore Where (CASE WHEN $1 > 0 THEN price > $1"+
		"                                  WHEN $2 > 0 THEN price < $2 "+
		"								   WHEN $4 > 0  THEN genre = $4"+
		"								   WHEN $3 not is null THEN name = $2 ELSE true END), amount > 0",
		request.Min_price, request.Max_price, request.Names, request.Genre)
	if err != nil {
		return nil, errors.ErrDBRequest
	}
	defer row.Close()
	var allBooks []b.Book
	for row.Next() {
		var book b.Book
		var i int
		err = row.Scan(&i, &book.Name, &book.Price, &book.Genre, &book.Amount)
		allBooks = append(allBooks, book)
	}
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return allBooks, err
}
