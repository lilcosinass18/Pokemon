package pokemons

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"pokemon-rest-api/pkg/internal/server"
	"pokemon-rest-api/pkg/internal/utils/xcmd"
)

const (
	connString = "postgresql://postgres:postgres@localhost:5432?dbname=postgres"
	address    = "[::]:8080"
)

func Exec(ctx context.Context, logger *zap.Logger) error {
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(
		func() error {
			err := xcmd.WaitInterrupted(ctx)
			return fmt.Errorf("interruption caught: %w", err)
		},
	)
	wg.Go(
		func() error {
			logger.Info("initializing server")

			srv, err := server.NewServer(connString, logger, true)
			if err != nil {
				return fmt.Errorf("failed to initialize pokemons server: %w", err)
			}

			return srv.Serve(ctx, address)
		},
	)

	return wg.Wait()
}
