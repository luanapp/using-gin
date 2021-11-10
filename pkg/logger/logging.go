package logger

import (
	"go.uber.org/zap"
)

func New() *zap.SugaredLogger {
	//l, _ := zap.NewProduction(zap.AddCaller())
	l, _ := zap.NewDevelopment(zap.AddCaller())
	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(l)

	return l.Sugar()
}
