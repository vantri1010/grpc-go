package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-go/greet/greetpb"
	"log"
	"net"
)

// struct that represents the rpc server
type server struct {
	greetpb.GreetServiceServer
}

func main() {
	fmt.Println("Server started")

	// listen on the port
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fail to listen: %v", err)
	}

	log.Printf("server started at %v", lis.Addr())

	// creating new gRPC server and register the server
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	// now gRPC server accepts the request from the client
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
