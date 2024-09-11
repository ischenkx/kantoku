package utils

import (
	"github.com/ischenkx/kantoku/pkg/lib/builder"
	"io"
	"log/slog"
)

func GetLogger(w io.Writer, service string) *slog.Logger {
	return builder.BuildServiceLogger(
		builder.BuildPrettySlogHandler(w, slog.LevelDebug),
		service,
	)
}
