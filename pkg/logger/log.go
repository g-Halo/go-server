package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Instance *zap.SugaredLogger

func InitLogger(logPath string, loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    1 << 30, // megabytes
		MaxBackups: 30,      // 最多保留300个备份
		MaxAge:     7,       // days
		Compress:   true,    // 是否压缩 disabled by default
	}

	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和日志
		level,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	logger := zap.New(core, caller, development)

	// register to logger Instance
	Instance = logger.Sugar()

	return logger
}

func Debug(args ...interface{}) {
	Instance.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Instance.Debugf(template, args...)
}

func Info(args ...interface{}) {
	Instance.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Instance.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Instance.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Instance.Warnf(template, args...)
}

func Error(args ...interface{}) {
	Instance.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Instance.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	Instance.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Instance.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	Instance.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	Instance.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	Instance.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Instance.Fatalf(template, args...)
}