package v1

import (
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
