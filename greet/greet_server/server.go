package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-go/greet/greetpb"
	"log"
	"net"
)

type server struct {
	greetpb.GreetServiceServer
}

func main() {
	fmt.Println("Hello, World!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fail to listen: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
