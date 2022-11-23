package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const keyLogger contextKey = "logger"

func BuildContext(ctx context.Context, l *StructuredLogger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func FromContext(ctx context.Context) *StructuredLogger {
	l, ok := ctx.Value(keyLogger).(*StructuredLogger)
	if ok {
		return l
	}
	return zap.L()
}
