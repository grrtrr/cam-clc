package hooks

import (
	"fmt"
	"os"

	"github.com/RackSec/srslog"
	"github.com/sirupsen/logrus"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer        *srslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", srslog.LOG_DEBUG, "")`
func NewSyslogHook(network, raddr string, priority srslog.Priority, tag string) (*SyslogHook, error) {
	w, err := srslog.Dial(network, raddr, priority, tag)
	if err != nil {
		return nil, err
	}

	w.SetFormatter(func(p srslog.Priority, hostname, tag, content string) string {
		return fmt.Sprintf("<%d>%s: %s", p, tag, content)
	})
	return &SyslogHook{w, network, raddr}, err
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit(line)
	case logrus.FatalLevel:
		return hook.Writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.Writer.Err(line)
	case logrus.WarnLevel:
		return hook.Writer.Warning(line)
	case logrus.InfoLevel:
		return hook.Writer.Info(line)
	case logrus.DebugLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
