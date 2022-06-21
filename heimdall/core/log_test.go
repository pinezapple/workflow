package core

// import (
// 	"testing"
// )

// func TestInitLoggerWithProduction(t *testing.T) {
// 	logger := InitLogger("heimdall", "debug", &ServerConfig{
// 		Environment: "PRODUCTION",
// 	})

// 	logger.Warn("Log warning from zap")
// }

// func TestInitLoggerWithDevelopment(t *testing.T) {
// 	logger := InitLogger("heimdall", "debug", &ServerConfig{
// 		Environment: "DEVELOPMENT",
// 	})

// 	logger.Warn("Log warning from zap")
// }

// func TestFatalDevelopmentLog(t *testing.T) {
// 	defer func() {
// 		if r := recover(); r == nil {
// 			logger.Error("The code did not panic")
// 		}
// 	}()

// 	logger := InitLogger("heidall", "debug", &ServerConfig{
// 		Environment: "DEVELOPMENT",
// 	})

// 	logger.DPanic("Fatal error")
// }

// func TestLoggerWithFormat(t *testing.T) {
// 	logger := InitLogger("heimdall", "debug", &ServerConfig{
// 		Environment: "DEVELOPMENT",
// 	})

// 	logger.Infof("Log formatted %v", "Hello")
// }
