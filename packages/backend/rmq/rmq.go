package rmq

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	// "path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

var pool *pgxpool.Pool
var connection *amqp.Connection

func ConnectRMQ() {
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")
	var err error
	connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	pool, err = pgxpool.Connect(context.Background(), "host=database user=postgres password=postgres dbname=streamvault sslmode=disable")

	if err != nil {
		fmt.Println("Error connecting to database")
	}

	err = MakeQueue("vods")
	if err != nil {
		fmt.Println("error creating queue")
	}

	fmt.Println("Successfully connected to RabbitMQ instance")
}

func CloseConnection() {
	connection.Close()
}

func MakeQueue(queuename string) error {
	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("error creating channel")
		return err
	}

	// Create a queue
	_, err = ch.QueueDeclare(
		queuename, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Println("error creating queue")
		return err
	}

	return nil

}

func PublishMessage(videoId string, queueName string) error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(videoId),
		})

	if err != nil {
		return err
	}

	return nil

}

func ConsumeMessages(queueName string) error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	for delivery := range msgs {
		fmt.Println("Hello fellas", string(delivery.Body))
		videoId := string(delivery.Body)

		err := VideoIdToHls(videoId)
		if err != nil {
			fmt.Println("Error converting video to HLS:", err)
		}

		err = delivery.Ack(false)
		if err != nil {
			// Handle acknowledgment error
			fmt.Println("Error acknowledging message:", err)
		}
	}

	return nil
}

func VideoIdToHls(videoId string) error {
	dirPath := fmt.Sprintf("/home/anurag/s3mnt/%s", videoId)
	filePattern := fmt.Sprintf("/home/anurag/s3mnt/vod/%s.mkv", videoId)
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	fmt.Println(filePattern, "file Pattern")
	cmd := exec.Command("ffmpeg",
		"-i", filePattern, // Input file path
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
		"-f", "hls",
		"-g", "20",
		"-hls_time", "5",
		"-hls_list_size", "0",
		"-progress", "pipe:1",
		fmt.Sprintf("%s/%s.m3u8", dirPath, videoId), // Output HLS file path
	)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing ffmpeg command:", err)
		return err
	}

	err = UpdateVodStatus(videoId)
	if err != nil {
		fmt.Println("Error updating vod status:", err)
		return err
	}

	return nil

}

func UpdateVodStatus(videoId string) error {
	_, err := pool.Exec(context.Background(), `
	UPDATE "Video"
	SET "isProcessed" = true
	WHERE "id" = $1
	`, videoId)
	if err != nil {
		return err
	}
	return nil
}
