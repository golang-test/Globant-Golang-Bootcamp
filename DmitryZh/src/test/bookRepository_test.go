package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/DmitriyZhevnov/library/src/app"
)

var a app.App

func TestMain(m *testing.M) {
	// a.Initialize("postgres", "root", "root", "5432", "postgres_test", "fullstack_api_test")
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "postgres_test", "5432", "root", "fullstack_api_test", "root")
	a.DB, err = sql.Open("postgres", DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", "postgres")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", "postgres")
	}
	os.Exit(m.Run())
}

func TestFindAll(t *testing.T) {
	// _, err := repository.FindAll(a.DB)
	// if err != nil {
	// 	t.Errorf("OK")
	// }

}

func TestFindById(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFindNonExistentId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books/1000", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Book not found'. Got '%s'", m["error"])
	}
}

func TestFindInvalidId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books/1fff", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid book Id" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid book Id'. Got '%s'", m["error"])
	}
}

func TestFindByName(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/name/book1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFindByGenre(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/genre/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFilterByPrices(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/price/10/20", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFilterByOnlyOnePrice(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/price/10", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdate(t *testing.T) {
	clearTable()
	addBook()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	payload := []byte(`{
		"name": "Some new name",
		"price": 15,
		"genre": 2,
		"amount": 10 }`)
	req, err := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] == originalBook["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalBook["name"], m["name"], m["name"])
	}

	if m["genre"] == originalBook["genre"] {
		t.Errorf("Expected the genre to change from '%v' to '%v'. Got '%v'", originalBook["genre"], m["genre"], m["genre"])
	}
	if m["amount"] == originalBook["amount"] {
		t.Errorf("Expected the amount to change from '%v' to '%v'. Got '%v'", originalBook["amount"], m["amount"], m["amount"])
	}
}

func TestUpdateWithInvalidData(t *testing.T) {
	clearTable()
	addBook()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	payload := []byte(`{
		"name": "Some new name",
		"price": -1,
		"genre": 2,
		"amount": -10 }`)
	req, err := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestUpdateWithMissedField(t *testing.T) {
	clearTable()
	addBook()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	payload := []byte(`{
		"name": "Some new name",
		"price": 10,
		"genre": 2 }`)
	req, err := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateWithInvalidData(t *testing.T) {
	clearTable()
	addBook()
	payload := []byte(`{
		"name": "Some new name",
		"price": -1,
		"genre": 2,
		"amount": -10 }`)
	req, err := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateWithMissedField(t *testing.T) {
	clearTable()
	addBook()
	payload := []byte(`{
		"name": "Some new name",
		"genre": 2,
		"amount": 10 }`)
	req, err := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreate(t *testing.T) {
	jsonStr := []byte(`{
			"name": "The Three Musketeers",
			"price": 10.44,
			"genre": 1,
			"amount": 5 }`)
	req, _ := http.NewRequest("POST", "/api/books", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	var i = response.Body
	_, err := strconv.ParseInt(i.String(), 10, 64)
	if err != nil {
		t.Errorf("Erorr! %v", err)
	}
}

func TestDelete(t *testing.T) {
	clearTable()
	addBook()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/books/1", nil)
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/api/books/1", nil)
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addBook() {
	a.DB.Exec(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('test_book', '25', '1', '30');`)
}

func clearTable() {
	a.DB.Exec("DELETE FROM library.book")
	a.DB.Exec("ALTER TABLE library.book AUTO_INCREMENT = 1")
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
