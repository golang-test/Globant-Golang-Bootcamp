package handlers

import (
	"errors"
	"strconv"

	"github.com/AA55hex/golang_bootcamp/server/db/connection"
	"github.com/AA55hex/golang_bootcamp/server/db/entity"
)

// FilterMap is map type for filter values
// using for BookFilter parsing
type FilterMap map[string]string

// PriceFilter is filter with min & max values of price
type PriceFilter struct {
	minPrice *float32
	maxPrice *float32
}

// BookFilter is filter for GetBooks func
type BookFilter struct {
	Name  string
	Price PriceFilter
	Genre *int32
}

// Parse into structure filter parameters
func (f *BookFilter) Parse(filters FilterMap) error {

	f.Name = filters["name"]

	err := f.Price.Parse(filters)
	if err != nil {
		return err
	}

	if filters["genre"] != "" {
		genre64, err := strconv.ParseInt(filters["genre"], 10, 32)
		if err != nil {
			return errors.New("bad genre")
		}
		genre32 := int32(genre64)
		f.Genre = &genre32
	}
	return nil
}

// (oject), nil - parsed
// nil, nil - good return without data (map string was empty)
// nil, err - parsing failed
func parseFloat32(filters FilterMap, key string) (*float32, bool) {
	if filters[key] != "" {
		buff64, err := strconv.ParseFloat(filters[key], 32)

		if err != nil {
			return nil, false
		}
		buff32 := float32(buff64)
		return &buff32, true
	}
	return nil, true
}

// Parse into structure filter parameters
// Returns nil on success
func (p *PriceFilter) Parse(filters FilterMap) error {
	var ok bool
	p.minPrice, ok = parseFloat32(filters, "minPrice")
	if !ok {
		return errors.New("bad minPrice")
	}

	p.maxPrice, ok = parseFloat32(filters, "maxPrice")
	if !ok {
		return errors.New("bad maxPrice")
	}

	return nil
}

// GetBooks create sql-query for db and returns result on success
func GetBooks(filter *BookFilter) ([]entity.Book, error) {
	query := connection.GetSession().SQL().SelectFrom("book")
	// create variable to determine next query function

	// check for 0 amoount
	query = query.Where("amount != 0")

	// name filtering
	if filter.Name != "" {
		query = query.And("name = ?", filter.Name)
	}

	// genre filtering
	if filter.Genre != nil {
		query = query.And("genre = ?", *filter.Genre)
	}

	// prcie filtering
	switch {
	case filter.Price.minPrice != nil && filter.Price.maxPrice != nil:
		query = query.And("price between ? and ?",
			*filter.Price.minPrice,
			*filter.Price.maxPrice)
	case filter.Price.minPrice != nil && filter.Price.maxPrice == nil:
		query = query.And("price >= ?", *filter.Price.minPrice)
	case filter.Price.minPrice == nil && filter.Price.maxPrice != nil:
		query = query.And("price <= ?", *filter.Price.maxPrice)
	}

	var result []entity.Book
	err := query.All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
