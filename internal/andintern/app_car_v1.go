package andintern

import (
	"context"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
	"github.com/kosdirus/andintern/internal/model"
)

func (c *Core) GetCars(ctx context.Context, carFilter *dataprovider.CarFilter) ([]*model.Car, error) {
	return c.carStore.GetListByFilter(ctx, carFilter)
}
