package repo

import log "github.com/sirupsen/logrus"

type gormLogger struct{}

func (l *gormLogger) Printf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}
