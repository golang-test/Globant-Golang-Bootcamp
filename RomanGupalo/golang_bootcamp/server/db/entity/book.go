package entity

import (
	"errors"

	"github.com/upper/db/v4"
)

// Book is main structure for interaction with database
type Book struct {
	Id     int32    `json:"id" db:"id,omitempty"`
	Name   *string  `json:"name" db:"name"`
	Price  *float32 `json:"price" db:"price"`
	Genre  *int     `json:"genre" db:"genre"`
	Amount *int     `json:"amount" db:"amount"`
}

// Validate full validate book for insertion to db
// with db requests
// Returns nil on success
func (b *Book) Validate(session db.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	// is not unique
	books := session.Collection("book")
	genres := session.Collection("genre")

	// check for name duplications
	name_duplications, _ := books.Find(db.Cond{"name": b.Name}).Count()
	if name_duplications != 0 {
		return errors.New("name is not unique")
	}

	// check for genre existence
	genre_existence, _ := genres.Find(db.Cond{"id": b.Genre}).Count()
	if genre_existence != 0 {
		return errors.New("bad genre id")
	}

	// simple validations
	err := b.SimpleValidate()
	return err
}

// SimpleValidate validation without db requests
// Returns nil on success
func (b *Book) SimpleValidate() error {
	switch {
	case b.Price == nil, *b.Price < 0:
		return errors.New("bad price")
	case b.Amount == nil, *b.Amount < 0:
		return errors.New("bad amount")
	case b.Name == nil:
		return errors.New("bad name")
	case b.Genre == nil:
		return errors.New("bad genre")
	default:
		return nil
	}
}

// Insert perform simple validation and trying to insert book object in database
// If the operation succeeds, updates current
// object with data from the newly inserted row
// Returns nil on success
func (b *Book) Insert(session db.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	if err := b.SimpleValidate(); err != nil {
		return errors.New("insert validation failed: " + err.Error())
	}

	books := session.Collection("book")
	err := books.InsertReturning(b)
	if err != nil {
		return errors.New("insertion failed: " + err.Error())
	}

	return nil
}

// Update perform simple validation and trying
// to update book object in database
// Returns nil on success
func (b *Book) Update(session db.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	books := session.Collection("book")

	res := books.Find(db.Cond{"id": b.Id})
	defer res.Close()

	// check for existence
	if count, _ := res.Count(); count != 1 {
		return errors.New("update validation failed: book not found")
	}

	// simple validations
	err := b.SimpleValidate()
	if err != nil {
		return errors.New("update validation failed: " + err.Error())
	}

	// try to update
	err = res.Update(b)
	if err != nil {
		return errors.New("update failed: " + err.Error())
	}
	return nil
}

// Delete validate and trying to delete book object from database
// Returns nil on success
func (b *Book) Delete(session db.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	books := session.Collection("book")

	res := books.Find(db.Cond{"id": b.Id})
	defer res.Close()

	// check for existence
	if count, _ := res.Count(); count != 1 {
		return errors.New("validation failed: book not found")
	}

	// try to delete
	if err := res.Delete(); err != nil {
		return errors.New("deleting failed: book not found")
	}

	return nil
}

// GetBook finds book by id
// Returns book, nil on success
func GetBook(book_id int32, session db.Session) (*Book, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	books := session.Collection("book")
	result := &Book{}

	err := books.Find(db.Cond{"id": book_id}).One(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteBook delete book by id
// Returns nil on success
func DeleteBook(book_id int32, session db.Session) error {
	if session == nil {
		return errors.New("session is nil")
	}

	book := Book{Id: int32(book_id)}
	err := book.Delete(session)
	return err
}

// BookEqual equals books
// Returns true on success
func BookEqual(l *Book, r *Book) bool {
	return l.Id == r.Id &&
		*l.Name == *r.Name &&
		*l.Price == *r.Price &&
		*l.Genre == *r.Genre &&
		*l.Amount == *r.Amount
}
