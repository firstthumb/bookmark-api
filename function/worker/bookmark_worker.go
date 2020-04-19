package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

// var service bookmark.Service

// func init() {
// 	service, err := di.CreateService()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("No SQS message passed to function")
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	lambda.Start(Handler)
}
