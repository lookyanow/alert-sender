package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"log"
)

// NewLogger method wraps zap logger
func NewLogger() *zap.Logger {
	zapCfg := zap.NewProductionConfig()
	//zapCfg.DisableCaller = true // we want see code line of a message
	//zapCfg.DisableStacktrace = true // we want stacktrace for Errors and above
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg.EncoderConfig.MessageKey = "message"
	zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatal("Initialize logger failed", err)
	}
	return logger
}
