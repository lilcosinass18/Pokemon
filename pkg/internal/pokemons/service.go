package pokemons

import (
	"context"

	"go.uber.org/zap"

	"pokemon-rest-api/listing"
	"pokemon-rest-api/pkg/internal/db"
)

type Service interface {
	Create(context.Context, *listing.Pokemon) (int, error)
	Fetch(context.Context) ([]*listing.Pokemon, error)
	Get(context.Context, int) (*listing.Pokemon, error)
	Delete(context.Context, int) (*listing.Pokemon, error)
}

type service struct {
	db     db.Storage
	logger *zap.Logger
}

func NewService(db db.Storage, logger *zap.Logger) Service {
	return &service{
		db:     db,
		logger: logger,
	}
}

func (m *service) Create(ctx context.Context, pokemon *listing.Pokemon) (int, error) {
	tx, err := m.db.BeginRW(ctx)
	defer db.Finish(tx, &err)

	id, err := tx.Pokemons.Insert(ctx, pokemon)
	if err != nil {
		return 0, err
	}

	return id[0], nil
}

func (m *service) Fetch(ctx context.Context) ([]*listing.Pokemon, error) {
	tx, err := m.db.BeginRO(ctx)
	defer db.Finish(tx, &err)

	pokemons, err := tx.Pokemons.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func (m *service) Get(ctx context.Context, id int) (*listing.Pokemon, error) {
	tx, err := m.db.BeginRO(ctx)
	defer db.Finish(tx, &err)

	pokemon, err := tx.Pokemons.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}

func (m *service) Delete(ctx context.Context, id int) (*listing.Pokemon, error) {
	tx, err := m.db.BeginRW(ctx)
	defer db.Finish(tx, &err)

	pokemon, err := tx.Pokemons.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}
