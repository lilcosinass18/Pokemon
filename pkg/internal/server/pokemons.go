package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"pokemon-rest-api/listing"
	"pokemon-rest-api/pkg/internal/utils/errs"
)

func perform[T any](ctx echo.Context, okStatus int, f func(ctx context.Context) (T, error)) error {
	result, err := f(ctx.Request().Context())
	if err != nil {
		return errs.Echo(err)
	}

	return ctx.JSON(okStatus, result)
}

func (m *Server) createPokemon(ctx echo.Context) (err error) {
	var pokemon listing.Pokemon
	if err = ctx.Bind(&pokemon); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return perform(
		ctx,
		http.StatusCreated,
		func(ctx context.Context) (int, error) {
			return m.pokemons.Create(ctx, &pokemon)
		},
	)
}

func (m *Server) fetchPokemons(ctx echo.Context) (err error) {
	return perform(ctx, http.StatusOK, m.pokemons.Fetch)
}

func (m *Server) getPokemon(ctx echo.Context) (err error) {
	id, err := WrapContext(ctx).ParamInt("id")
	if err != nil {
		return err
	}

	return perform(
		ctx,
		http.StatusOK,
		func(ctx context.Context) (_ *listing.Pokemon, err error) {
			return m.pokemons.Get(ctx, id)
		},
	)
}

func (m *Server) deletePokemon(ctx echo.Context) (err error) {
	id, err := WrapContext(ctx).ParamInt("id")
	if err != nil {
		return err
	}

	return perform(
		ctx,
		http.StatusOK,
		func(ctx context.Context) (_ *listing.Pokemon, err error) {
			return m.pokemons.Delete(ctx, id)
		},
	)
}
