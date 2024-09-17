package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger/jsonformater"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger/textformatter"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultLogger *logrus.Logger = logrus.New()
)

func NewLogger(config *config.LoggerConfig) (*logrus.Logger, error) {
	logger := logrus.New()

	logger.SetFormatter(getFormatter(config))

	logFile, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(mw)

	logger.SetLevel(getLevel(config.Level))

	logger.SetReportCaller(config.ReportCaller)

	logger.Info("Init logger successfully!")

	return logger, nil
}

func InitLogger() error {
	var (
		config = config.GetConfig().LoggerConfig
	)

	defaultLogger.SetFormatter(getFormatter(&config))

	logFile := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, //days
		LocalTime:  true,
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	defaultLogger.SetOutput(mw)

	defaultLogger.SetLevel(getLevel(config.Level))

	defaultLogger.SetReportCaller(config.ReportCaller)

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

func getFormatter(config *config.LoggerConfig) logrus.Formatter {
	timeFormat := "2006-01-02T15:04:05.000"
	hostName := os.Getenv("HOSTNAME")
	if hostName == "" {
		hostName = config.Service
	}

	switch config.Format {
	case "json":
		return &jsonformater.JSONFormatter{
			PrettyPrint:     false,
			TimestampFormat: timeFormat,
			DataKey:         "metadata",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				function := f.Function
				file := f.File

				funcIdx := strings.LastIndexByte(function, '.')
				fileIdx := strings.LastIndexByte(file, '/')

				return function[funcIdx+1:], fmt.Sprintf("%s:%d", file[fileIdx+1:], f.Line)
			},
			SourceGetter: func() string {
				return hostName
			},
		}
	case "text":
		return &textformatter.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: timeFormat,
			ForceFormatting: true,
			Name:            hostName,
		}
	default:
		return &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: timeFormat,
		}
	}
}
