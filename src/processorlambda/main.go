package main

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kounkou/hasherprovider"
)

var (
	ErrEmptyNodes            = errors.New("node list is empty")
	ErrEmptyHashKeys         = errors.New("hash keys list is empty")
	ErrEmptyHashType         = errors.New("hash type is empty")
	ErrFailureCreatingHasher = errors.New("hasher creation failure")
	ErrFailureHashingKey     = errors.New("hash failure")
)

type AlgorithmType int64

type Event struct {
	Nodes       map[string]string `json:"nodes"`
	HashKeys    []string          `json:"hashKeys"`
	HashingType string            `json:"hashingType"`
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

func handleRequest(event Event) (string, error) {

	hasherprovider := hasherprovider.HasherProvider{}

	h, err := hasherprovider.GetHasher(m[event.HashingType])

	if h == nil || err != nil {
		return "", ErrFailureCreatingHasher
	}

	h.SetReplicas(1)

	for v := range event.Nodes {
		h.AddNode(v)
	}

	var r []string

	for _, v := range event.HashKeys {
		t, err := h.Hash(v, 0)

		if err != nil {
			return "", ErrFailureHashingKey
		}

		r = append(r, t)
	}

	if len(event.Nodes) == 0 {
		return "", ErrEmptyNodes
	} else if len(event.HashKeys) == 0 {
		return "", ErrEmptyHashKeys
	} else if event.HashingType == "" {
		return "", ErrEmptyHashType
	}

	return strings.Join(r, ", "), nil
}

func Handler(event Event) (Response, error) {

	result, err := handleRequest(event)

	statusCode := 200

	if err != nil {
		if err == ErrEmptyNodes {
			statusCode = 400
			result = ErrEmptyNodes.Error()
		} else if err == ErrEmptyHashKeys {
			statusCode = 400
			result = ErrEmptyHashKeys.Error()
		} else if err == ErrEmptyHashType {
			statusCode = 400
			result = ErrEmptyHashType.Error()
		} else if err == ErrFailureCreatingHasher {
			statusCode = 400
			result = ErrFailureCreatingHasher.Error()
		}
	}

	return Response{
		StatusCode:      int(statusCode),
		IsBase64Encoded: false,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            result,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
