package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
)

// New for own logger configuration; need to do mo features
func New(level string) *zap.SugaredLogger {
	file, err := os.OpenFile("info.log", os.O_TRUNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sampleJSON := []byte(`{
       "level" : "` + level + `",
       "encoding": "json",
       "outputPaths":["info.log"],
       "errorOutputPaths":["stderr"],
       "encoderConfig": {
           "messageKey":"message",
           "levelKey":"level",
           "levelEncoder":"lowercase"
       }
   }`)

	var cfg zap.Config

	if err := json.Unmarshal(sampleJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

var singleLogger = New("info")

// Default if you use only this method, you guaranteed have only 1 logger
// aka singleton
func Default() *zap.SugaredLogger {
	return singleLogger
}
