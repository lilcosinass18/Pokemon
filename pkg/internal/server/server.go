package server

import (
	"context"
	"net"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"pokemon-rest-api/pkg/internal/db"
	"pokemon-rest-api/pkg/internal/pokemons"
	middleware2 "pokemon-rest-api/pkg/internal/server/middleware"
	"pokemon-rest-api/pkg/internal/utils/errs"
)

type Server struct {
	server   *echo.Echo
	storage  db.Storage
	pokemons pokemons.Service
	logger   *zap.Logger
}

func NewServer(dbConnString string, logger *zap.Logger, debug bool) (_ *Server, err error) {
	defer errs.Format(&err, "initialize server")

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Use(middleware.RequestID())
	server.Use(middleware.Recover())
	server.Use(middleware2.AccessLog(logger))
	if debug {
		server.Use(middleware.CORS())
	}

	server.Debug = debug

	storage, err := db.NewStorage(dbConnString, logger)
	if err != nil {
		return nil, err
	}

	s := &Server{
		server:   server,
		storage:  storage,
		pokemons: pokemons.NewService(storage, logger),
		logger:   logger,
	}
	s.initRouter()

	return s, nil
}

func (m *Server) initRouter() {
	m.server.POST("/pokemons", m.createPokemon)
	m.server.GET("/pokemons", m.fetchPokemons)
	m.server.GET("/pokemons/:id", m.getPokemon)
	m.server.DELETE("/pokemons/:id", m.deletePokemon)
}

func (m *Server) Serve(ctx context.Context, address string) error {
	if err := m.storage.CreateSchema(ctx); err != nil {
		return err
	}

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(
		func() error {
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return err
			}

			m.logger.Info("exposed HTTP server", zap.String("address", listener.Addr().String()))
			defer m.logger.Info("stopped HTTP server", zap.String("address", listener.Addr().String()))

			m.server.Listener = listener
			return m.server.Start(address)
		},
	)
	wg.Go(
		func() error {
			<-ctx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			m.logger.Info("shutting down HTTP server")
			return m.server.Shutdown(ctx)
		},
	)

	return wg.Wait()
}
