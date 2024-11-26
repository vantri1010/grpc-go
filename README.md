
An Golang-RPC example implementation
===

This project demonstrates a gRPC client in Go that interacts with a gRPC server. It includes multiple operations such as Unary RPC, Error Unary, Server Streaming, Client Streaming, Bi-Directional Streaming, and Unary RPC with deadlines.

> You can run directly on the same host.
> I prefer to use docker via docker compose or Dev Environments as the server and the host machine (windows) as the client to highlight the special features of RPC.

## Features
There are 3 services: Greet, Calculator and Blog implement features:
- Unary RPC (doUnary): Simple request-response interaction with the server.
- Error Unary RPC (doErrUnary): Handles errors from the server during Unary RPC.
- Server Streaming RPC (doServerStreaming): The server sends multiple responses for a single client request.
- Client Streaming RPC (doClientStreaming): The client sends multiple requests to the server, which responds once.
- Bi-Directional Streaming RPC (doBiDiStreaming): Both client and server send and receive streams of messages in real time.
- Unary With Deadline (doUnaryWithDeadline): Unary RPC with a timeout to control long-running requests.
- CRUD with PRC protocol connect to MongoDB

## Prerequisites
- [Go](https://go.dev/) installed on your system
- [Docker](https://www.docker.com/) and Docker compose installed (optional)
- [MongoDB](https://hub.docker.com/_/mongo) install (locally or on docker)
- Protocol Buffers compiler installed ([protoc](https://grpc.io/docs/protoc-installation/))
- ssl fully generated (using `./ssl/instruction.sh`)
- protoc file generated (using `./generate.sh`)

## Installation
Install the dependencies and devDependencies and start the server.
```
git clone https://github.com/vantri1010/grpc-go
cd grpc-go
go mod tidy
```
optional docker compose install :
```compose
docker-compose up 
```
### Blog usage
This will start a go server and connect to MongoDB (store blog data).
> You can run it in normal way in host OS or inside a docker cli inside app/grpc after docker compose up
```server blog
go run blog/blog_server/server.go
```
Compile and run the client side with appropriate flags:
- Create:
```create blog
go run blog/blog_client/client.go -operation=create -author="bash client" -title="on window" -content="to docker"
```
- Read:
```read blog
go run blog/blog_client/client.go -operation=read -id="BlogID just created"
```
- Update:
```update blog
go run blog/blog_client/client.go -operation=update -id="Copied BlogID" -author="Update author" -title="UpdatedTitle" -content="UpdatedContent"
```

- Delete:
```delete blog
go run blog/blog_client/client.go -operation=delete -id="6745fce9ca67268d0cfd1a4a"
```

### Calculator usage
Start a server
> You can run the calculator server side in normal way in host OS or inside a docker cli inside app/grpc after docker compose up
```server calculator
go run calculator/calculator_server/server.go
```
Compile and run the client side with appropriate flags:
- Unary RPC:
```doUnary
go run calculator/calculator_client/client.go -operation=doUnary
```
- Server Streaming RPC:
```doServerStreaming
go run calculator/calculator_client/client.go -operation=doServerStreaming
```
- Client Streaming RPC:
```doClientStreaming
go run calculator/calculator_client/client.go -operation=doClientStreaming
```
- Bi-Directional Streaming RPC:
```doBiDiStreaming
go run calculator/calculator_client/client.go -operation=doBiDiStreaming
```

### Greet usage
Start a server
> You can run the calculator server side in normal way in host OS or inside a docker cli inside app/grpc after docker compose up
```greet server
go run greet/greet_server/server.go
```
- Unary RPC:
```doUnary
go run greet/greet_client/client.go -operation=doUnary
```
- Error Unary RPC:
```doErrUnary
go run greet/greet_client/client.go -operation=doErrUnary
```
- Server Streaming RPC:
```doServerStreaming
go run greet/greet_client/client.go -operation=doServerStreaming
```
- Client Streaming RPC:
```doClientStreaming
go run greet/greet_client/client.go -operation=doClientStreaming
```
- Bi-Directional Streaming RPC:
```doBiDiStreaming
go run greet/greet_client/client.go -operation=doBiDiStreaming
```
- Unary With Deadline (default 3 seconds, customizable):
```doUnaryWithDeadline 
go run greet/greet_client/client.go -operation=doUnaryWithDeadline -deadline=5
```