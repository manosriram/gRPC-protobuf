package main

import (
	"context"
	"fmt"
	"grpc-ex/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func greet_everyone_driver(c greetpb.GreetServiceClient) {
	fmt.Println("Starting biDi streaming")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Err while receiving response from LongGreet() %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Mano",
				LastName:  "Sriram",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Virat",
				LastName:  "Kohli",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greet: &greetpb.Greeting{
				FirstName: "Michael",
				LastName:  "Jordan",
			},
		},
	}

	wc := make(chan struct{})

	go func() {
		for _, req := range requests {
			stream.Send(req)
			fmt.Printf("Sending: %v", req)
			time.Sleep(2 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				close(wc)
				break
			}
			if err != nil {
				log.Fatalf("Failed to listen %v", err)
				close(wc)
			}
			fmt.Printf("Received: %v", msg.GetResult())
		}
	}()

	<-wc
}

func greet_with_deadline_driver(c greetpb.GreetServiceClient, seconds time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greet: &greetpb.Greeting{
			FirstName: "Mano",
			LastName:  "Sriram",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), seconds)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			if status.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit, deadline exceeded")
			} else {
				fmt.Printf("Some error happened %v", err)
			}
		} else {
			log.Fatalf("Err while calling GreetWithDeadline() rpc %v", err)
		}
	}
	log.Printf("Response from GreetWithDeadline() %+v", res)
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
	// long_greet_driver(c)
	// greet_everyone_driver(c)
	greet_with_deadline_driver(c, 1*time.Second)
	greet_with_deadline_driver(c, 5*time.Second)
}
