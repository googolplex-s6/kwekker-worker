package main

import (
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/v2/kwek"
	"go.uber.org/zap"
	"kwekker-worker/util"
)

type Worker struct {
	logger *zap.SugaredLogger
	config *util.Config
}

func NewWorker(logger *zap.SugaredLogger, config *util.Config) *Worker {
	return &Worker{
		logger: logger,
		config: config,
	}
}

func (w *Worker) Initialize() {
	createKwekChannel := make(chan *kwekker_protobufs.CreateKwek)
	updateKwekChannel := make(chan *kwekker_protobufs.UpdateKwek)
	deleteKwekChannel := make(chan *kwekker_protobufs.DeleteKwek)

	rabbitMQWorker := NewRabbitMQWorker(w.logger, &w.config.RabbitMQ)
	go rabbitMQWorker.Listen(createKwekChannel, updateKwekChannel, deleteKwekChannel)

	for {
		select {
		case createKwek := <-createKwekChannel:
			w.logger.Debug("Received create kwek request", "kwek", createKwek)
		case updateKwek := <-updateKwekChannel:
			w.logger.Debug("Received update kwek request", "kwek", updateKwek)
		case deleteKwek := <-deleteKwekChannel:
			w.logger.Debug("Received delete kwek request", "kwek", deleteKwek)
		}
	}
}
