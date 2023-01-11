package main

import (
	"context"
	"fmt"
	"github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"kwekker-worker/pkg/config"
	"log"
	"time"
)

const kwekGuid = "f9d30d37-63a8-44a9-b2c3-3a45eb0701bc"

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

	deleteKwek(ch)
}

func deleteKwek(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"kwek.delete",
		true,
		false,
		false,
		false,
		nil,
	)

	err = ch.QueueBind(q.Name, "kwek.delete", "kweks", false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue", err)
	}

	if err != nil {
		log.Fatal("Failed to declare queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newKwek := &kwek.DeleteKwek{
		KwekGuid: kwekGuid,
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
}
