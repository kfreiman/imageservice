package logging

import (
	"log"

	"go.uber.org/zap"
)

// Logger is a simplified interface over zap.Logger. https://github.com/uber-go/zap/issues/381
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// NewLogger creates Logger
func NewLogger(isDev bool) Logger {
	config := zap.Config{}
	if isDev {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	return logger.Sugar()
}
