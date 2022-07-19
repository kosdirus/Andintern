package model

type Car struct {
	Id    int    `json:"id,omitempty"`
	Brand string `json:"brand,omitempty" db:"brand"`
	Price int    `json:"price,omitempty" db:"price"`
}
