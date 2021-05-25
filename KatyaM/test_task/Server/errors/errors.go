package errors

import "errors"

var (
	ErrNotFound     = errors.New("error: not found")
	ErrSameName     = errors.New("error: create new book with name already exist")
	ErrWhileMarshal = errors.New("error: incorrect json marshal")
	ErrDBRequest    = errors.New("error: incorrect database request")
	ErrWhileUpdate  = errors.New("error: nothing to update")
)
