package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

type Logger = zap.Logger
type SugaredLogger = zap.SugaredLogger

var globalMu sync.Mutex
var globalLogger *Logger
var globalSugaredLogger *SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ReplaceGlobal(l)
}

func ReplaceGlobal(l *Logger) {
	globalMu.Lock()
	zap.ReplaceGlobals(l)
	globalLogger = zap.L()
	globalSugaredLogger = zap.S()
	globalMu.Unlock()
}

// G returns the global logger
func G() *Logger {
	return globalLogger
}

// S returns the global simple logger
func S() *SugaredLogger {
	return globalSugaredLogger
}

type Options struct {
	Development bool
	TimeHidden  bool
	Filename    string
}

func NewLogger(optFns ...func(*Options)) *Logger {
	options := &Options{}
	for _, fn := range optFns {
		fn(options)
	}

	var config zap.Config
	var zapOptions []zap.Option
	if options.Development {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	var encoderConfig zapcore.EncoderConfig
	if options.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	if options.TimeHidden {
		encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {}
		config.EncoderConfig = encoderConfig
	}

	if options.Filename != "" {
		writer := NewFileWriter(options.Filename)
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		zapOptions = append(zapOptions, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				core,
				zapcore.NewCore(encoder, zapcore.AddSync(writer), config.Level),
			)
		}))
	}

	l, err := config.Build(zapOptions...)
	if err != nil {
		panic(err)
	}
	return l
}

type Func = func(msg string, fields ...Field)

type Stringer interface {
	LogString() string
}

type StdLogger interface {
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
}

func Debugln(args ...interface{}) {
	globalSugaredLogger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	globalSugaredLogger.Infoln(args...)
}

func Warnln(args ...interface{}) {
	globalSugaredLogger.Warnln(args...)
}

func Errorln(args ...interface{}) {
	globalSugaredLogger.Errorln(args...)
}

func Panicln(args ...interface{}) {
	globalSugaredLogger.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	globalSugaredLogger.Fatalln(args...)
}

func Debugf(template string, args ...interface{}) {
	globalSugaredLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	globalSugaredLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	globalSugaredLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	globalSugaredLogger.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	globalSugaredLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	globalSugaredLogger.Fatalf(template, args...)
}
