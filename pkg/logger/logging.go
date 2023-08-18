package logger

import (
	"os"

	"github.com/luanapp/gin-example/config/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	highPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
)

func NewLogger() (l *zap.Logger) {
	if env.IsProduction() {
		jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), highPriority)
		l = zap.New(core, zap.AddCaller())
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
