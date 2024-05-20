package builder

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/logging/prefixed"
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"time"
)

func newLogger(w io.Writer) *slog.Logger {
	coloredHandler := tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})

	prefixedHandler := prefixed.NewHandler(coloredHandler,
		&prefixed.HandlerOptions{
			PrefixKeys: []string{"service"},
		})
	return slog.New(prefixedHandler)
}

func withLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func extractLogger(ctx context.Context, defaultLogger *slog.Logger) *slog.Logger {
	val := ctx.Value("logger")
	if val != nil {
		if logger, ok := val.(*slog.Logger); ok {
			return logger
		}
	}
	return defaultLogger
}
