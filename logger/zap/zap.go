package zap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger
type apiLogger struct {
	sugarLogger *zap.SugaredLogger
	level       int
}

// App Logger constructor
func NewApiLogger(level int) *apiLogger {
	return &apiLogger{level: level}
}

// For mapping config logger to app logger levels
var loggerLevelMap = map[int]zapcore.Level{
	1: zapcore.DebugLevel,
	2: zapcore.InfoLevel,
	3: zapcore.WarnLevel,
	4: zapcore.ErrorLevel,
	5: zapcore.DPanicLevel,
	6: zapcore.PanicLevel,
	7: zapcore.FatalLevel,
}

func (l *apiLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// Init logger
func (l *apiLogger) InitLogger() {
	logLevel := l.getLoggerLevel()
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)
	logger := zap.New(core)
	l.sugarLogger = logger.Sugar()
}

// Logger methods

func (l *apiLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *apiLogger) Infow(template string, args ...interface{}) {
	l.sugarLogger.Infow(template, args...)
}

func (l *apiLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *apiLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *apiLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *apiLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *apiLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *apiLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *apiLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
