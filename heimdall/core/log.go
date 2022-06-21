package core

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *Logger

type Logger struct {
	name string
}

func GetLogger() *Logger {
	return logger
}

func InitLogger(name, level, env string) (log *Logger) {
	fmt.Printf("Logger config: %v %v %v\n", name, level, env)
	logger = &Logger{name: name}

	var zapConfig zap.Config
	if strings.ToLower(env) == "development" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	if err := zapConfig.Level.UnmarshalText([]byte(level)); err != nil {
		panic("Can not create logger with level: " + level)
	}

	zapConfig.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	return logger
}

// Fatal log fatal error
func (lg *Logger) Fatal(message string) {
	zap.L().Fatal(message)
}

// Panic log panic error
func (lg *Logger) Panic(message string) {
	zap.L().Panic(message)
}

// DPanic log panic error and in Development
// config it then panics
func (lg *Logger) DPanic(message string) {
	zap.L().DPanic(message)
}

// Error log error message
func (lg *Logger) Error(message string) {
	zap.L().Error(message)
}

// Warn log warn message
func (lg *Logger) Warn(message string) {
	zap.L().Warn(message)
}

// Info log info message
func (lg *Logger) Info(message string) {
	zap.L().Info(message)
}

// Debug log debug message
func (lg *Logger) Debug(message string) {
	zap.L().Debug(message)
}

// Fatalf log fatal error
func (lg *Logger) Fatalf(template string, args ...interface{}) {
	zap.S().Fatalf(template, args...)
}

// Panicf log panic error
func (lg *Logger) Panicf(template string, args ...interface{}) {
	zap.S().Panicf(template, args...)
}

// DPanicf log panic error and in Development
// config it then panics
func (lg *Logger) DPanicf(template string, args ...interface{}) {
	zap.S().DPanicf(template, args...)
}

// Errorf log error message
func (lg *Logger) Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

// Warnf log warn message
func (lg *Logger) Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

// Infof log info message
func (lg *Logger) Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

// Debugf log debug message
func (lg *Logger) Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}
