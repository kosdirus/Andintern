package andintern

import (
	"github.com/kosdirus/andintern/internal/config"
	"github.com/kosdirus/andintern/internal/database"
)

type Core struct {
	config *config.Config
	db     *database.Client
}

func NewCore(config *config.Config, db *database.Client) *Core {
	return &Core{
		config: config,
		db:     db,
	}
}

func (c Core) DB() *database.Client {
	return c.db
}
