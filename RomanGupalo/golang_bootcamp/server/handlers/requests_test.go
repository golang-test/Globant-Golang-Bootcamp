package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/AA55hex/golang_bootcamp/server/config"
	"github.com/AA55hex/golang_bootcamp/server/db/connection"
	"github.com/AA55hex/golang_bootcamp/server/db/entity"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

type request_tester struct {
	Method  string
	Request string
	Body    io.Reader
}

var books db.Collection
var router *mux.Router

// init session && create router
func init() {
	config.LoadConfigs("../configs.env")
	db_settings := mysql.ConnectionURL{
		Database: config.MySQL.Database,
		Host:     config.MySQL.Host,
		User:     config.MySQL.User,
		Password: config.MySQL.Password,
	}
	session, err := connection.OpenSession(&db_settings, 50)
	if err != nil {
		log.Fatal(err)
	}

	books = session.Collection("book")

	router = mux.NewRouter()
	router.HandleFunc("/books/{id:[0-9]+}", GetBookByIDHandler).Methods("GET")
	router.HandleFunc("/books/{id:[0-9]+}", UpdateBookHandler).Methods("PUT")
	router.HandleFunc("/books/{id:[0-9]+}", DeleteBookHandler).Methods("DELETE")
	router.HandleFunc("/books", GetBooksByFilterHandler).Methods("GET")
	router.HandleFunc("/books/new", CreateBookHandler).Methods("POST")
}

// simple way for getting test book object
func getTestBook() *entity.Book {
	name := "ThisBookForTestsAndOnlyForTestsNotForYou!!!!!!!"
	price := float32(10)
	genre := int(1)
	amount := int(100)
	book := entity.Book{
		Name:   &name,
		Price:  &price,
		Genre:  &genre,
		Amount: &amount,
	}
	return &book
}

// simple way for inserting & getting test book
func createTestBook() (*entity.Book, error) {
	book := getTestBook()
	err := books.InsertReturning(book)
	return book, err
}

// clean database from testing objects
func clearDatabase() {
	books := connection.GetSession().Collection("book")
	res := books.Find(db.Cond{"name": *getTestBook().Name})
	res.Delete()
}

// returns served response
func serveHTTP(method string, url string, body io.Reader, content_type string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, body)
	if content_type != "" {
		req.Header.Add("Content-Type", content_type)
	}
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr, nil
}

func TestGetBookByIdRequestOnSuccess(t *testing.T) {
	clearDatabase()
	book, err := createTestBook()
	require.NoError(t, err, err)

	str_id := strconv.FormatInt(int64(book.Id), 10)

	rr, err := serveHTTP("GET", "/books/"+str_id, nil, no_content)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusFound,
		"handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)

	expected := book
	actual := entity.Book{}
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	require.NoError(t, err, "handler returned unexpected body: \ngot: %v",
		rr.Body.String())
	require.True(t, entity.BookEqual(expected, &actual),
		"handler returned unexpected body: \ngot: %v\n want: %v",
		actual, expected)
}

func TestGetBookByIdRequestOnNotFound(t *testing.T) {
	clearDatabase()
	book, err := createTestBook()
	require.NoError(t, err, err)
	str_id := strconv.FormatInt(int64(book.Id), 10)
	book.Delete(connection.GetSession())

	rr, err := serveHTTP("GET", "/books/"+str_id, nil, no_content)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusNotFound,
		"handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
}

func TestCreateBookHandlerOnSuccess(t *testing.T) {
	clearDatabase()
	const testBookJSON = `{ "name": "ThisBookForTestsAndOnlyForTestsNotForYou!!!!!!!", "price": 9999, "genre": 1, "amount": 9999 }`
	reader := strings.NewReader(testBookJSON)
	rr, err := serveHTTP("POST", "/books/new", reader, json_content)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusCreated,
		"handler returned wrong status code: got %v want %v\nbody:%v",
		rr.Code, http.StatusOK, rr.Body.String())

	id, err := strconv.ParseInt(rr.Body.String(), 10, 32)
	require.NoError(t, err, err)

	book, err := entity.GetBook(int32(id), connection.GetSession())
	require.NoError(t, err, err)

	expected := "ThisBookForTestsAndOnlyForTestsNotForYou!!!!!!!"
	actual := *book.Name
	require.Equal(t, expected, actual,
		"handler returned unexpected body: \ngot: %v\n want: %v",
		actual, expected)
}

func TestCreateBookHandlerOnBadRequest(t *testing.T) {
	clearDatabase()
	const jsonBadBody = `{ "price": 9999, "genre": 1, "amount": 9999 }`
	reader := strings.NewReader(jsonBadBody)
	rr, err := serveHTTP("POST", "/books/new", reader, json_content)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusBadRequest,
		"handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
}

func getBookReader(book *entity.Book) io.Reader {
	book_json, _ := json.Marshal(book)
	return bytes.NewReader(book_json)
}
func TestUpdateBookHandlerOnSuccess(t *testing.T) {
	clearDatabase()
	book, err := createTestBook()
	require.NoError(t, err, err)
	str_id := strconv.FormatInt(int64(book.Id), 10)

	good_id := book.Id
	book.Id = book.Id + 1
	*book.Price = 0

	rr, err := serveHTTP("PUT", "/books/"+str_id, getBookReader(book), json_content)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusOK,
		"handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	require.Equal(t, rr.Body.String(), str_id,
		"handler returned wrong id: got %v want %v", rr.Body.String(), str_id)

	// Check row from db
	db_book, err := entity.GetBook(good_id, connection.GetSession())
	require.NoError(t, err, err)

	expected := *book.Price
	actual := *db_book.Price
	require.Equal(t, expected, actual,
		"database returned bad record: \ngot: %v\n want: %v",
		actual, expected)
}

func TestUpdateBookHandlerOnBadRequest(t *testing.T) {
	clearDatabase()
	book, err := createTestBook()
	require.NoError(t, err, err)
	str_id := strconv.FormatInt(int64(book.Id), 10)

	// test function to avoid code duplication
	test_func := func() {
		rr, err := serveHTTP("PUT", "/books/"+str_id, getBookReader(book), json_content)
		require.NoError(t, err, err)

		require.Equal(t, rr.Code, http.StatusBadRequest,
			"handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	book.Price = nil
	test_func()

	book.Price = new(float32)
	*book.Price = -1
	test_func()

	*book.Price = 0
	*book.Amount = -1
	test_func()
}

func TestDeleteBookHandlerOnSuccess(t *testing.T) {
	clearDatabase()
	book, err := createTestBook()
	require.NoError(t, err, err)
	str_id := strconv.FormatInt(int64(book.Id), 10)

	rr, err := serveHTTP("DELETE", "/books/"+str_id, nil, text_plain)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusNoContent,
		"handler returned wrong status code: got %v want %v", rr.Code, http.StatusNoContent)
}

func TestDeleteBookHandlerOnNotFound(t *testing.T) {
	clearDatabase()
	book, err := createTestBook()
	require.NoError(t, err, err)
	str_id := strconv.FormatInt(int64(book.Id), 10)
	book.Delete(connection.GetSession())

	rr, err := serveHTTP("DELETE", "/books/"+str_id, nil, text_plain)
	require.NoError(t, err, err)

	require.Equal(t, rr.Code, http.StatusNotFound,
		"handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
}

func makeBook(name string, price float32, genre int, amount int) *entity.Book {
	book := entity.Book{
		Name:   &name,
		Price:  &price,
		Genre:  &genre,
		Amount: &amount,
	}
	return &book
}
func TestGetBooksByFilterHandler(t *testing.T) {
	clearDatabase()
	// creating test data set of 100 objects
	const capacity = 100
	test_books := make([]*entity.Book, 0, capacity)
	for i := 0; i < capacity; i++ {
		name := "only_test_book_" + strconv.FormatInt(int64(i), 10)
		price := rand.Float32() * 100
		genre := rand.Uint32()%3 + 1
		amount := rand.Uint32() % 2
		buff := makeBook(name, price, int(genre), int(amount))
		test_books = append(test_books, buff)
		(*buff).Insert(connection.GetSession())
	}

	// clearing after testing
	defer func() {
		for i := range test_books {
			(test_books[i]).Delete(connection.GetSession())
		}
	}()

	// create map that consist of the book names & expected status
	// expected status = true -> book must be in returning test set
	get_expected := func(name string, lp *float32, rp *float32, genre *int) (result map[string]bool) {
		result = map[string]bool{}
		name_test := name == ""
		lp_test := lp == nil
		rp_test := rp == nil
		genre_test := genre == nil
		for _, item := range test_books {
			is_expected := *item.Amount != 0 &&
				(genre_test || *item.Genre == *genre) &&
				(name_test || *item.Name == name) &&
				(lp_test || *item.Price >= *lp) &&
				(rp_test || *item.Price <= *rp)

			result[*item.Name] = is_expected
		}
		return
	}

	// function for testing filtering
	test_func := func(name string, lp *float32, rp *float32, genre *int) bool {
		// creating request
		req, err := http.NewRequest("GET", "/books", nil)
		values := req.URL.Query()
		values.Set("name", name)
		if lp != nil {
			lp_str := strconv.FormatFloat(float64(*lp), 'f', 6, 32)
			values.Add("minPrice", lp_str)
		}
		if rp != nil {
			rp_str := strconv.FormatFloat(float64(*rp), 'f', 6, 32)
			values.Add("maxPrice", rp_str)
		}
		if genre != nil {
			genre_str := strconv.FormatInt(int64(*genre), 10)
			values.Add("genre", genre_str)
		}
		req.URL.RawQuery = values.Encode()
		// serving
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		require.NoError(t, err, err)

		expected := get_expected(name, lp, rp, genre)

		// creating actual slice
		actual := []entity.Book{}
		err = json.NewDecoder(rr.Body).Decode(&actual)

		// creating actual map for more simpler validation
		actual_map := map[string]bool{}
		for _, item := range actual {
			actual_map[*item.Name] = true
		}

		// validate
		for str, exp_status := range expected {
			_, act_ok := actual_map[str]
			if act_ok != exp_status {
				return false
			}
		}
		return true
	}

	// testing
	var status bool
	params := struct {
		name  string
		lp    *float32
		rp    *float32
		genre *int
	}{}

	status = test_func(params.name, params.lp, params.rp, params.genre)
	require.True(t, status, "Bad response with params: %v", params)

	params.lp = new(float32)
	*params.lp = 10
	status = test_func(params.name, params.lp, params.rp, params.genre)
	require.True(t, status, "Bad response with params: %v", params)

	params.rp = new(float32)
	*params.rp = 20
	status = test_func(params.name, params.lp, params.rp, params.genre)
	require.True(t, status, "Bad response with params: %v", params)

	params.name = "only_test_book_1"
	params.lp = nil
	params.genre = new(int)
	*params.genre = 2
	status = test_func(params.name, params.lp, params.rp, params.genre)
	require.True(t, status, "Bad response with params: %v", params)
}
