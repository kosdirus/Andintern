package dataprovider

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	*sqlx.Tx
}

type Txer interface {
	New() (*Tx, error)
}

func EndTransaction(tx *Tx, err error) error {
	if err == nil {
		if cerr := tx.Commit(); cerr != nil {
			return fmt.Errorf("can't commit transaction: %w", cerr)
		}
		//log.Println("commit OK")

		return nil
	}

	//log.Println("rolling back transaction")
	if rerr := tx.Rollback(); rerr != nil {
		return fmt.Errorf("can't roll back transaction: %w", rerr)
	}

	return err
}
