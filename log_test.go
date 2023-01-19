package log

import "testing"

func TestWriteToFile(t *testing.T) {
	ReplaceGlobal(NewLogger(func(options *Options) {
		options.Development = true
		options.ConsoleTimeHidden = true
		options.Filename = "testdata/test.log"
	}))
	Debugln("debug message")
	Infoln("info message")
	Errorln("this is another error")
}
