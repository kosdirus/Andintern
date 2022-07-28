package pg

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	httpv1 "github.com/kosdirus/andintern/internal/api/http/handler/v1"
	"github.com/kosdirus/andintern/internal/database"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
	"github.com/kosdirus/andintern/internal/model"
)

// CarStore is car postgres store
type CarStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

func NewCarStore(db *database.Client, txer dataprovider.Txer) dataprovider.CarStore {
	return &CarStore{
		db:        db,
		txer:      txer,
		tableName: "andintern.cars",
	}
}

func (s *CarStore) WithTx(tx *dataprovider.Tx) dataprovider.CarStore {
	return &CarStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func getCarCond(f *dataprovider.CarFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	lt := make(sq.Lt)

	var cond sq.Sqlizer = eq

	if f.ID != 0 {
		eq["cars.id"] = f.ID
	}

	if f.Brand != "" {
		eq["cars.brand"] = f.Brand
	}

	if f.Price != 0 {
		lt["cars.price"] = f.Price
		cond = lt
	}

	return cond
}

func (s CarStore) GetByFilter(ctx context.Context, filter *dataprovider.CarFilter) (*model.Car, error) {
	cars, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(cars) == 0:
		return nil, nil
	case len(cars) == 1:
		return cars[0], nil
	default:
		return nil, fmt.Errorf("fetched more than 1 car")
	}
}

func (s CarStore) GetListByFilter(ctx context.Context, filter *dataprovider.CarFilter) ([]*model.Car, error) {
	b := sq.Select(
		"cars.id",
		"cars.brand",
		"cars.price",
	).
		From(s.tableName).
		Where(getCarCond(filter))

	query, args, err := b.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("creating sql query for getting persons by filter: %w", err)
	}

	cars := make([]*model.Car, 0)
	if err = sqlx.SelectContext(ctx, s.db, &cars, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return cars, nil
		}
		return nil, fmt.Errorf("selecting cars by filter from database with query %q, error: %w", query, err)
	}

	return cars, nil
}

func (s CarStore) Insert(ctx context.Context, car *model.Car) error {
	query, args, err := sq.Insert(s.tableName).
		SetMap(map[string]interface{}{
			"brand": car.Brand,
			"price": car.Price,
		}).
		Suffix("RETURNING id;").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("can't create query SQL for inserting car: %w", err)
	}

	row := s.db.QueryRowxContext(ctx, query, args...)
	if err = row.Err(); err != nil {
		return fmt.Errorf("can't execute SQL query for inserting car: %w", err)
	}

	if err = row.Scan(&car.Id); err != nil {
		return fmt.Errorf("can't scan inserted car id: %w", err)
	}

	return nil
}

func (s CarStore) Update(ctx context.Context, car *httpv1.CarToUpdate) error {
	updates := make(map[string]interface{})
	if car.Brand.Set && car.Brand.Valid {
		updates["brand"] = car.Brand.Value
	}
	if car.Price.Set && car.Price.Valid {
		updates["price"] = car.Price.Value
	}

	query, args, err := sq.Update(s.tableName).
		SetMap(updates).
		Where(sq.Eq{"id": car.Id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("creating sql query for updating car: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't execute SQL query for updating car")
	}

	return nil
}

func (s CarStore) DeleteByFilter(ctx context.Context, filter *dataprovider.CarFilter) error {
	query, args, err := sq.
		Delete(s.tableName).
		Where(getCarCond(filter)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't create delete car query")
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't delete person from database")
	}

	return nil
}
