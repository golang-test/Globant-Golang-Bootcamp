package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AA55hex/golang_bootcamp/server/db/connection"
	"github.com/AA55hex/golang_bootcamp/server/db/entity"
	"github.com/gorilla/mux"
)

const (
	json_content = "application/json"
	text_plain   = "text/plain"
	no_content   = ""
)

// for simple json responsing
func jsonResponse(w http.ResponseWriter, httpStatus int, jsonBody interface{}) error {
	w.WriteHeader(httpStatus)
	w.Header().Add("Content-Type", json_content)
	err := json.NewEncoder(w).Encode(jsonBody)
	return err
}

// for simple text responsing
func textResponse(w http.ResponseWriter, httpStatus int, body []byte) {
	w.WriteHeader(httpStatus)
	w.Header().Add("Content-Type", text_plain)
	w.Write(body)
}

// validate Content-Type header
func validateContent(r *http.Request, content_type string) error {
	content := r.Header.Get("Content-Type")
	if content != content_type {
		return errors.New("Content-Type != " + content_type)
	}
	return nil
}

// trying unmarshal to book variable from request
func tryUnmarshalBook(r *http.Request) (*entity.Book, error) {
	if err := validateContent(r, json_content); err != nil {
		return nil, err
	}
	book := &entity.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetBookByIDHandler is http handler for GET /books/{id:[0-9]+}
func GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book_id, _ := strconv.ParseInt(vars["id"], 10, 32)
	book, _ := entity.GetBook(int32(book_id), connection.GetSession())
	if book == nil {
		textResponse(w, http.StatusNotFound, nil)
		return
	}

	jsonResponse(w, http.StatusFound, book)
}

// GetBooksByFilterHandler is http handler for GET /books with filtres
func GetBooksByFilterHandler(w http.ResponseWriter, r *http.Request) {
	// Getting and parsing filters
	filters := FilterMap{}
	var filter BookFilter
	filters["name"] = r.URL.Query().Get("name")
	filters["minPrice"] = r.URL.Query().Get("minPrice")
	filters["maxPrice"] = r.URL.Query().Get("maxPrice")
	filters["genre"] = r.URL.Query().Get("genre")

	err := filter.Parse(filters)
	if err != nil {
		textResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	books, err := GetBooks(&filter)
	if err != nil {
		textResponse(w, http.StatusInternalServerError, []byte("Server error"))
		return
	}

	jsonResponse(w, http.StatusFound, books)
}

// CreateBookHandler is http handler for POST /books/new
func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	book, err := tryUnmarshalBook(r)
	if err != nil {
		textResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	err = book.Insert(connection.GetSession())
	if err != nil {
		textResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	id := []byte(strconv.FormatInt(int64(book.Id), 10))
	textResponse(w, http.StatusCreated, id)
}

// UpdateBookHandler is http handler for PUT /books/{id:[0-9]+}
func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	book, err := tryUnmarshalBook(r)
	if err != nil {
		textResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	vars := mux.Vars(r)
	book_id, _ := strconv.ParseInt(vars["id"], 10, 32)
	book.Id = int32(book_id)

	err = book.Update(connection.GetSession())
	if err != nil {
		textResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	id := []byte(strconv.FormatInt(int64(book.Id), 10))
	textResponse(w, http.StatusOK, id)
}

// DeleteBookHandler is http handler for DELETE /books/{id:[0-9]+}
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book_id, _ := strconv.ParseInt(vars["id"], 10, 32)
	err := entity.DeleteBook(int32(book_id), connection.GetSession())
	if err != nil {
		textResponse(w, http.StatusNotFound, []byte(err.Error()))
	}
	w.WriteHeader(http.StatusNoContent)
}
