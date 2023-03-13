package logger

import (
	"schedrestd/common"
	"schedrestd/config"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"os"
	"path"
)

// AipLogger defines one new type
type AipLogger = *logrus.Logger

var defaultLogger AipLogger

// NewLogger for DI
func NewLogger(conf *config.Config) AipLogger {

	hostname := os.Getenv(common.HostName)
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	logFileName := fmt.Sprintf("%s.log.%s", common.ConfigFileName, hostname)

	fileName := path.Join(conf.LogDir, logFileName)

	// Create new logger
	logger := logrus.New()

	// Set the rotate rule
	writer, _ := rotatelogs.New(
		fileName + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithRotationSize(52428800),
		rotatelogs.WithRotationCount(10),
	)
	logger.SetOutput(writer)

	// Set the log level
	logger.SetLevel(getLogLevel(conf.LogLevel))

	// Set the
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	return logger
}

func getLogLevel(level string) logrus.Level {
	switch level {
	case common.Fatal:
		return logrus.FatalLevel
	case common.Error:
		return logrus.ErrorLevel
	case common.Warn:
		return logrus.WarnLevel
	case common.Info:
		return logrus.InfoLevel
	case common.Debug:
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}

// SetDefault sets the default logger for testing and outside fx container
func SetDefault() (AipLogger, error) {
	log := NewLogger(&config.Config{LogDir: config.GetLogDir()})
	if log == nil {
		return nil, fmt.Errorf("failed to set default logger")
	}

	defaultLogger = log
	return defaultLogger, nil
}

// GetDefault gets the default logger for testing and outside fx container
func GetDefault() AipLogger {
	if defaultLogger == nil {
		// actually only reach in test
		if _, err := SetDefault(); err != nil {
			fmt.Printf("Failed to get logger: %v", err.Error())
		}
	}
	return defaultLogger
}

// GetLogger gets the gin logger
func GetLogger(c *gin.Context) AipLogger {
	log, ok := c.Get(common.LoggerName)
	if ok {
		switch log.(type) {
		case AipLogger:
			return log.(AipLogger)
		}
	}
	return defaultLogger
}
