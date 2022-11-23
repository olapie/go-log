package log

import (
	"go.uber.org/zap"
	"log"
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

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
}
