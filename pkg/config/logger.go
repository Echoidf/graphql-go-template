package config

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	conf *Config
)

func init() {
	conf = GetConfig()
	setLogger()
}

func setLogger() {
	var zapOptions []zap.Option
	host, _ := os.Hostname()
	zapOptions = append(zapOptions,
		zap.AddCaller(),
		//zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(zap.String("hostname", host), zap.Int("pid", os.Getpid())))

	encodingConfig := zap.NewDevelopmentEncoderConfig()

	if conf.Environment != "development" {
		encodingConfig = zap.NewProductionEncoderConfig()
	}

	//encodingConfig.ConsoleSeparator = " "
	encodingConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodingConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)

	var writeSyncers []zapcore.WriteSyncer

	writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout), getLogWriter())

	// 日志级别
	level, _ := zapcore.ParseLevel(conf.Logger.Level)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encodingConfig),
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		level,
	)
	// 替换全局的日志记录器
	zap.ReplaceGlobals(zap.New(consoleCore, zapOptions...))
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.Logger.LogPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
