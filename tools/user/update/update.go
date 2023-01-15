package main

import (
	"context"
	"fmt"
	"github.com/googolplex-s6/kwekker-protobufs/v3/user"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kwekker-worker/pkg/config"
	"log"
	"time"
)

const userId = "553bb4d0-332e-401f-8e9d-e44b44aa0532"

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

	err = ch.ExchangeDeclare("user-exchange", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange", err)
	}

	updateUserQueue(ch)
}

func updateUserQueue(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"user.update",
		true,
		false,
		false,
		false,
		nil,
	)

	err = ch.QueueBind(q.Name, "user.update", "user-exchange", false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue", err)
	}

	if err != nil {
		log.Fatal("Failed to declare queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	username := "newUsername"
	email := "newemail@new.com"

	newUser := &user.UpdateUser{
		UserId:    userId,
		Username:  &username,
		Email:     &email,
		UpdatedAt: timestamppb.Now(),
	}

	body, err := proto.Marshal(newUser)
	if err != nil {
		log.Fatal("Failed to marshal protobuf", err)
	}

	err = ch.PublishWithContext(ctx,
		"user-exchange",
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
