package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(mode string) error {
	var cfg zap.Config

	if mode == "release" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
