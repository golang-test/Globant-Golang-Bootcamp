package service

import (
	"database/sql"
	"testing"

	"github.com/DmitriyZhevnov/library/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) FindAll(db *sql.DB) (book []entities.Book, err error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Book), args.Error(1)
}

func TestFindAll(t *testing.T) {
	var floatValue float64 = 100
	mockRepository := new(MockRepository)
	book := entities.Book{Name: "book", Amount: 10, Price: floatValue, GenreId: 1}
	book2 := entities.Book{Name: "book2", Amount: 10, Price: floatValue, GenreId: 1}
	mockRepository.On("FindAll").Return([]entities.Book{book, book2}, nil)
	testService := NewBookService(mockRepository)
	result, _ := testService.FindAll(nil)
	mockRepository.AssertExpectations(t)
	assert.Equal(t, "book", result[0].Name)
	assert.Equal(t, 10, result[0].Amount)
	assert.Equal(t, floatValue, result[0].Price)
	assert.Equal(t, 1, result[0].GenreId)
	assert.Equal(t, "book2", result[1].Name)
}

func (mock *MockRepository) FindById(db *sql.DB, id int) (book []entities.Book, err error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entities.Book), args.Error(1)
}

func TestFindById(t *testing.T) {
	// var floatValue float64 = 100
	// mockRepository := new(MockRepository)
	// book2 := entities.Book{Id: 2, Name: "book2", Amount: 2, Price: floatValue, GenreId: 2}
	// mockRepository.On("FindById").Return([]entities.Book{book2}, nil)
	// testService := NewBookService(mockRepository)
	// result, err := testService.FindById(nil, "2")
	// mockRepository.AssertExpectations(t)
	// assert.Equal(t, 2, result[0].Id)
	// assert.Equal(t, "book2", result[0].Name)
	// assert.Equal(t, 2, result[0].Amount)
	// assert.Equal(t, floatValue, result[0].Price)
	// assert.Equal(t, 2, result[0].GenreId)
	// assert.Nil(t, err)
}

func (mock *MockRepository) FindByName(db *sql.DB, name string) (book []entities.Book, err error) {
	return nil, nil
}
func (mock *MockRepository) FilterByGenre(db *sql.DB, id int) (book []entities.Book, err error) {
	return nil, nil
}
func (mock *MockRepository) FilterByPrices(db *sql.DB, min, max float64) (book []entities.Book, err error) {
	return nil, nil
}
func (mock *MockRepository) Create(db *sql.DB, book *entities.Book) error {
	args := mock.Called()
	return args.Error(1)
}

func TestCreate(t *testing.T) {
	mockRepository := new(MockRepository)
	book := entities.Book{Name: "book", Amount: 10, Price: 100, GenreId: 1}
	mockRepository.On("Create").Return(&book, nil)
	testService := NewBookService(mockRepository)
	err := testService.Create(nil, &book)
	mockRepository.AssertExpectations(t)
	assert.Nil(t, err)
}

func (mock *MockRepository) Update(db *sql.DB, id int64, book *entities.Book) error {
	args := mock.Called()
	return args.Error(1)
}

func TestUpdate(t *testing.T) {
	// mockRepository := new(MockRepository)
	// book := entities.Book{Id: 1, Name: "book", Amount: 10, Price: 100, GenreId: 1}
	// newBook := entities.Book{Id: 1, Name: "book", Amount: 10, Price: 100, GenreId: 1}
	// mockRepository.On("Update").Return(nil)
	// testService := NewBookService(mockRepository)
	// err := testService.Update(nil, "1", &newBook)
	// mockRepository.AssertExpectations(t)
	// assert.Equal(t, book.Id, newBook.Id)
	// assert.Equal(t, book.Name, newBook.Name)
	// assert.Equal(t, book.Amount, newBook.Amount)
	// assert.Equal(t, book.Price, newBook.Price)
	// assert.Equal(t, book.GenreId, newBook.GenreId)
	// assert.Nil(t, err)
}

func (mock *MockRepository) Delete(db *sql.DB, id int64) (int64, error) {
	return 0, nil
}

func TestValidateBookWithEmptyName(t *testing.T) {
	testService := NewBookService(nil)
	book := entities.Book{Price: 1, GenreId: 1, Amount: 1}
	err := testService.Validate(&book)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid resquest payload", err.Error())
}

func TestValidateBookWithEmptyPrice(t *testing.T) {
	testService := NewBookService(nil)
	book := entities.Book{Name: "book", GenreId: 1, Amount: 1}
	err := testService.Validate(&book)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid resquest payload", err.Error())
}

func TestValidateBookWithEmptyGenreId(t *testing.T) {
	testService := NewBookService(nil)
	book := entities.Book{Name: "book", Price: 100, Amount: 1}
	err := testService.Validate(&book)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid resquest payload", err.Error())
}

func TestValidateBookWithEmptyAmount(t *testing.T) {
	testService := NewBookService(nil)
	book := entities.Book{Name: "book", Price: 100, GenreId: 1}
	err := testService.Validate(&book)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid resquest payload", err.Error())
}

func TestValidateEmptyBook(t *testing.T) {
	testService := NewBookService(nil)
	book := entities.Book{}
	err := testService.Validate(&book)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid resquest payload", err.Error())
}
