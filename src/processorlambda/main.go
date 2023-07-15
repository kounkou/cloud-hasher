package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	StatusCode      int               `json:"statusCode"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

func Handler() (Response, error) {
	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            `{"messsage":"Hello World from lambda"}`,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
