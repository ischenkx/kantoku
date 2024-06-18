package utils

import (
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"io"
	"log/slog"
)

func GetLogger(w io.Writer, service string) *slog.Logger {
	return platform.BuildServiceLogger(
		platform.BuildPrettySlogHandler(w, slog.LevelDebug),
		service,
	)
}
