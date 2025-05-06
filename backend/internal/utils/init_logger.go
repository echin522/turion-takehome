package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger creates a production logger. Less pretty, more speedy
func InitLogger() *zap.Logger {
	return zap.Must(zap.NewProduction())
}

// InitDevLogger creates a prettified logger. Lots of pretty colors
func InitDevLogger() *zap.Logger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.Must(zapConfig.Build())
}
