package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	// Set up a new AWS session with the LocalStack configuration.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"), // LocalStack endpoint
	}))

	// Create an SQS service client
	svc := sqs.New(sess)

	queueURL := "http://localhost:4566/000000000000/job-queue"

	// Send message to the SQS queue
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String("Hello, LocalStack!"),
		QueueUrl:     &queueURL,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Message Sent:", *result.MessageId)
}
