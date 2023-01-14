package main

import (
	"go.uber.org/zap"
	"kwekker-worker/pkg/config"
	"kwekker-worker/pkg/worker"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	conf, err := config.LoadConfig()
	if err != nil {
		sugaredLogger.Fatalln("Unable to load configuration; is the .env file present and valid?", err)
	}

	w := worker.NewWorker(sugaredLogger, *conf)
	w.Initialize()
}
