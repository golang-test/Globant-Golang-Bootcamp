package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFloat32(t *testing.T) {
	filters := make(FilterMap)
	filters["success"] = "123"
	filters["empty"] = ""
	filters["bad_value"] = "not float"

	buff, ok := parseFloat32(filters, "success")
	require.True(t, buff != nil && ok,
		"Wrong return of parsing. \nExpected: parsed value, true \n Actual: %v, %v", buff, ok)

	buff, ok = parseFloat32(filters, "empty")
	require.True(t, buff == nil && ok,
		"Wrong return of parsing. \nExpected: nil, true \n Actual: %v, %v", buff, ok)

	buff, ok = parseFloat32(filters, "bad_value")
	require.True(t, buff == nil && !ok,
		"Wrong return of parsing. \nExpected: nil, false \n Actual: %v, %v", buff, ok)
}

func TestPriceFilterParse(t *testing.T) {
	filters := make(FilterMap)
	filters["minPrice"] = "123"
	filters["maxPrice"] = "1000"

	f := PriceFilter{}
	err := f.Parse(filters)
	require.NoError(t, err, "Unexpected error on success testing")
	require.NotNil(t, f.minPrice, "Wrong minPrice. Expected value on returning")
	require.NotNil(t, f.maxPrice, "Wrong maxPrice. Expected value on returning")

	filters["minPrice"] = "bad"
	err = f.Parse(filters)
	require.Error(t, err, "Expected error on bad testing")
	require.Nil(t, f.minPrice, "Wrong minPrice. Expected value on returning")

}
