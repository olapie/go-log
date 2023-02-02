package log

import (
	"github.com/natefinch/lumberjack/v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ io.WriteCloser = (*FileWriter)(nil)

type FileWriterOptions struct {
	lumberjack.Options
	MaxSize int64
}

type FileWriter struct {
	ll *lumberjack.Roller
}

func NewFileWriter(filename string, optFns ...func(options *FileWriterOptions)) *FileWriter {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = filepath.Join(os.Args[0], "log")
	}

	opts := &FileWriterOptions{
		Options: lumberjack.Options{
			MaxBackups: 32,
			MaxAge:     30 * time.Hour * 24, // 28 days
			LocalTime:  false,
			Compress:   true,
		},
		MaxSize: 512 * 1024 * 1024,
	}

	for _, fn := range optFns {
		fn(opts)
	}

	ll, err := lumberjack.NewRoller(filename, opts.MaxSize, &opts.Options)

	if err != nil {
		log.Fatalln(err)
	}

	return &FileWriter{
		ll: ll,
	}
}

func (f *FileWriter) Close() error {
	return f.ll.Close()
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
	return f.ll.Write(p)
}
