package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DmitriyZhevnov/library/src/entities"
	"github.com/DmitriyZhevnov/library/src/repository"
	service "github.com/DmitriyZhevnov/library/src/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

var (
	bookRepository repository.BookRepository = repository.NewPostgresBookRepository()
	bookService    service.BookService       = service.NewBookService(bookRepository)
)

func (a *App) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable ",
		DbHost, DbPort, DbUser, DbName, DbPassword)
	var err error
	a.DB, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/books", a.FindAll).Methods("GET")
	a.Router.HandleFunc("/api/books/{id}", a.FindById).Methods("GET")
	a.Router.HandleFunc("/api/books/name/{name}", a.FindByName).Methods("GET")
	a.Router.HandleFunc("/api/books/genre/{id:[0-9]+}", a.FilterByGenre).Methods("GET")
	a.Router.HandleFunc("/api/books/price/{minPrice:[0-9]+}/{maxPrice:[0-9]+}", a.FilterByPrices).Methods("GET")
	a.Router.HandleFunc("/api/books/price/{minPrice:[0-9]+}", a.FilterByPricesWithOneParameter).Methods("GET")
	a.Router.HandleFunc("/api/books", a.Create).Methods("POST")
	a.Router.HandleFunc("/api/books/{id:[0-9]+}", a.Update).Methods("PUT")
	a.Router.HandleFunc("/api/books/{id:[0-9]+}", a.Delete).Methods("DELETE")
}

func (a *App) FindAll(w http.ResponseWriter, r *http.Request) {
	books, err := bookService.FindAll(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, books)
}

func (a *App) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book := []entities.Book{}
	var err error
	if book, err = bookService.FindById(a.DB, vars["id"]); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) FindByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book := []entities.Book{}
	var err error
	if book, err = bookService.FindByName(a.DB, vars["name"]); err != nil {
		respondWithError(w, http.StatusNotFound, "Book not found")
		return
	}
	if book == nil {
		var array [5]string
		respondWithJSON(w, http.StatusOK, array)
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) FilterByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book := []entities.Book{}
	var err error
	if book, err = bookService.FilterByGenre(a.DB, vars["id"]); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if book == nil {
		var array [5]string
		respondWithJSON(w, http.StatusOK, array)
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) FilterByPricesWithOneParameter(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusBadRequest, "Please, enter the second parameter")
}

func (a *App) FilterByPrices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book := []entities.Book{}
	var err error
	if book, err = bookService.FilterByPrices(a.DB, vars["minPrice"], vars["maxPrice"]); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if book == nil {
		var array [5]string
		respondWithJSON(w, http.StatusOK, array)
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	var book entities.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil || bookService.Validate(&book) != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	err = bookService.Create(a.DB, &book)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, book.Id)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var book entities.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil || bookService.Validate(&book) != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	if err := bookService.Update(a.DB, vars["id"], &book); err != nil {
		respondWithError(w, http.StatusNotModified, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, err := bookService.Delete(a.DB, vars["id"]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusNoContent, map[string]string{"result": "Successfully deleted"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
