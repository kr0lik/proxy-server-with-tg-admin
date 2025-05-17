package socks5

import (
	"fmt"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func (cl *Logger) Errorf(format string, v ...interface{}) {
	cl.logger.Info("Socks5.server", "internal", fmt.Sprintf(format, v...))
}
