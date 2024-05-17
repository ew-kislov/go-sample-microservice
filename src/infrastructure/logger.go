package infrastructure

import (
	"os"

	"github.com/sirupsen/logrus"
)

type ContextHook struct{}

func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	requestId, ok := entry.Context.Value("requestId").(string)
	if ok {
		entry.Data["requestId"] = requestId
	}
	return nil
}

func CreateLogger(config Config) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	if config.Env == Production {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.AddHook(&ContextHook{})

	return logger
}
