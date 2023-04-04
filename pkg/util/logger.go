package util

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	prefix string
	log    *logrus.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix + ": ",
		log:    logrus.New(),
	}
}

func (l *Logger) Infof(fmt string, args ...interface{}) {
	l.log.Infof(l.prefix+fmt, args...)
}
