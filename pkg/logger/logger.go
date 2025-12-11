package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

const (
	ServiceTypeDev  = "dev"
	ServiceTypeProd = "prod"
)

func Init(serviceType string, logLevel string, logFile string) error {
	var config zap.Config

	// 按照部署环境配置 zap 的配置项
	switch serviceType {
	case ServiceTypeProd:
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000")
	case ServiceTypeDev:
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000")
	default:
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000")
	}

	// 设置 Log 等级
	switch logLevel {
	case LogLevelDebug:
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case LogLevelWarn:
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LogLevelError:
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case LogLevelFatal:
		config.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	case LogLevelInfo:
		fallthrough
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 设置日志持久化目录
	if logFile != "" {
		logDir := filepath.Dir(logFile)

		err := os.MkdirAll(logDir, 0755)

		if err != nil {
			return err
		}

		config.OutputPaths = []string{"stdout", logFile}
		config.ErrorOutputPaths = []string{"stderr", logFile}
	}

	// 构建 zap 的对象
	logger, err := config.Build(zap.AddCallerSkip(1))

	if err != nil {
		return err
	}

	globalLogger = logger

	return nil
}

func GetLogger() *zap.Logger {
	if globalLogger == nil {
		logger, err := zap.NewDevelopment()

		if err != nil {
			//todo: 后续处理
			panic(err)
		}

		globalLogger = logger
	}
	return globalLogger
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
