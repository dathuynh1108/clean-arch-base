package logger

import (
	"io"
	"os"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/sirupsen/logrus"
)

var (
	defaultLogger *logrus.Logger = logrus.New()
)

func InitLogger() error {
	var (
		config = config.GetConfig()
	)

	defaultLogger.SetFormatter(getFormatter(config.LoggerConfig.Format))

	logFile, err := os.OpenFile(config.LoggerConfig.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	defaultLogger.SetOutput(mw)

	defaultLogger.SetLevel(getLevel(config.LoggerConfig.Level))

	defaultLogger.Info("Init logger successfully!")
	return nil
}

func GetLogger() *logrus.Logger {
	return defaultLogger
}

func getLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	default:
		return logrus.ErrorLevel
	}
}

func getFormatter(format string) logrus.Formatter {
	switch format {
	case "json":
		return &logrus.JSONFormatter{
			PrettyPrint:     true,
			TimestampFormat: "2006-01-02T15:04:05.000",
		}
	case "text":
		return &TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000",
			ForceFormatting: true,
		}
	default:
		return &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000",
		}
	}
}
