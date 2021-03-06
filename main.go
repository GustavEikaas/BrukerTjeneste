package main

import (
	shopify "./api_caller"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// The following code subscribes to rabbitMQ and waits for messages
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ, retrying in 5 seconds..")
		time.Sleep(5 * time.Second)
		main()
	}


	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if string(d.Body) == "update-user" {
				shopify.ShopifyCall()
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
