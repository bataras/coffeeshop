package util

import (
	"github.com/sirupsen/logrus"
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
	return &Logger{
		prefix: prefix + ": ",
		log:    log,
	}
}

func (l *Logger) Infof(fmt string, args ...interface{}) {
	l.log.Infof(l.prefix+fmt, args...)
}
