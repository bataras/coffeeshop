package util

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	prefix string
	log    *logrus.Logger
}

func NewLogger(prefix string) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.99999",
	})
	log.SetOutput(os.Stdout)
	return &Logger{
		prefix: prefix + ": ",
		log:    log,
	}
}

func (l *Logger) Infof(fmt string, args ...interface{}) {
	l.log.Infof(l.prefix+fmt, args...)
}
func (l *Logger) Errorf(fmt string, args ...interface{}) {
	l.log.Errorf(l.prefix+fmt, args...)
}
