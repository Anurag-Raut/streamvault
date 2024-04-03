// // Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// // SPDX-License-Identifier: Apache-2.0

// // snippet-start:[gov2.sqs.Hello]

// package subtitle

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/sqs"

	
// )

// // main uses the AWS SDK for Go V2 to create an Amazon Simple Queue Service
// // (Amazon SQS) client and list the queues in your account.
// // This example uses the default settings specified in your shared credentials
// // and config files.
// func StartSubtitleServer() {

// 	println("starting subtitle server")
// 	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
// 	if err != nil {
// 		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
// 		fmt.Println(err)
// 		return
// 	}
// 	sqsClient := sqs.NewFromConfig(sdkConfig)
// 	fmt.Println("Let's list the queues for your account.")

// 	qUrl, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
// 		QueueName: aws.String("streamvault"),
// 	})

// 	fmt.Println(qUrl.QueueUrl)
// 	if err != nil {
// 		log.Fatalf("Got an error getting the queue URL: %v", err)
// 	}

// 	retrieve(*sqsClient, *qUrl.QueueUrl)

// }

// func retrieve(sqsClient sqs.Client, qUrl string) {

// 	for {
// 		result, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
// 			QueueUrl:            aws.String(*aws.String(qUrl)),
// 			MaxNumberOfMessages: 10,
// 			WaitTimeSeconds:     20,
// 		})

// 		if err != nil {
// 			log.Fatalf("Got an error receiving messages: %v", err)
// 		}

// 		fmt.Println("Messages:", len(result.Messages))
// 		for _, message := range result.Messages {
// 			fmt.Println("  Message ID:     " + *message.MessageId)
// 			fmt.Println("  deleting : " + *message.ReceiptHandle)
// 			// fmt.Println("  Message Body:   " + *message.Body)

// 			_, err := sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
// 				QueueUrl:      aws.String(qUrl),
// 				ReceiptHandle: message.ReceiptHandle,
// 			})
// 			if err != nil {
// 				log.Fatalf("Got an error deleting the message: %v", err)
// 			}
// 			fmt.Printf("doneee  ")

// 		}

// 	}

// }

// // snippet-end:[gov2.sqs.Hello]
