package worker

import (
	"context"
	kwekproto "github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	userproto "github.com/googolplex-s6/kwekker-protobufs/v3/user"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"kwekker-worker/pkg/config"
	database "kwekker-worker/pkg/db"
	"kwekker-worker/pkg/rabbitmq"
)

type Worker struct {
	logger *zap.SugaredLogger
	config config.Config
	dbconn *pgx.Conn
}

func NewWorker(logger *zap.SugaredLogger, config config.Config) *Worker {
	return &Worker{
		logger: logger,
		config: config,
	}
}

func (w *Worker) Initialize() {
	ch := make(chan proto.Message)

	rabbitMQWorker := rabbitmq.NewRabbitMQWorker(w.logger, w.config.RabbitMQ)
	go rabbitMQWorker.ListenToQueues(config.QueueList, ch)

	db := database.NewDB(w.logger, w.config.Postgres)
	w.dbconn = db.Connect()
	defer w.dbconn.Close(context.Background())

	for {
		select {
		case data := <-ch:
			switch data.(type) {
			case *kwekproto.CreateKwek:
				w.handleCreateKwek(data.(*kwekproto.CreateKwek))
			case *kwekproto.UpdateKwek:
				w.handleUpdateKwek(data.(*kwekproto.UpdateKwek))
			case *kwekproto.DeleteKwek:
				w.handleDeleteKwek(data.(*kwekproto.DeleteKwek))
			case *userproto.CreateUser:
				w.handleCreateUser(data.(*userproto.CreateUser))
			case *userproto.UpdateUser:
				w.handleUpdateUser(data.(*userproto.UpdateUser))
			case *userproto.DeleteUser:
				w.handleDeleteUser(data.(*userproto.DeleteUser))
			default:
				w.logger.Error("Unknown type received from channel")
			}
		}
	}
}
