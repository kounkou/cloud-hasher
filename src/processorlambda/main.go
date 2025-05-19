package main

import (
	"encoding/json"
	"errors"
	"fmt"

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
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid JSON: %v", err),
		}, nil
	}

	results, err := handleRequest(event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error: %v", err),
		}, nil
	}

	// Convert results to JSON
	responseJSON, _ := json.Marshal(results)

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
