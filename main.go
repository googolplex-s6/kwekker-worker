package main

import (
	"go.uber.org/zap"
	"kwekker-worker/util"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	config, err := util.LoadConfig()
	if err != nil {
		sugaredLogger.Fatalln("Unable to load configuration; is the .env file present and valid?", err)
	}

	worker := NewWorker(sugaredLogger, *config)

	worker.Initialize()
}
