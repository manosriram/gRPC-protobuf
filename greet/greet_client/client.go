package main

import (
	"context"
	"fmt"
	"grpc-ex/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func greet_driver(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greet: &greetpb.Greeting{
			FirstName: "Mano",
			LastName:  "Sriram",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Err while calling Greet() rpc %v", err)
	}
	log.Printf("Response from Greet() %+v", res)
}

func greet_many_times_driver(c greetpb.GreetServiceClient) {
	greet_many_request := &greetpb.GreetManyTimesRequest{
		Greet: &greetpb.Greeting{
			FirstName: "Mano",
			LastName:  "Sriram",
		},
	}
	res_stream, _ := c.GreetManyTimes(context.Background(), greet_many_request)
	for {
		msg, err := res_stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(msg)
	}
}

func long_greet_driver(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Client streaming rpc")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Err while calling LongGreet() rpc %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Mano",
				LastName:  "Sriram",
			},
		},
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Virat",
				LastName:  "Kohli",
			},
		},
		&greetpb.LongGreetRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Michael",
				LastName:  "Jordan",
			},
		},
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Err while receiving response from LongGreet() %v", err)
	}
	fmt.Println(resp)
}

func main() {
	fmt.Println("This is client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// greet_many_times_driver(c)
	long_greet_driver(c)
}
