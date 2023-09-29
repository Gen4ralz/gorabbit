package main

import (
	"context"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connectToRabbitMQ()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer rabbitConn.Close()

	// We want to talk to RabbitMQ, but we need a way to communicate with it.
	ch, err := rabbitConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// We need a place to send our messages to RabbitMQ, like a mailbox.
	q, err := ch.QueueDeclare(
		"task_queue",	// name, This is the name of our mailbox.
		true,		// durable, We want our mailbox to stay forever.
		false,		// auto-delete, We don't want to delete the mailbox when it's not used.
		false,		// exclusive, Our mailbox is not just for us, others can use it too.
		false,		// no-wait, We don't want to wait when we send a message.
		nil,		// args, No special settings for our mailbox.
	)
	failOnError(err, "Failed to declare a queue")

	// Now, we want to send a message to our mailbox.
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	body := "Do some important work!"	// This is our message.
	err = ch.PublishWithContext(ctx,
		"",		// exchange, We don't need to explain where to send the message, RabbitMQ knows.
		q.Name,	// routing key, We send the message to our mailbox.
		false,	// mandatory, Set to 'false' to avoid returning undeliverable messages
		false,	// immediate, We don't need an immediate response.
		amqp.Publishing{
			ContentType: "text/plain",	// We tell RabbitMQ that our message is plain text.
			Body: []byte(body),			// We put our "Hello World!" message here.
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

// This function helps us connect to RabbitMQ.
func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts float64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ is not ready..")
			counts++
		} else {
			log.Println("Connect to RabbitMQ!")
			connection = conn
			break
		}
		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("back off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}

// This function helps us handle errors.
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}