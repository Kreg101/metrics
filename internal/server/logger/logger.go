package logger

import "go.uber.org/zap"

func new() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	//logger.WithOptions()
	return logger.Sugar()
}

var logger = new()

func Default() *zap.SugaredLogger {
	return logger
}
