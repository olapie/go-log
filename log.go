package log

import (
	"log"

	"go.uber.org/zap"
)

type StructuredLogger = zap.Logger
type SimpleLogger = zap.SugaredLogger

// G returns the global logger
func G() *StructuredLogger {
	return zap.L()
}

// S returns the global simple logger
func S() *SimpleLogger {
	return zap.S()
}

type Stringer interface {
	LogString() string
}

func init() {
	z, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	zap.ReplaceGlobals(z)
}
