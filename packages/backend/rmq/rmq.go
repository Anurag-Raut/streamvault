package rmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var connection *amqp.Connection

func ConnectRMQ() {
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")
	var err error
	connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	

	fmt.Println("Successfully connected to RabbitMQ instance")
}

func CloseConnection() {
	connection.Close()
}

func MakeQueue(streamId string) error {
	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("error creating channel")
		return err
	}

	// Create a queue
	_, err = ch.QueueDeclare(
		streamId, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		true,     // no-wait
		nil,      // arguments
	)
	if err != nil {
		fmt.Println("error creating queue")
		return err
	}

	return nil

}

func PublishMessage(streamId string, message string) error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",       // exchange
		streamId, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		return err
	}

	return nil

}



func ConsumeMessages(streamId string)  error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(1, 0, false)
    if err != nil {
        return err
    }

	msgs, err := ch.Consume(
		streamId, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return err
	}

	for delivery := range msgs {
		fmt.Println("Hello fellas",string(delivery.Body))
		err := delivery.Ack(false)
        if err != nil {
            // Handle acknowledgment error
            fmt.Println("Error acknowledging message:", err)
        }
	}

	return  nil
}



