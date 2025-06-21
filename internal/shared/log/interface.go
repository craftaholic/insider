package log

import "context"

var BaseLogger Log

type ctxKey struct{}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) Log {
	return BaseLogger.FromCtx(ctx)
}

// WithFields return a child Logger associated with
// other metadata Fields.
func WithFields(fields ...any) Log {
	return BaseLogger.WithFields(fields...)
}

type Log interface {
	FromCtx(ctx context.Context) Log
	WithCtx(ctx context.Context) context.Context
	WithFields(fields ...any) Log
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Panic(msg string, fields ...any)
	DPanic(msg string, fields ...any)
	Fatal(msg string, fields ...any)
}
