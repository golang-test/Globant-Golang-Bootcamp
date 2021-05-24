package repository

import (
	"fmt"
	"github.com/SavinskiPavel/bookstore/internal/database"
	"github.com/SavinskiPavel/bookstore/internal/model"
	"log"
	"strconv"
)

func FindAllBooks() []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Println(err)
	}
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	defer rows.Close()
	return books
}

func SaveBook(name string, price float64, genre int, amount int) (int, error) {
	db := database.InitDB()
	result, err := db.Exec(fmt.Sprintf("INSERT INTO books (name, price, genre, amount) VALUES ('%s', %f, %d, %d);", name, price, genre, amount))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func FindBookById(id string) model.Book {
	db := database.InitDB()
	row, err := db.Query(fmt.Sprintf("SELECT * FROM books WHERE id = '%s'", id))
	if err != nil {
		log.Println(err)
	}
	defer row.Close()
	var b model.Book
	for row.Next() {
		err = row.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
		}
	}
	return b
}

func UpdateBook(id string, name string, price float64, genre int, amount int) {
	db := database.InitDB()
	integerId, _ := strconv.Atoi(id)
	ins, err := db.Query(fmt.Sprintf("UPDATE books SET name='%s', price=%f, genre=%d, amount=%d WHERE id=%d", name, price, genre, amount, integerId))
	if err != nil {
		log.Println(err)
	}
	defer ins.Close()
}

func DeleteBook(id string) {
	db := database.InitDB()
	integerId, _ := strconv.Atoi(id)
	del, err := db.Query(fmt.Sprintf("DELETE FROM books WHERE id=%d", integerId))
	if err != nil {
		log.Println(err)
	}
	defer del.Close()
}

func FindBooksByGenre(genre int) []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books WHERE genre=?", genre)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}

func FindBooksById(id int) []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books WHERE id=?", id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}

func FindBooksByName(name string) []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books WHERE name=?", name)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}

func FindBooksByPrices(min, max float64) []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books WHERE price >= ? AND price <= ?", min, max)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}

func FindBooksByMaxPrice(max float64) []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books WHERE price <= ?", max)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}

func FindBooksByMinPrice(min float64) []model.Book {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM books WHERE price >= ?", min)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		b := model.Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.Genre, &b.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}
