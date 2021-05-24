package handler

import (
	"fmt"
	"github.com/SavinskiPavel/bookstore/internal/database"
	_ "github.com/SavinskiPavel/bookstore/internal/database"
	"github.com/SavinskiPavel/bookstore/internal/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "Test book")
	form.Add("price", "10.56")
	form.Add("amount", "5")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	newBook := repository.FindBookById("4")
	assert.Equal(t, 302, w.Result().StatusCode)
	assert.Equal(t, "", w.Body.String())
	assert.NotEqual(t, beforeSave, afterSave)
	assert.Equal(t, "Test book", newBook.Name)
	assert.Equal(t, 10.56, newBook.Price)
	assert.Equal(t, 5, newBook.Amount)
	assert.Equal(t, 1, newBook.Genre)
}

func TestEmptyFieldSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "")
	form.Add("price", "")
	form.Add("amount", "5")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	assert.Equal(t, 400, w.Result().StatusCode)
	assert.Equal(t, "all fields must be filled\n", w.Body.String())
	assert.Equal(t, beforeSave, afterSave)
}

func TestIncorrectPriceSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "test book")
	form.Add("price", "10d58")
	form.Add("amount", "5")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	assert.Equal(t, 400, w.Result().StatusCode)
	assert.Equal(t, "price: incorrect data\n", w.Body.String())
	assert.Equal(t, beforeSave, afterSave)
}

func TestIncorrectAmountSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "test book")
	form.Add("price", "10")
	form.Add("amount", "5fff")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	assert.Equal(t, 400, w.Result().StatusCode)
	assert.Equal(t, "amount: incorrect data\n", w.Body.String())
	assert.Equal(t, beforeSave, afterSave)
}
func TestInvalidPriceSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "Book")
	form.Add("price", "-56")
	form.Add("amount", "2")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	assert.Equal(t, 400, w.Result().StatusCode)
	assert.Equal(t, "the entered data must be >=0\n", w.Body.String())
	assert.Equal(t, beforeSave, afterSave)
}

func TestInvalidAmountSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "Book")
	form.Add("price", "56")
	form.Add("amount", "-2")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	assert.Equal(t, 400, w.Result().StatusCode)
	assert.Equal(t, "the entered data must be >=0\n", w.Body.String())
	assert.Equal(t, beforeSave, afterSave)
}

func TestBookExistsSave(t *testing.T) {
	database.InitDB()
	booksBefore := repository.FindAllBooks()
	beforeSave := len(booksBefore)
	form := url.Values{}
	form.Add("name", "Test book")
	form.Add("price", "10")
	form.Add("amount", "5")
	form.Add("genre", "1")
	r, err := http.NewRequest("POST", "/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	save(w, r)
	booksAfter := repository.FindAllBooks()
	afterSave := len(booksAfter)
	assert.Equal(t, 400, w.Result().StatusCode)
	assert.Equal(t, "a book with this name already exists\n", w.Body.String())
	assert.Equal(t, beforeSave, afterSave)
}

func TestUpdate(t *testing.T) {
	database.InitDB()
	form := url.Values{}
	form.Add("name", "Updated test book")
	form.Add("price", "7.31")
	form.Add("amount", "11")
	form.Add("genre", "3")
	r, err := http.NewRequest("POST", "/4/update", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	r.Form = form
	w := httptest.NewRecorder()
	update(w, r)
	assert.Equal(t, 302, w.Result().StatusCode)
	assert.Equal(t, "", w.Body.String())

}

func TestDelete(t *testing.T) {
	database.InitDB()
	r, err := http.NewRequest("DELETE", fmt.Sprintf("/4/delete"), nil)
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	w := httptest.NewRecorder()
	del(w, r)
	assert.Equal(t, 204, w.Result().StatusCode)
	assert.Equal(t, "", w.Body.String())
}
