package handler

import (
	_ "bytes"
	_ "database/sql"
	"encoding/json"
	_ "errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
	"io/ioutil"
	_ "io/ioutil"
	_ "log"
	"net/http"
	"strconv"
	"strings"
	_ "strings"
	b "test_task/Server/database_interface"
	"test_task/Server/errors"
	s "test_task/Server/service"
)

func NewBookHandler(l *zap.Logger, s s.BookService, db b.Database) *BookHandler {
	return &BookHandler{l: l, s: s, db: db}
}

type BookHandler struct {
	l  *zap.Logger
	s  s.BookService
	db b.Database
}

func (h *BookHandler) Register(mux *http.ServeMux) {
	mux.Handle("/", h)
}

func (h BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Info("start")
	if r.Method == http.MethodPost {
		if strings.Contains(r.RequestURI, "/create") {
			h.l.Info("create")
			jsonReq := b.Book{}
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			err = json.Unmarshal(body, &jsonReq)
			if err != nil {
				_, _ = w.Write([]byte(errors.ErrWhileMarshal.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp, err := h.db.CreateBook(r.Context(), jsonReq)
			if err != nil {
				_, _ = w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			jsonRes, err := json.Marshal(resp)
			if err != nil {
				_, _ = w.Write([]byte(errors.ErrWhileMarshal.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_, err = w.Write(jsonRes)
			if err != nil {
				_, _ = w.Write([]byte(errors.ErrWhileMarshal.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else if strings.Contains(r.RequestURI, "/update") {
			h.l.Info("update")
			marsh_id := r.URL.Query()["id"][0]
			id, _ := strconv.Atoi(marsh_id)
			h.l.Info(marsh_id)
			jsonReq := b.Book{}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				h.l.Error("error")
				_, _ = w.Write([]byte(errors.ErrDBRequest.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			h.l.Info(string(body))
			defer r.Body.Close()
			err = json.Unmarshal(body, &jsonReq)
			h.l.Error("err")
			h.l.Info(fmt.Sprint(jsonReq))
			err = h.db.UpdateBook(r.Context(), id, jsonReq)
			if err != nil {
				h.l.Error("err update book")
				_, _ = w.Write([]byte(errors.ErrWhileUpdate.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else if strings.Contains(r.RequestURI, "/get_book") {
			h.l.Info("get_book")
			marsh_id := r.URL.Query()["id"][0]
			id, _ := strconv.Atoi(marsh_id)
			jsonReq := b.Book{}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				h.l.Error("error")
				_, _ = w.Write([]byte(errors.ErrWhileMarshal.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			defer r.Body.Close()
			err = json.Unmarshal(body, &jsonReq)
			book, err := h.db.GetBook(r.Context(), id)
			if err != nil {
				h.l.Error("err get book")
				_, _ = w.Write([]byte(errors.ErrDBRequest.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			book_json, _ := json.Marshal(book)
			_, _ = w.Write(book_json)
			w.WriteHeader(http.StatusOK)
		} else if strings.Contains(r.RequestURI, "/delete") {
			h.l.Info("delete")
			marsh_id := r.URL.Query()["id"][0]
			h.l.Info(marsh_id)
			id, _ := strconv.Atoi(marsh_id)
			defer r.Body.Close()
			err := h.db.DeleteBook(r.Context(), id)
			if err != nil {
				h.l.Error("err delete book")
				_, _ = w.Write([]byte(errors.ErrDBRequest.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			h.l.Info("wrong request")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	if r.Method == http.MethodGet {
		if strings.Contains(r.RequestURI, "/get_book") {
			h.l.Info("get_book")
			marsh_id := r.URL.Query()["id"][0]
			h.l.Info(marsh_id)
			id, _ := strconv.Atoi(marsh_id)
			jsonReq := b.Book{}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				h.l.Error("error")
				_, _ = w.Write([]byte(errors.ErrWhileMarshal.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			defer r.Body.Close()
			err = json.Unmarshal(body, &jsonReq)
			book, err := h.db.GetBook(r.Context(), id)
			if err != nil {
				h.l.Error("err get book")
				_, _ = w.Write([]byte(errors.ErrDBRequest.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			book_json, _ := json.Marshal(book)
			_, _ = w.Write(book_json)
			w.WriteHeader(http.StatusOK)
		} else if strings.Contains(r.RequestURI, "/get_all") {
			h.l.Info("get_all")
			jsonReq := b.AllBooksRequest{}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				h.l.Error("error")
				_, _ = w.Write([]byte(errors.ErrWhileMarshal.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			defer r.Body.Close()
			err = json.Unmarshal(body, &jsonReq)
			books, err := h.db.GetAllBooks(r.Context(), jsonReq)
			if err != nil {
				h.l.Error("err get all books")
				_, _ = w.Write([]byte(errors.ErrDBRequest.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			books_json, _ := json.Marshal(books)
			_, _ = w.Write(books_json)
			w.WriteHeader(http.StatusOK)
		} else {
			h.l.Info("wrong request")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
