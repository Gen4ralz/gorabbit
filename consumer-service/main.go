package main

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connectToRabbitMQ()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer rabbitConn.Close()

	// Declare Channel
	ch, err := rabbitConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare Queue
	q, err := ch.QueueDeclare(
		"hello",	// name
		false,		// durable
		false,		// delete when unused
		false,		// exclusive
		false,		// no-wait
		nil,		// arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,		// queue
		"",			// consumer
		true,		// auto-ack
		false,		// exclusive
		false,		// no-local
		false,		// no-wait
		nil,		// args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}