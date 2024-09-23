package logging

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

const (
	pkg      = "logs"
	fullPath = "logs/log.txt"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

type LogrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func NewLogger() *LogrusLogger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if _, err := os.Stat(pkg); os.IsNotExist(err) {
		err := os.Mkdir(pkg, 0755)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	return &LogrusLogger{logger: logger}
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}
