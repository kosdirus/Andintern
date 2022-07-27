package dataprovider

import (
	"context"
	"fmt"
	"github.com/kosdirus/andintern/internal/model"
)

type CarStore interface {
	GetByFilter(ctx context.Context, filter *CarFilter) (*model.Car, error)
	GetListByFilter(ctx context.Context, filter *CarFilter) ([]*model.Car, error)
	Insert(ctx context.Context, car *model.Car) error
	Update(ctx context.Context, car *model.Car) error
	DeleteByFilter(ctx context.Context, filter *CarFilter) error
}

// CarFilter is a filter for car postgres db.
type CarFilter struct {
	ID    int
	Brand string
	Price int
}

func (f CarFilter) String() string {
	if f.ID != 0 {
		return fmt.Sprintf("id: '%d'", f.ID)
	}

	if f.Brand != "" {
		return fmt.Sprintf("brand: '%s'", f.Brand)
	}

	if f.Price != 0 {
		return fmt.Sprintf("price: '%d'", f.Price)
	}

	return fmt.Sprintf("id: %d; brand: %s; price: %d", f.ID, f.Brand, f.Price)
}

// NewCarFilter creates new instance of CarFilter.
func NewCarFilter() *CarFilter {
	return &CarFilter{}
}

// ByID filters by car.id.
func (f *CarFilter) ByID(id int) *CarFilter {
	f.ID = id
	return f
}

// ByBrand filters by car.brand.
func (f *CarFilter) ByBrand(brand string) *CarFilter {
	f.Brand = brand
	return f
}

// ByPrice filters by car.price.
func (f *CarFilter) ByPrice(price int) *CarFilter {
	f.Price = price
	return f
}
