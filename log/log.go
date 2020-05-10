package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

const format = log.Ldate | log.Ltime | log.Lshortfile

const (
	LevelDebug = iota
	LevelInfo
	LevelError
)

var (
	level uint8
	debug *log.Logger
	info  *log.Logger
	err   *log.Logger
)

func init() {
	level = LevelInfo

	SetOutput(os.Stdout)
	SetErrOutput(os.Stderr)
}

func GetErrLogger() *log.Logger {
	return err
}

func Debug(format string, v ...interface{}) {
	if level <= LevelDebug {
		_ = debug.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if level <= LevelInfo {
		_ = info.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if level <= LevelError {
		_ = err.Output(2, fmt.Sprintf(format, v...))
	}
}

func SetLevel(l uint8) {
	level = l
}

func SetOutput(w io.Writer) {
	debug = log.New(w, "[DEBUG] ", format)
	info = log.New(w, "[INFO] ", format)
}

func SetErrOutput(w io.Writer) {
	err = log.New(w, "[ERROR] ", format)
}
