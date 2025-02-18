package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

// Init 初始化日志
func Init(isDebug bool, logPath string) {
	writeSyncer := getLogWriter(isDebug, logPath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(isDebug bool, logPath string) zapcore.WriteSyncer {
	if isDebug {
		return zapcore.AddSync(os.Stdout)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename: logPath,
		Compress: false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func Debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	if ce := myCheck(zap.DebugLevel, msg); ce != nil {
		ce.Write()
	}
}

// Info logs interface in Info loglevel.
func Info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	if ce := myCheck(zap.InfoLevel, msg); ce != nil {
		ce.Write()
	}
}

// Warn logs interface in warning loglevel
func Warn(v ...interface{}) {
	msg := fmt.Sprint(v...)
	if ce := myCheck(zap.WarnLevel, msg); ce != nil {
		ce.Write()
	}
}

// Error logs interface in Error loglevel
func Error(v ...interface{}) {
	msg := fmt.Sprint(v...)
	if ce := myCheck(zap.ErrorLevel, msg); ce != nil {
		ce.Write()
	}
}

// Debugf logs interface in debug loglevel with formating string
func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if ce := myCheck(zap.DebugLevel, msg); ce != nil {
		ce.Write()
	}
}

// Infof logs interface in Infof loglevel with formating string
func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if ce := myCheck(zap.InfoLevel, msg); ce != nil {
		ce.Write()
	}
}

// Warnf logs interface in warning loglevel with formating string
func Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if ce := myCheck(zap.WarnLevel, msg); ce != nil {
		ce.Write()
	}
}

// Errorf logs interface in Error loglevel with formating string
func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if ce := myCheck(zap.ErrorLevel, msg); ce != nil {
		ce.Write()
	}
}

// myCheck 添加go id
func myCheck(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	msg = GetGoID() + "  " + msg
	return logger.Check(lvl, msg)
}
