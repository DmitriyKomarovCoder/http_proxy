package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(logFilePath string) (*Logger, error) {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})

	consoleHandler := logrus.New()
	consoleHandler.SetOutput(os.Stdout)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, file)

	l.SetOutput(mw)

	return &Logger{
		Logger: l,
	}, nil
}
