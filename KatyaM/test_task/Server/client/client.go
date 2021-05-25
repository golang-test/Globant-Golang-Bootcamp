package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	b "test_task/Server/database_interface"
	errors2 "test_task/Server/errors"
)

type BookClient struct {
	l        *zap.Logger
	endpoint string
}

func NewBookClient(l *zap.Logger, endpoint string) *BookClient {
	return &BookClient{l: l, endpoint: endpoint}
}

func (c *BookClient) Create(ctx context.Context, name string, price int, genre int, amount int) (int, error) {
	book := b.Book{Name: name, Price: price, Genre: genre, Amount: amount}
	book_json, err := json.Marshal(book)
	if err != nil {
		return 0, errors2.ErrWhileMarshal
	}
	body := bytes.NewBufferString(string(book_json))
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/create", body)
	rq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(rq)
	r, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(r), "error") {
		return 0, errors.New(string(r))
	}
	var id int
	err = json.Unmarshal(r, &id)
	c.l.Info("resp : " + string(r))
	c.l.Info("get response")
	if err != nil {
		return 0, err
	}
	return id, err
}

func (c *BookClient) Update(ctx context.Context, id int, name string, price int, genre int, amount int) error {
	book := b.Book{Name: name, Price: price, Genre: genre, Amount: amount}
	book_json, err := json.Marshal(book)
	body := bytes.NewBufferString(string(book_json))
	s := strconv.Itoa(id)
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/update?id="+s, body)
	rq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return errors2.ErrWhileUpdate
	}
	r, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(r), "error") {
		return errors.New(string(r))
	}
	return err
}

func (c *BookClient) Delete(ctx context.Context, id int) error {
	s := strconv.Itoa(id)
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/delete?id="+s, *new(io.Reader))
	rq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(rq)
	r, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(r), "error") {
		return errors.New(string(r))
	}
	return err
}

func (c *BookClient) GetBook(ctx context.Context, id int) (b.BookWithGenre, error) {
	s := strconv.Itoa(id)
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/get_book?id="+s, *new(io.Reader))
	if err != nil {
		return b.BookWithGenre{}, err
	}
	rq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(rq)
	r, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(r), "error") {
		return b.BookWithGenre{}, errors.New(string(r))
	}
	var book b.BookWithGenre
	err = json.Unmarshal(r, &book)
	return book, err
}

func (c *BookClient) GetAllBooks(ctx context.Context, request b.AllBooksRequest) ([]b.Book, error) {
	req_json, err := json.Marshal(request)
	body := bytes.NewBufferString(string(req_json))
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/get_all?", body)
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(rq)
	r, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(r), "error") {
		return nil, errors.New(string(r))
	}
	var books []b.Book
	err = json.Unmarshal(r, &books)
	return books, err
}
