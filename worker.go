package main

import (
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/v2/kwek"
	"go.uber.org/zap"
	"kwekker-worker/util"
)

func Initialize(logger *zap.SugaredLogger, config *util.Config) {
	createKwekChannel := make(chan *kwekker_protobufs.CreateKwek)
	updateKwekChannel := make(chan *kwekker_protobufs.UpdateKwek)
	deleteKwekChannel := make(chan *kwekker_protobufs.DeleteKwek)

	rabbitMQWorker := NewRabbitMQWorker(logger, &config.RabbitMQ)
	go rabbitMQWorker.Listen(createKwekChannel, updateKwekChannel, deleteKwekChannel)

	for {
		select {
		case createKwek := <-createKwekChannel:
			logger.Info("Received create kwek request", "kwek", createKwek)
		case updateKwek := <-updateKwekChannel:
			logger.Info("Received update kwek request", "kwek", updateKwek)
		case deleteKwek := <-deleteKwekChannel:
			logger.Info("Received delete kwek request", "kwek", deleteKwek)
		}
	}
}
