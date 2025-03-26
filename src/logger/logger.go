package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func Initialize() {
	// Setup Zap
	prodConfig := zap.NewProductionConfig()
	prodConfig.Encoding = "console"
	prodConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	prodConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	prodConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	prodConfig.EncoderConfig.CallerKey = zapcore.OmitKey
	prodConfig.EncoderConfig.StacktraceKey = zapcore.OmitKey
	log, err := prodConfig.Build()
	if err != nil {
		fmt.Printf("Error initializing logger: %s\n", err)
		os.Exit(1)
	}
	logger = log.Sugar()
}

func Get() *zap.SugaredLogger {
	return logger
}
