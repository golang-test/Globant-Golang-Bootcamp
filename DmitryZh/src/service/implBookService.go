package service

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/DmitriyZhevnov/library/src/entities"
	"github.com/DmitriyZhevnov/library/src/repository"
	"gopkg.in/go-playground/validator.v9"
)

var repo repository.BookRepository

type service struct{}

func NewBookService(reposit repository.BookRepository) BookService {
	repo = reposit
	return &service{}
}

func (*service) Validate(book *entities.Book) error {
	validate := validator.New()
	err := validate.Struct(book)
	if err != nil {
		return errors.New("Invalid resquest payload")
	}
	return nil
}

func (*service) FindAll(db *sql.DB) (book []entities.Book, err error) {
	return repo.FindAll(db)
}

func (*service) FindById(db *sql.DB, id string) (book []entities.Book, err error) {
	idBook, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("Invalid book Id")
	}
	book, err = repo.FindById(db, idBook)
	if err != nil {
		return nil, err
	}
	return
}

func (*service) FindByName(db *sql.DB, name string) (book []entities.Book, err error) {
	return repo.FindByName(db, name)
}

func (*service) FilterByGenre(db *sql.DB, id string) (book []entities.Book, err error) {
	idGenre, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("Invalid book Id")
	}
	book, err = repo.FilterByGenre(db, idGenre)
	if err != nil {
		return nil, err
	}
	return
}

func (*service) FilterByPrices(db *sql.DB, min, max string) (book []entities.Book, err error) {
	minPrice, err := strconv.ParseFloat(min, 64)
	maxPrice, err := strconv.ParseFloat(max, 64)
	if err != nil {
		return nil, errors.New("Invalid prices")
	}
	book, err = repo.FilterByPrices(db, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	return
}

func (*service) Create(db *sql.DB, book *entities.Book) error {
	return repo.Create(db, book)
}

func (*service) Update(db *sql.DB, id string, book *entities.Book) (err error) {
	idBook, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("Invalid book Id")
	}
	book.Id = idBook
	if err = repo.Update(db, int64(idBook), book); err != nil {
		return errors.New("There is no book with this id")
	}
	return nil
}

func (*service) Delete(db *sql.DB, id string) (int64, error) {
	idBook, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("Invalid book Id")
	}
	var rowsAffected int64
	if rowsAffected, err = repo.Delete(db, int64(idBook)); err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, errors.New("There is no book with this id")
	}
	return 1, nil
}
