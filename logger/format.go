package logger

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// ANSI escape sequences for coloured logging
const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	cyan    = 36
)

// Formatter is a modified logrus.TextFormatter
type Formatter struct {
	EnableColours   bool // Enable coloured logging
	EnableTimestamp bool // Print a timestamp in front of each lime
}

// Implements logrus.Formatter
func (d *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var (
		levelColor      int
		keys            []string = make([]string, 0, len(entry.Data))
		b                        = &bytes.Buffer{}
		line                     = entry.Message
		enableTimestamp          = d.EnableTimestamp
		enableColours            = d.EnableColours
	)

	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Logger != nil && isTerminal(entry.Logger.Out) {
		enableTimestamp = true
	} else {
		enableTimestamp = false
		enableColours = false
	}

	sort.Strings(keys)

	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = cyan
	case logrus.InfoLevel:
		// black
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel: // general error
		levelColor = red
	case logrus.FatalLevel: // aborts program
		levelColor = red
	case logrus.PanicLevel: // like Fatal, but with stack trace
		levelColor = red
	default:
		levelColor = blue
	}

	if enableColours {
		if enableTimestamp {
			line = fmt.Sprintf("%s %s", time.Now().Format("15:04:05.0"), line)
		}
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", levelColor, line)
		for _, k := range keys {
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%+v", levelColor, k, entry.Data[k])
		}
	} else {
		if enableTimestamp {
			fmt.Fprintf(b, "%-5.5s %s", time.Now().Format("15:04:05.0"), line)
		} else {
			fmt.Fprintf(b, "%s", line)
		}

		for _, k := range keys {
			fmt.Fprintf(b, " %s=%+v", k, entry.Data[k])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
