package model

type Car struct {
	Id    int    `json:"id"`
	Brand string `json:"brand" db:"brand"`
	Price int    `json:"price" db:"price"`
}
