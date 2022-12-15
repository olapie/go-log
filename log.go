package log

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"

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

// InitLoggerWithoutTime removes time from log as CloudWatch already includes time
func InitLoggerWithoutTime(development bool) {
	var config zap.Config
	if development {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {

		}
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig = encoderConfig
	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {

		}
		config = zap.NewProductionConfig()
		config.EncoderConfig = encoderConfig
	}

	logger, err := config.Build()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
}

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
}
