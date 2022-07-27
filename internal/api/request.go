package api

import (
	"errors"
	"fmt"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
	"net/http"
	"strconv"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJSON = "application/json"
)

var NoneQueryParamProvided = errors.New("none of required fields provided")

func ParseCarFilter(r *http.Request) (*dataprovider.CarFilter, error) {
	if idStr, ok := ParseQueryParam(r, "id"); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("wrong id '%s' query parameter: should be integer", idStr)
		}
		return dataprovider.NewCarFilter().ByID(id), nil
	}

	if brand, ok := ParseQueryParam(r, "brand"); ok {
		return dataprovider.NewCarFilter().ByBrand(brand), nil
	}

	if priceLowerThanStr, ok := ParseQueryParam(r, "priceLowerThan"); ok && r.Method == "DELETE" {
		priceLowerThan, err := strconv.Atoi(priceLowerThanStr)
		if err != nil {
			return nil, fmt.Errorf("wrong priceLowerThan '%s' query parameter: should be integer", priceLowerThanStr)
		}
		return dataprovider.NewCarFilter().ByPrice(priceLowerThan), nil
	}

	if r.Method == "GET" {
		return dataprovider.NewCarFilter(), nil
	}

	return nil, NoneQueryParamProvided
}

func ParseQueryParam(r *http.Request, field string) (string, bool) {
	q := r.URL.Query()

	param := q.Get(field)
	if param == "" {
		//return "", fmt.Errorf("empty param with field %q", field)
		return "", false
	}

	return param, true
}
