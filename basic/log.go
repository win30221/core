package basic

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setLog() {
	logLevel, err := zap.ParseAtomicLevel(LogMode)
	if err != nil {
		log.Fatalf("取得的 log_mode 參數為 %s ，但只允許 debug, info, warn, error, dpanic, panic, fatal，請檢查 /system/log_mode 或 /service/<server_name>/log_mode 底下的配置", LogMode)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	customConfig := &zap.Config{
		Level:             logLevel,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}

	logger, err := customConfig.Build()
	if err != nil {
		log.Fatalln(err)
	}
	zap.ReplaceGlobals(logger)
}
