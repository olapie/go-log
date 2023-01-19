package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const keyLogger contextKey = "logger"

func BuildContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(keyLogger).(*Logger)
	if ok {
		return l
	}
	return zap.L()
}
