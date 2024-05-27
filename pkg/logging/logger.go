package logging

import (
	"os"

	"github.com/ew-kislov/go-sample-microservice/pkg/cfg"
	"github.com/sirupsen/logrus"
)

type ContextHook struct{}

func (*ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (*ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	requestId, ok := entry.Context.Value("requestId").(string)
	if ok {
		entry.Data["requestId"] = requestId
	}

	return nil
}

func CreateLogger(config *cfg.Config) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	if config.Env == cfg.Production {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.AddHook(&ContextHook{})

	return logger
}
