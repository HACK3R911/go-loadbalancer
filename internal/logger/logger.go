package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger struct {
	*zap.SugaredLogger
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	formatted := fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
	enc.AppendString(formatted)
}

func New(level string) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = customTimeEncoder
	//encoderConfig.MessageKey = "message"
	//encoderConfig.LevelKey = "level"
	//encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Цветные уровни

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	sugaredLogger := baseLogger.Sugar()
	return &Logger{sugaredLogger}, nil
}
func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}
