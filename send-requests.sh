#!/bin/bash

# Set script to exit immediately on error
set -e

if [ $# -eq 0 ]; then
    echo "Error: No arguments provided. Please provide the api gateway endpoint"
    exit 1
fi

curl -X POST -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" -d '{"nodes":["node1","node2"],"hashKeys":["node1"],"hashingType":"CONSISTENT_HASHING"}' | jq
