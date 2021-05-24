package handler

import (
	"fmt"
	"github.com/SavinskiPavel/bookstore/internal/model"
	"github.com/SavinskiPavel/bookstore/internal/repository"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

func HandleRequest() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainPage).Methods("GET")
	r.HandleFunc("/create", createPage).Methods("GET")
	r.HandleFunc("/save", save).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", infoPage).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/update", updatePage).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/update", update).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/delete", del).Methods("POST")
	r.HandleFunc("/search", searchPage).Methods("GET")
	r.HandleFunc("/find_by_id", findById).Methods("POST")
	r.HandleFunc("/find_by_name", findByName).Methods("POST")
	r.HandleFunc("/filter_by_price", filterByPrice).Methods("POST")
	r.HandleFunc("/filter_by_genre", filterByGenre).Methods("POST")
	http.Handle("/", r)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	books := repository.FindAllBooks()
	t, err := template.ParseFiles("web/template/mainPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	_ = t.Execute(w, books)
}

func createPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/createPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	data := model.Select{
		ListGenres: repository.FindAllGenresByValue(),
	}
	_ = t.Execute(w, data)
}

func save(w http.ResponseWriter, r *http.Request) {
	books := repository.FindAllBooks()

	if r.FormValue("name") == "" || r.FormValue("price") == "" || r.FormValue("genre") == "" || r.FormValue("amount") == "" {
		http.Error(w, "all fields must be filled", 400)
		return
	}
	name := r.FormValue("name")
	price, e := strconv.ParseFloat(r.FormValue("price"), 64)
	if e != nil {
		http.Error(w, "price: incorrect data", 400)
		return
	}
	genre, e := strconv.Atoi(r.FormValue("genre"))
	if e != nil {
		http.Error(w, "genre: incorrect data", 400)
		return
	}
	amount, e := strconv.Atoi(r.FormValue("amount"))
	if e != nil {
		http.Error(w, "amount: incorrect data", 400)
		return
	}
	b := model.Book{Name: name, Price: price, Genre: genre, Amount: amount}
	for _, book := range books {
		if b.Name == book.Name {
			http.Error(w, "a book with this name already exists", 400)
			return
		}
	}

	err := b.Validate()
	if err != nil {
		http.Error(w, "the entered data must be >=0", 400)
	} else {
		roundPrice := math.Round(price*100) / 100
		_, err = repository.SaveBook(name, roundPrice, genre, amount)
		if err != nil {
			log.Println(w, err)
		}
	}
	http.Redirect(w, r, "/", 302)
}

func infoPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/infoPage.html")
	if err != nil {
		_, _ = fmt.Fprintln(w, err)
	}
	vars := mux.Vars(r)
	book := repository.FindBookById(vars["id"])
	if book.Name == "" {
		http.Error(w, "404 page not found", 404)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = t.Execute(w, book)
}

func updatePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/updatePage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	vars := mux.Vars(r)
	data := model.Select{
		ListGenres: repository.FindAllGenresByValue(),
		Book:       repository.FindBookById(vars["id"]),
	}
	_ = t.Execute(w, data)
}

func update(w http.ResponseWriter, r *http.Request) {
	books := repository.FindAllBooks()
	vars := mux.Vars(r)
	if r.FormValue("name") == "" || r.FormValue("price") == "" || r.FormValue("genre") == "" || r.FormValue("amount") == "" {
		http.Error(w, "all fields must be filled", 400)
		return
	}
	name := r.FormValue("name")
	price, e := strconv.ParseFloat(r.FormValue("price"), 64)
	if e != nil {
		http.Error(w, "price: incorrect data", 400)
		return
	}
	genre, e := strconv.Atoi(r.FormValue("genre"))
	if e != nil {
		http.Error(w, "genre: incorrect data", 400)
		return
	}
	amount, e := strconv.Atoi(r.FormValue("amount"))
	if e != nil {
		http.Error(w, "amount: incorrect data", 400)
		return
	}
	b := model.Book{Name: name, Price: price, Genre: genre, Amount: amount}
	for _, book := range books {
		if b.Name == book.Name {
			http.Error(w, "a book with this name already exists", 400)
			return
		}
	}
	err := b.Validate()
	if err != nil {
		http.Error(w, "the entered data must be >=0", 400)
	} else {
		roundPrice := math.Round(price*100) / 100
		repository.UpdateBook(vars["id"], name, roundPrice, genre, amount)
	}
	http.Redirect(w, r, "/", 302)
}

func del(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repository.DeleteBook(vars["id"])
	w.WriteHeader(http.StatusNoContent)
}

func searchPage(w http.ResponseWriter, r *http.Request) {

	var books []model.Book
	b := repository.FindAllBooks()
	t, err := template.ParseFiles("web/template/searchPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	for _, i := range b {
		if i.Amount == 0 {
			continue
		}
		books = append(books, i)
	}
	data := model.Select{
		ListGenres: repository.FindAllGenresByValue(),
		Books:      books,
	}
	_ = t.Execute(w, data)
}

func findById(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/filterPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	books := repository.FindBooksById(id)
	_ = t.Execute(w, books)
}

func findByName(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/filterPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	name := r.FormValue("name")
	books := repository.FindBooksByName(name)
	_ = t.Execute(w, books)
}

func filterByPrice(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("web/template/filterPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	minString := r.FormValue("min")
	maxString := r.FormValue("max")
	if minString != "" && maxString != "" {
		min, e := strconv.ParseFloat(r.FormValue("min"), 64)
		if e != nil {
			http.Error(w, "price: incorrect data", 400)
			return
		}
		max, e := strconv.ParseFloat(r.FormValue("max"), 64)
		if e != nil {
			http.Error(w, "price: incorrect data", 400)
			return
		}
		books := repository.FindBooksByPrices(min, max)
		_ = t.Execute(w, books)
	} else if minString == "" && maxString != "" {
		max, e := strconv.ParseFloat(r.FormValue("max"), 64)
		if e != nil {
			http.Error(w, "price: incorrect data", 400)
			return
		}
		books := repository.FindBooksByMaxPrice(max)
		_ = t.Execute(w, books)
	} else if minString != "" && maxString == "" {
		min, e := strconv.ParseFloat(r.FormValue("min"), 64)
		if e != nil {
			http.Error(w, "price: incorrect data", 400)
			return
		}
		books := repository.FindBooksByMinPrice(min)
		_ = t.Execute(w, books)
	} else if minString == "" && maxString == "" {

		var books []model.Book
		_ = t.Execute(w, books)
	}
}

func filterByGenre(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/filterPage.html")
	if err != nil {
		http.Error(w, "404 page not found", 404)
	}
	genre, e := strconv.Atoi(r.FormValue("genre"))
	if e != nil {
		http.Error(w, "genre: incorrect data", 400)
		return
	}
	books := repository.FindBooksByGenre(genre)
	_ = t.Execute(w, books)
}
