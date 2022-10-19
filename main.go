package main

import (
	"go.uber.org/zap"
	"kwekker-worker/util"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()

	_, err := util.LoadConfig()
	if err != nil {
		sugaredLogger.Fatalln("Unable to load configuration; is the .env file present and valid?", err)
	}
}
