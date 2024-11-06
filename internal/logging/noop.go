package logging

type NoopLogger struct{}

func (l NoopLogger) Info(message string) {}

func (l NoopLogger) Warn(message string) {}

func (l NoopLogger) Error(message string) {}
