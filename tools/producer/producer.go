package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kwekker-worker/pkg/config"
	"log"
	"math/rand"
	"time"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Unable to load configuration; is the .env file present and valid?", err)
	}

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d%s",
			conf.RabbitMQ.Username,
			conf.RabbitMQ.Password,
			conf.RabbitMQ.Host,
			conf.RabbitMQ.Port,
			conf.RabbitMQ.Vhost,
		),
	)

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal("Failed to open channel", err)
	}

	defer ch.Close()

	err = ch.ExchangeDeclare("kweks", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange", err)
	}

	createKwekQueue(ch)
}

func createKwekQueue(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"kwek.create",
		true,
		false,
		false,
		false,
		nil,
	)

	err = ch.QueueBind(q.Name, "kwek.create", "kweks", false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue", err)
	}

	if err != nil {
		log.Fatal("Failed to declare queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rand.Seed(time.Now().UnixNano())

	for {
		newKwek := &kwek.CreateKwek{
			KwekGuid: uuid.New().String(),
			Text:     "Foo bar",
			UserId:   "12abc",
			PostedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		}

		body, err := proto.Marshal(newKwek)
		if err != nil {
			log.Fatal("Failed to marshal protobuf", err)
		}

		err = ch.PublishWithContext(ctx,
			"kweks",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/protobuf",
				Body:         body,
			},
		)

		if err != nil {
			log.Fatal("Failed to publish message", err)
		}

		log.Println("Published message")

		//time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}
}
