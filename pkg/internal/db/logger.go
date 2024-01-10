package db

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type pgxTracer struct {
	logger *zap.Logger
}

func newPGXTracer(logger *zap.Logger) *pgxTracer {
	return &pgxTracer{logger: logger}
}

func (m *pgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	marshalled, _ := json.Marshal(data.Args)

	m.logger.Info("executing query", zap.String("query", data.SQL), zap.String("args", string(marshalled)))
	return ctx
}

func (m *pgxTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		m.logger.Error("failed to execute query", zap.Error(data.Err))
		return
	}

	m.logger.Info("query executed successfully")
}
