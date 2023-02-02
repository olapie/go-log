package log

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger

var globalMu sync.Mutex
var globalLogger *Logger
var globalSugaredLogger *zap.SugaredLogger
var globalSugaredLogger2 *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ReplaceGlobal(l)
}

func ReplaceGlobal(l *Logger) {
	if l == nil {
		panic("logger is nil")
	}
	globalMu.Lock()
	zap.ReplaceGlobals(l)
	globalLogger = zap.L()
	globalSugaredLogger = zap.S()
	globalSugaredLogger2 = zap.S().WithOptions(zap.AddCallerSkip(1))
	globalMu.Unlock()
}

// G returns the global logger
func G() *Logger {
	return globalLogger
}

type FileEncoding int

const (
	FileEncodingJSON = iota
	FileEncodingConsole
)

type Options struct {
	Development       bool
	ConsoleTimeHidden bool
	Console           bool // output to console
	Filename          string
	FileEncoder       FileEncoding
	MaxFileSize       int // megabytes
	MaxFileAge        int // days
	MaxFileBackups    int
}

func NewLogger(optFns ...func(*Options)) *Logger {
	options := &Options{
		Console:        true,
		MaxFileSize:    256, // 256MB
		MaxFileAge:     14,  // 14 days
		MaxFileBackups: 32,
	}
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

	if options.ConsoleTimeHidden {
		config.EncoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {}
	}

	if options.Filename != "" {
		writer := NewFileWriter(options.Filename, func(opts *FileWriterOptions) {
			opts.MaxSize = int64(options.MaxFileSize) * 1024 * 1024
			opts.MaxAge = time.Hour * 24 * time.Duration(options.MaxFileAge)
			opts.MaxBackups = options.MaxFileBackups
		})
		var encoder zapcore.Encoder
		if options.FileEncoder == FileEncodingConsole {
			encoder = zapcore.NewConsoleEncoder(config.EncoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(config.EncoderConfig)
		}
		zapOptions = append(zapOptions, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			fileCore := zapcore.NewCore(encoder, zapcore.AddSync(writer), config.Level)
			if options.Console {
				return zapcore.NewTee(core, fileCore)
			} else {
				return fileCore
			}
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
	globalSugaredLogger2.Debugln(args...)
}

func Infoln(args ...interface{}) {
	globalSugaredLogger2.Infoln(args...)
}

func Warnln(args ...interface{}) {
	globalSugaredLogger2.Warnln(args...)
}

func Errorln(args ...interface{}) {
	globalSugaredLogger2.Errorln(args...)
}

func Panicln(args ...interface{}) {
	globalSugaredLogger2.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	globalSugaredLogger2.Fatalln(args...)
}

func Debugf(template string, args ...interface{}) {
	globalSugaredLogger2.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	globalSugaredLogger2.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	globalSugaredLogger2.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	globalSugaredLogger2.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	globalSugaredLogger2.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	globalSugaredLogger2.Fatalf(template, args...)
}
