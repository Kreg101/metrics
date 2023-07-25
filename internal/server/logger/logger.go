package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"os"
)

func New(level string) *zap.SugaredLogger {
	file, err := os.OpenFile("info.log", os.O_TRUNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
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

func Default() *zap.SugaredLogger {
	return singleLogger
}
