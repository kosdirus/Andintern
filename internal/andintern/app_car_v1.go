package andintern

import (
	"context"
	"fmt"
	httpv1 "github.com/kosdirus/andintern/internal/api/http/handler/v1"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
	"github.com/kosdirus/andintern/internal/model"
)

func (c *Core) GetCars(ctx context.Context, carFilter *dataprovider.CarFilter) ([]*model.Car, error) {
	return c.carStore.GetListByFilter(ctx, carFilter)
}

func (c *Core) GetCar(ctx context.Context, carFilter *dataprovider.CarFilter) (*model.Car, error) {
	return c.carStore.GetByFilter(ctx, carFilter)
}

func (c *Core) CreateCar(ctx context.Context, carToCreate *httpv1.CarToInsert) (*model.Car, error) {
	car := &model.Car{
		Brand: carToCreate.Brand,
		Price: carToCreate.Price,
	}

	if err := c.carStore.Insert(ctx, car); err != nil {
		return nil, err
	}

	return car, nil
}

func (c *Core) UpdateCar(ctx context.Context, carToUpdate *httpv1.CarToUpdate) (*model.Car, error) {
	if err := c.carStore.Update(ctx, carToUpdate); err != nil {
		return nil, err
	}

	car, err := c.GetCar(ctx, dataprovider.NewCarFilter().ByID(carToUpdate.Id))
	if err != nil {
		return nil, err
	}

	return car, nil
}

func (c *Core) DeleteCar(ctx context.Context, filter *dataprovider.CarFilter) error {
	carsToDelete, err := c.carStore.GetListByFilter(ctx, filter)
	if err != nil {
		return err
	}

	if len(carsToDelete) == 0 {
		return fmt.Errorf("in database no cars with such parameter: %s", filter)
	}

	err = c.carStore.DeleteByFilter(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
