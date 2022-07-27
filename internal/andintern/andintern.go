package andintern

import (
	"github.com/kosdirus/andintern/internal/config"
	"github.com/kosdirus/andintern/internal/database"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
)

type Core struct {
	config   *config.Config
	db       *database.Client
	carStore dataprovider.CarStore
	txer     dataprovider.Txer
}

func NewCore(
	config *config.Config,
	db *database.Client,
	carStore dataprovider.CarStore,
	txer dataprovider.Txer,
) *Core {
	return &Core{
		config:   config,
		db:       db,
		carStore: carStore,
		txer:     txer,
	}
}

func (c Core) DB() *database.Client {
	return c.db
}
