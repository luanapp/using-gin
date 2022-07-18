package logger

import (
	"os"

	"go.uber.org/zap"
)

func NewLogger() (l *zap.Logger) {
	env := os.Getenv("USING_GIN_ENV")
	if env == "production" {
		l, _ = zap.NewProduction(zap.AddCaller())
	} else {
		l, _ = zap.NewDevelopment(zap.AddCaller())
	}

	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(l)

	return l
}

func New() *zap.SugaredLogger {
	return NewLogger().Sugar()
}
