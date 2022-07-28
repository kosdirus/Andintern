package api

import (
	"encoding/json"
	"errors"
	"fmt"
	httpv1 "github.com/kosdirus/andintern/internal/api/http/handler/v1"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
	"net/http"
	"strconv"
)

const (
	queryParamID             = "id"
	queryParamBrand          = "brand"
	queryParamPriceLowerThan = "priceLowerThan"

	HeaderContentType   = "Content-Type"
	MimeApplicationJSON = "application/json"
)

var NoneQueryParamProvided = errors.New("none of required fields provided")

func ParseCarFilter(r *http.Request) (*dataprovider.CarFilter, error) {
	if idStr, ok := ParseQueryParam(r, queryParamID); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("wrong id '%s' query parameter: should be integer", idStr)
		}
		return dataprovider.NewCarFilter().ByID(id), nil
	}

	if brand, ok := ParseQueryParam(r, queryParamBrand); ok {
		return dataprovider.NewCarFilter().ByBrand(brand), nil
	}

	if priceLowerThanStr, ok := ParseQueryParam(r, queryParamPriceLowerThan); ok && r.Method == "DELETE" {
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

func ParseCreateCarRequest(r *http.Request) (*httpv1.CarToInsert, error) {
	carToCreate := httpv1.CarToInsert{}
	err := json.NewDecoder(r.Body).Decode(&carToCreate)
	if err != nil {
		return nil, err
	}

	return &carToCreate, nil
}

func ParseUpdateCarRequest(r *http.Request) (*httpv1.CarToUpdate, error) {
	idStr, ok := ParseQueryParam(r, queryParamID)
	if !ok {
		return nil, fmt.Errorf("missing id in URL query")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("wrong id '%s' query parameter: should be integer", idStr)
	}

	carToUpdate := httpv1.CarToUpdate{}
	err = json.NewDecoder(r.Body).Decode(&carToUpdate)
	if err != nil {
		return nil, err
	}
	carToUpdate.Id = id

	return &carToUpdate, nil
}
