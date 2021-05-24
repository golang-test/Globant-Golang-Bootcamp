package model

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type Book struct {
	Id     int     `json: "id"`
	Name   string  `json: "name"`
	Price  float64 `json: "price"`
	Genre  int     `json: "genre" `
	Amount int     `json: "amount"`
}

func (b Book) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Length(1, 100)),
		validation.Field(&b.Price, validation.Min(0.0)),
		validation.Field(&b.Genre, validation.Min(0)),
		validation.Field(&b.Amount, validation.Min(0)),
	)

}
