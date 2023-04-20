package log

import (
	"go.uber.org/zap"
	"kando-backend/config"
)

var Logger *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Logger = logger.Sugar()
}

func SetupLogging() {
	var logger *zap.Logger
	if config.C.IsProduction() {
		logger, _ = zap.NewProduction()
	} else if config.C.IsStaging() {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	Logger = logger.Sugar()
}
