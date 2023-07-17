package main

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kounkou/hasherprovider"
)

type AlgorithmType int64

type Event struct {
	Nodes       []string `json:"nodes"`
	HashKeys    []string `json:"hashKeys"`
	HashingType string   `json:"hashingType"`
	Replicas    int      `json:"replicas"`
}

type Response struct {
	StatusCode      int               `json:"statusCode"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
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

func Handler(event Event) (Response, error) {
	results, err := handleRequest(event)
	if err != nil {
		var statusCode int
		var result string

		switch err {
		case ErrEmptyNodes, ErrEmptyHashKeys, ErrEmptyHashType, ErrFailureCreatingHasher, ErrFailureHashingKey, ErrInvalidReplicas, ErrUnknownHashType:
			statusCode = 400
			result = err.Error()
		default:
			statusCode = 500
			result = "Internal Server Error"
		}

		return Response{
			StatusCode:      statusCode,
			IsBase64Encoded: false,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            result,
		}, nil
	}

	return Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            strings.Join(results, ", "),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
