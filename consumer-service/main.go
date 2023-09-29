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
		"task_queue",	// name
		true,		// durable, Make the queue durable to survive server restarts
		false,		// delete when unused
		false,		// exclusive
		false,		// no-wait
		nil,		// arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Implement Fair Dispatch
	err = ch.Qos(
		1,     // Per consumer limit
		0,     // Per consumer global
		false, // Apply to this channel only
	)
	failOnError(err, "Failed to set QoS")

	// We're ready to receive messages from our mailbox. We sit and wait for messages to come.
	msgs, err := ch.Consume(
		q.Name,		// queue, We tell RabbitMQ which mailbox to listen to.
		"",			// consumer, We're just one of many listeners, so we don't need a special name.
		false,		// auto-ack, We don't want to automatically acknowledge messages. --> change in tutorial 2
		false,		// exclusive, We're not the only one who can listen to this mailbox.
		false,		// no-local, We don't need to hear messages we sent ourselves.
		false,		// no-wait, We don't want to wait, just listen for new messages.
		nil,		// args, No special settings for listening.
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{} // We create a special channel that will keep our program running forever.

	go func() {
		// Here, we have a little helper who listens for messages in the mailbox and tells us when they arrive.
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Simulate task processing
			time.Sleep(2 * time.Second)
			log.Printf("Done: %s", d.Body)

			// Acknowledge the message to remove it from the queue
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// We wait forever until someone stops our program (by pressing CTRL+C).
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