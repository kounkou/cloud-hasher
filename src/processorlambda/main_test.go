package main

import (
	// "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	event := Event{
		Nodes: []string{
			"server1",
			"server2",
			"server3",
		},
		HashKeys:    []string{"key1", "key2", "key3"},
		HashingType: "CONSISTENT_HASHING",
		Replicas:    1,
	}

	// Test with valid event
	results, err := handleRequest(event)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 3, len(results))

	// Test with empty nodes
	event.Nodes = []string{}
	results, err = handleRequest(event)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrEmptyNodes.Error())
	assert.Nil(t, results)

	// Test with empty hash keys
	event.Nodes = []string{
		"server1",
		"server2",
		"server3",
	}
	event.HashKeys = []string{}
	results, err = handleRequest(event)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrEmptyHashKeys.Error())
	assert.Nil(t, results)

	// Test with empty hashing type
	event.HashKeys = []string{"key1", "key2", "key3"}
	event.HashingType = ""
	results, err = handleRequest(event)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrUnknownHashType.Error())
	assert.Nil(t, results)

	// Test with failure creating hasher
	event.HashingType = "INVALID_HASHING_TYPE"
	results, err = handleRequest(event)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrUnknownHashType.Error())
	assert.Nil(t, results)

	// Test with invalid number of replicas
	event.HashingType = "CONSISTENT_HASHING"
	event.Replicas = -1
	results, err = handleRequest(event)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidReplicas.Error())
	assert.Nil(t, results)
}

func TestHandler(t *testing.T) {
	event := Event{
		Nodes: []string{
			"server1",
			"server2",
			"server3",
		},
		HashKeys:    []string{"9", "jacques", "test1"},
		HashingType: "CONSISTENT_HASHING",
		Replicas:    1,
	}

	// Test with valid event
	response, err := Handler(event)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "server2, server1, server1", response.Body)

	// Test with empty nodes
	event.Nodes = []string{}
	response, err = Handler(event)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, ErrEmptyNodes.Error(), response.Body)

	// Test with empty hash keys
	event.Nodes = []string{
		"server1",
		"server2",
		"server3",
	}
	event.HashKeys = []string{}
	response, err = Handler(event)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, ErrEmptyHashKeys.Error(), response.Body)

	// Test with empty hashing type
	event.HashKeys = []string{"key1", "key2", "key3"}
	event.HashingType = ""
	response, err = Handler(event)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, ErrUnknownHashType.Error(), response.Body)

	// Test with failure creating hasher
	event.HashingType = "INVALID_HASHING_TYPE"
	response, err = Handler(event)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, ErrUnknownHashType.Error(), response.Body)
}
