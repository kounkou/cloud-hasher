# cloud-hasher


[![license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kounkou/hasherprovider/blob/master/LICENSE)

#### 1. Description

Certain applications encounter challenges in efficiently allocating entities, such as events, to designated stores or nodes. 
Consider the scenario where a request needs to be directed to a server from a pool of available servers. 
The crucial question arises: which server should be selected, and what are the underlying reasons behind the choice?

I have previously addressed these inquiries in a project that can be accessed at this link: https://github.com/kounkou/hasherprovider. Within this project, I offer an educational implementation of the Hasherprovider in a cloud environment.

The objective is to enhance understanding of the allocation process and its significance within such applications.

#### 2. System Design

The following is the architectural overview of the application. Please note that this architecture is intended solely for educational purposes and does not encompass all the best practices typically employed when designing such a project.

<img width="755" alt="Screenshot 2025-05-19 at 10 15 48 PM" src="https://github.com/user-attachments/assets/9dadf7ec-c3c2-474e-9d95-d7f6e83c5ca6" />

#### 3. Installation

To install the entire stack locally, you will need to have : 

- Docker
- Localstack

#### 4. Deploy

```bash
bash install.sh
```

#### 5. Launch basic tests

Here are both **successful** and **failure scenario** `curl` commands for your README, clearly labeled for demonstration and testing:


##### ✅ Successful Request

This command sends a valid payload with all required fields:

```bash
curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1","node2"],"hashKeys":["node1"],"hashingType":"CONSISTENT_HASHING", "replicas":3}' | jq
```

**Response:**

```json
{
  "statusCode": 200,
  "body": "[\"node1\"]"
}
```

##### ❌ Failure Scenarios

###### 1. ❌ Missing `nodes`

```bash
curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":[],"hashKeys":["key1"],"hashingType":"CONSISTENT_HASHING", "replicas":3}' | jq
```

**Expected Response:**

```json
{
  "statusCode": 400,
  "body": "node list is empty"
}
```


###### 2. ❌ Missing `hashKeys`

```bash
curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1"],"hashKeys":[],"hashingType":"CONSISTENT_HASHING", "replicas":3}' | jq
```

**Expected Response:**

```json
{
  "statusCode": 400,
  "body": "hash keys list is empty"
}
```


###### 3. ❌ Invalid `hashingType`

```bash
curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1"],"hashKeys":["key1"],"hashingType":"UNKNOWN", "replicas":3}' | jq
```

**Expected Response:**

```json
{
  "statusCode": 400,
  "body": "unknown hashing type"
}
```


###### 4. ❌ Negative `replicas`

```bash
curl -X POST \
     -H "Content-Type: application/json" "https://$1.execute-api.localhost.localstack.cloud:4566/prod" \
     -d '{"nodes":["node1"],"hashKeys":["key1"],"hashingType":"CONSISTENT_HASHING", "replicas":-1}' | jq
```

**Expected Response:**

```json
{
  "statusCode": 400,
  "body": "replicas number should be positive or 0"
}
```

#### 6. Other Testing

To perform test for the stack please run the following command

```bash
$ npm test
```

Sample result 

```bash
> cloud-hasher@0.1.0 test
> jest

 PASS  test/cloud-hasher.test.ts (9.265 s)
  CloudHasherStack
    ✓ SQS Queue Created (219 ms)
    ✓ Lambda Created (114 ms)
    ✓ APIGateway Created (112 ms)

Test Suites: 1 passed, 1 total
Tests:       3 passed, 3 total
Snapshots:   0 total
Time:        9.481 s
Ran all test suites.
```
