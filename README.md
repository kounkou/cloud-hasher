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

<img width="412" alt="Screenshot 2023-03-05 at 5 25 40 PM" src="https://github.com/kounkou/cloud-hasher/assets/2589171/72113ed7-f402-447a-a9e8-a41ac48075af">

#### 3. Installation

To install the entire stack locally, you will need to have : 

- Docker
- Localstack

#### 4. Usage

To perform the tests locally, after installing above dependencies, launch Docker desktop, then in another terminal, launch localstack with : 

```bash
docker run --rm -it -p 4566:4566 -p 4571:4571 -v /var/run/docker.sock:/var/run/docker.sock localstack/localstack
```

Then deploy the application

```bash
cdklocal bootstrap aws://000000000000/us-east-1 && cdklocal synth && cdklocal deploy
```

Here is a sample request JSON file containing the structure of an input.

```bash
$ cat request.json
'{
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
'
```

Here is an example request :

```bash
$ curl -X POST -H "Content-Type: application/json" https://h7p2dwjxxk.execute-api.localhost.localstack.cloud:4566/prod/ -d '{
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
```

Here is an example response to the above request :

```bash
{
  "statusCode": 200,
  "isBase64Encoded": false,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "server2, server1, server1"
}
```

#### 5. Error handling

Error handling helps the user design the code to be robust. The following errors are supported :

- Node list provided is empty

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"John"}' https://9jooqblp52.execute-api.localhost.localstack.cloud:4566/prod/ | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   131  100   116  100    15    702     90 --:--:-- --:--:-- --:--:--   823
{
  "statusCode": 400,
  "isBase64Encoded": false,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "node list is empty"
}
```

- Hashing type is empty

```bash
$ curl -X POST -H "Content-Type: application/json" https://h7p2dwjxxk.execute-api.localhost.localstack.cloud:4566/prod/ -d '{
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
  "hashingTypes": "CONSISTENT_HASHING"
}
' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   310  100   116  100   194     24     40  0:00:04  0:00:04 --:--:--    34
{
  "statusCode": 400,
  "isBase64Encoded": false,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "hash type is empty"
}
```

- Empty hashkeys

```bash
$ curl -X POST -H "Content-Type: application/json" https://h7p2dwjxxk.execute-api.localhost.localstack.cloud:4566/prod/ -d '{
  "nodes": {
    "node1": "server1",
    "node2": "server2",
    "node3": "server3"
  },
  "hashKey": [
    "key1",
    "server1",
    "server3"
  ],
  "hashingType": "CONSISTENT_HASHING"
}
' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   313  100   121  100   192    583    926 --:--:-- --:--:-- --:--:--  1638
{
  "statusCode": 400,
  "isBase64Encoded": false,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "hash keys list is empty"
}
```

#### 6. Testing

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
