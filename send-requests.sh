#!/bin/bash

# Set script to exit immediately on error
# set -e

if [ $# -eq 0 ]; then
    echo "Error: No arguments provided. Please provide the api gateway endpoint"
    exit 1
fi

echo "✅ Successful Request"
echo "This command sends a valid payload with all required fields:"

curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1","node2"],"hashKeys":["node1"],"hashingType":"CONSISTENT_HASHING", "replicas":3}' | jq

echo "❌ Failure Scenarios"
echo "1. ❌ Missing nodes"

curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":[],"hashKeys":["key1"],"hashingType":"CONSISTENT_HASHING", "replicas":3}' | jq

echo ""
echo "2. ❌ Missing hashKeys"

curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1"],"hashKeys":[],"hashingType":"CONSISTENT_HASHING", "replicas":3}' | jq

echo ""
echo "3. ❌ Invalid hashingType"

curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1"],"hashKeys":["key1"],"hashingType":"UNKNOWN", "replicas":3}' | jq

echo ""
echo "4. ❌ Negative replicas"

curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1"],"hashKeys":["key1"],"hashingType":"CONSISTENT_HASHING", "replicas":-1}' | jq

echo ""
