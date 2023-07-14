# cloud-hasher

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
docker run --rm -it -p 4566:4566 -p 4571:4571 localstack/localstack
```

Then deploy the application

```bash
cdklocal bootstrap aws://000000000000/us-east-1
cdklocal synth
cdklocal deploy
```

Here is a sample request JSON file containing the structure of an input.

```bash
$ cat request.json
$
{
  "nodes": {
    "1": "server1",
    "2": "server2",
    "3": "server3"
  },
  "toHash": "1"
}
```

Here is an example request :

```bash
$ curl -X POST -H "Content-Type: application/json" http://api-gateway-demo-endpoint.execute-api.com/servers -d request.json
```

Here is an example response to the above request :

```bash
{
  "Id": "1f2be9ee-03b1-45fe-90d2-b73e27561e9b",
  "response": "server1",
  "error": "",
}
```

