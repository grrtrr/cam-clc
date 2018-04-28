// Package logger wraps logrus for use with syslog.
package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/RackSec/srslog"
	"github.com/Sirupsen/logrus"
	"github.com/grrtrr/clccam/logger/hooks"
	"golang.org/x/crypto/ssh/terminal"
)

// GLOBAL VARIABLES
var (
	log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
)

// Init sets up the global logging defaults
// @serviceName: syslog service name, or "" to disable logging to syslog
func Init(serviceName string) {
	if serviceName != "" {
		hook, err := hooks.NewSyslogHook("", "", srslog.LOG_DEBUG|srslog.LOG_USER, serviceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: failed to initialize syslog: %s", err)
			os.Exit(-1)
		}
		log.Hooks.Add(hook)

		log.Out = ioutil.Discard
	}
}

// WriterLevel exposes the WriterLevel function of $log
func Writer() *io.PipeWriter {
	return log.WriterLevel(logrus.DebugLevel)
}

// SetLevel adjusts the global logging level to @level.
func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

func Infof(format string, args ...interface{}) {
	log.Infof(resolveLocation(format), args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(resolveLocation(format), args...)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(resolveLocation(format), args...)
}

var Warnf = Warningf // Alias

func Errorf(format string, args ...interface{}) {
	log.Errorf(resolveLocation(format), args...)
}

// Fatalf logs a formatted message and then exits the program.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(resolveLocation(format), args...)
}

// Panicf is like Fatalf, but additionally prints a stack trace.
func Panicf(format string, args ...interface{}) {
	var buf = make([]byte, 16384)

	runtime.Stack(buf, false /* do not print stacks of other goroutines */)
	log.Fatalf(resolveLocation(format)+"\n%s", append(args, string(buf))...)
}

// resolveLocation prepends file/line information
func resolveLocation(format string) string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("[%s:%d] %s", path.Base(file), line, format)
	}
	return format
}

// isTerminal returns true if @w is writing to a terminal
func isTerminal(w io.Writer) bool {
	if v, ok := w.(*os.File); ok {
		return terminal.IsTerminal(int(v.Fd()))
	}
	return false
}
