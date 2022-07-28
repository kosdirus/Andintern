package v1

import (
	"encoding/json"
	"github.com/kosdirus/andintern/internal/model"
)

// Car describes car for V1CarDTO
type Car struct {
	Id    int    `json:"id"`
	Brand string `json:"brand"`
	Price int    `json:"price"`
}

func ToCar(car *model.Car) *Car {
	return &Car{
		Id:    car.Id,
		Brand: car.Brand,
		Price: car.Price,
	}
}

// CarToInsert describes car for create depending on its ID.
type CarToInsert struct {
	Id    int    `json:"id"`
	Brand string `json:"brand"`
	Price int    `json:"price"`
}

// CarToUpdate describes car for update depending on its ID.
type CarToUpdate struct {
	Id    int        `json:"id"`
	Brand JSONString `json:"brand"`
	Price JSONInt    `json:"price"`
}

type JSONInt struct {
	Value int
	Valid bool
	Set   bool
}

func (i *JSONInt) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true

	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}

	// The key isn't set to null
	var temp int
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}

type JSONString struct {
	Value string
	Valid bool
	Set   bool
}

func (s *JSONString) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	s.Set = true

	if string(data) == "null" {
		// The key was set to null
		s.Valid = false
		return nil
	}

	// The key isn't set to null
	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	s.Value = temp
	s.Valid = true
	return nil
}
