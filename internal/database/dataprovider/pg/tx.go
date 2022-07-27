package pg

import (
	"fmt"
	"github.com/kosdirus/andintern/internal/database"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
)

func NewTxManager(db *database.Client) dataprovider.Txer {
	return &TxManager{
		db: db,
	}
}

type TxManager struct {
	db *database.Client
}

func (txm *TxManager) New() (*dataprovider.Tx, error) {
	sqltx, err := txm.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%w, creating tx", err)
	}

	return &dataprovider.Tx{Tx: sqltx}, nil
}
