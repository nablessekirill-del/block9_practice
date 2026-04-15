package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func New(level string) (*logrus.Logger, func() error, error) {
	log := logrus.New()
	parsed, err := logrus.ParseLevel(level)
	if err != nil {
		parsed = logrus.InfoLevel
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join("logs", fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	log.SetLevel(parsed)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return log, logFile.Close, nil
}
