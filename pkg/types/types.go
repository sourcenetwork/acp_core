package types

var _ Logger = (*noopLogger)(nil)

type Logger interface {
	Info(msg string, args ...any)

	Warn(msg string, args ...any)

	Error(msg string, args ...any)

	Debug(msg string, args ...any)
}

type noopLogger struct{}

func (l *noopLogger) Info(msg string, args ...any) {}

func (l *noopLogger) Warn(msg string, args ...any) {}

func (l *noopLogger) Error(msg string, args ...any) {}

func (l *noopLogger) Debug(msg string, args ...any) {}

func NoopLogger() Logger { return &noopLogger{} }
