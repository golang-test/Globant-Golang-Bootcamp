package service_interface_go

import "context"

type Book struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Genre  int    `json:"genre"`
	Amount int    `json:"amount"`
}

type BookWithGenre struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Genre  string `json:"genre"`
	Amount int    `json:"amount"`
}

type AllBooksRequest struct {
	Min_price int    `json:"min___price"`
	Max_price int    `json:"max___price"`
	Names     string `json:"names"`
	Genre     int    `json:"genre"`
}

type Database interface {
	CreateBook(ctx context.Context, book Book) (int, error)
	UpdateBook(ctx context.Context, id int, book Book) error
	GetBook(ctx context.Context, id int) (BookWithGenre, error)
	GetAllBooks(ctx context.Context, request AllBooksRequest) ([]Book, error)
	DeleteBook(ctx context.Context, id int) error
}
