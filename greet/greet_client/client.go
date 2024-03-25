package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-go/greet/greetpb"
	"log"
)

func main() {
	// creates a client connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer conn.Close()
	// created gRPC client
	client := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client: %f", client)

}
