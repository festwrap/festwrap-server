package logging

import "log/slog"

type BaseLogger struct {
	logger *slog.Logger
}

func (l BaseLogger) Info(message string) {
	l.logger.Info(message)
}

func (l BaseLogger) Warn(message string) {
	l.logger.Warn(message)
}

func (l BaseLogger) Error(message string) {
	l.logger.Error(message)
}

func NewBaseLogger(logger *slog.Logger) BaseLogger {
	return BaseLogger{logger: logger}
}
