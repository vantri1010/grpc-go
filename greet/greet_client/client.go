package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc-go/greet/greetpb"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	// Define CLI flags
	operation := flag.String("operation", "", "Type of gRPC operation: doUnary, doErrUnary, doServerStreaming, doClientStreaming, doBiDiStreaming, doUnaryWithDeadline")
	deadline := flag.Int("deadline", 3, "Deadline in seconds for doUnaryWithDeadline operation (default: 3 seconds)")
	flag.Parse()

	if *operation == "" {
		fmt.Println("Error: operation flag is required")
		flag.Usage()
		return
	}

	// TLS and connection setup
	tls := true
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	// creates a client connection to the server
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer conn.Close()

	// created gRPC client
	client := greetpb.NewGreetServiceClient(conn)

	// Handle operations based on CLI input
	switch *operation {
	case "doUnary":
		doUnary(client)
	case "doErrUnary":
		doErrUnary(client)
	case "doServerStreaming":
		doServerStreaming(client)
	case "doClientStreaming":
		doClientStreaming(client)
	case "doBiDiStreaming":
		doBiDiStreaming(client)
	case "doUnaryWithDeadline":
		doUnaryWithDeadline(client, time.Duration(*deadline)*time.Second)
	default:
		log.Fatalf("Unknown operation: %v", *operation)
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Tri",
			LastName:  "Nguyen Van",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doErrUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do an error Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "",
			LastName:  "Nguyen Van",
		},
	}

	res, err := c.Greet(context.Background(), req)

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent an empty string of first name!")
				return
			}
		} else {
			log.Fatalf("Big Error calling Greet: %v", err)
			return
		}
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Tri",
			LastName:  "Nguyen Van",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tri",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nguyen",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Van",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tien",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sinh",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our LongGreetRequests slice and send each message individually to the server.
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	// Close the stream and retrieve the response.
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tri",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nguyen",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Van",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tien",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sinh",
			},
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start a goroutine to send messages to the server.
	go func() {
		defer wg.Done()
		// Loop through the greetings and send them to the server.
		for _, req := range requests {
			time.Sleep(1000 * time.Millisecond)
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
		}
		stream.CloseSend()
	}()

	// Start a goroutine to receive messages from the server.
	go func() {
		defer wg.Done()
		// Loop until the stream is closed.
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
	}()

	// Wait for both goroutines to complete.
	wg.Wait()
}

// doUnaryWithDeadline sends a greetings.GreetWithDeadlineRequest to the server and prints the response.
// The context is cancelled after the given timeout.
func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Tri",
			LastName:  "Nguyen",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
