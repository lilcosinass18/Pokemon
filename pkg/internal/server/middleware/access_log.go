package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func AccessLog(logger *zap.Logger) echo.MiddlewareFunc {
	logFields := func(now time.Time, ctx echo.Context) []zap.Field {
		fields := []zap.Field{
			zap.Int("status", ctx.Response().Status),
			zap.String("method", ctx.Request().Method),
			zap.String("url", ctx.Request().URL.String()),
			zap.String("remote_addr", ctx.RealIP()),
			zap.Int64("response_size", ctx.Response().Size),
			zap.Duration("request_time", time.Since(now)),
		}

		requestID := ctx.Request().Header.Get(echo.HeaderXRequestID)
		if requestID == "" {
			requestID = ctx.Response().Header().Get(echo.HeaderXRequestID)
		}
		if len(requestID) > 0 {
			fields = append(fields, zap.String("request_id", requestID))
		}

		return fields
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			now := time.Now()
			if err := next(ctx); err != nil {
				ctx.Error(err)

				if ctx.Response().Status/100 == 5 {
					logger.Error("request completed", logFields(now, ctx)...)
				} else {
					logger.Warn("request completed", logFields(now, ctx)...)
				}
				return err
			}

			logger.Info("request completed", logFields(now, ctx)...)
			return nil
		}
	}
}
