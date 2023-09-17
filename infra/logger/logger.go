package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Initialize(isDebug bool) {
	var logger *zap.Logger
	if isDebug {
		logger = NewDevelopmentLogger()
	} else {
		logger = NewProductionLogger()
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}

func NewProductionLogger() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		// Filename:   "/var/log/mxshop/userweb/info.log",
		Filename:   "./log/info.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "ts",
			CallerKey:     "caller",
			StacktraceKey: "trace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendInt64(int64(d) / 1000000)
			},
		}),
		w,
		zap.InfoLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func NewDevelopmentLogger() *zap.Logger {
	w := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "ts",
			CallerKey:     "caller",
			StacktraceKey: "trace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendInt64(int64(d) / 1000000)
			},
		}),
		w,
		zap.DebugLevel,
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
