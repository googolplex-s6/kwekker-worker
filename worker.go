package main

import (
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/kwek"
	"go.uber.org/zap"
	"kwekker-worker/util"
)

func Initialize(logger *zap.Logger, config *util.Config) {
	kwekChannel := make(chan *kwekker_protobufs.Kwek, 10)

	rabbitMQWorker := NewRabbitMQWorker(logger, &config.RabbitMQ)
	go rabbitMQWorker.Listen(kwekChannel)

	for {
		kwek := <-kwekChannel
		logger.Info("Received kwek", zap.String("kwek", kwek.String()))
	}
}
