package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kounkou/hasherprovider"
)

type Event struct {
	Nodes       []string `json:"nodes"`
	HashKeys    []string `json:"hashKeys"`
	HashingType string   `json:"hashingType"`
	Replicas    int      `json:"replicas"`
}

var m = map[string]int{
	"CONSISTENT_HASHING": 0,
	"RANDOM_HASHING":     1,
	"UNIFORM_HASHING":    2,
}

var (
	ErrEmptyNodes            = errors.New("node list is empty")
	ErrEmptyHashKeys         = errors.New("hash keys list is empty")
	ErrEmptyHashType         = errors.New("hash type is empty")
	ErrFailureCreatingHasher = errors.New("hasher creation failure")
	ErrFailureHashingKey     = errors.New("hash failure")
	ErrInvalidReplicas       = errors.New("replicas number should be positive or 0")
	ErrUnknownHashType       = errors.New("unknown hashing type")
)

type JsonResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func handleRequest(event Event) ([]string, error) {
	_, ok := m[event.HashingType]

	if !ok {
		return nil, ErrUnknownHashType
	} else if len(event.Nodes) == 0 {
		return nil, ErrEmptyNodes
	} else if len(event.HashKeys) == 0 {
		return nil, ErrEmptyHashKeys
	} else if event.HashingType == "" {
		return nil, ErrEmptyHashType
	} else if event.Replicas < 0 {
		return nil, ErrInvalidReplicas
	}

	hasherProvider := hasherprovider.HasherProvider{}
	hasher, err := hasherProvider.GetHasher(m[event.HashingType])
	if err != nil || hasher == nil {
		return nil, ErrFailureCreatingHasher
	}

	hasher.SetReplicas(event.Replicas)

	for _, v := range event.Nodes {
		hasher.AddNode(v)
	}

	var results []string
	for _, node := range event.Nodes {
		hasher.AddNode(node)
	}

	for _, key := range event.HashKeys {
		hashedKey, err := hasher.Hash(key, 0)
		if err != nil {
			return nil, ErrFailureHashingKey
		}
		results = append(results, hashedKey)
	}

	return results, nil
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event Event

	err := json.Unmarshal([]byte(req.Body), &event)
	if err != nil {
		responseBody, _ := json.Marshal(JsonResponse{
			StatusCode: 400,
			Body:       err.Error(),
		})

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	results, err := handleRequest(event)
	if err != nil {
		responseBody, _ := json.Marshal(JsonResponse{
			StatusCode: 400,
			Body:       err.Error(),
		})

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Convert results to JSON
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		responseBody, _ := json.Marshal(JsonResponse{
			StatusCode: 400,
			Body:       err.Error(),
		})

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(responseBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	jsonString := string(jsonBytes)

	responseJSON, _ := json.Marshal(JsonResponse{
		StatusCode: 200,
		Body:       string(jsonString),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
