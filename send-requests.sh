#!/bin/bash

# Set script to exit immediately on error
set -e

if [ $# -eq 0 ]; then
    echo "Error: No arguments provided. Please provide the api gateway endpoint"
    exit 1
fi

# request 1 : add servers

curl -X POST -H "Content-Type: application/json" $1 -d '{
  "nodes": {
    "node1": "server1",
    "node2": "server2",
    "node3": "server3"
  },
  "hashKeys": [
    "key1",
    "server1",
    "server3"
  ],
  "hashingType": "CONSISTENT_HASHING"
}
' | jq
