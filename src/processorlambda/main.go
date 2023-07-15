package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string `json:"name"`
}

type Response struct {
	StatusCode      int               `json:"statusCode"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

func Handler(event Event) (Response, error) {

	name := event.Name

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            fmt.Sprintf("Hello %s!", name),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
