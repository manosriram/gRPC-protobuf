package main

import (
	"context"
	"fmt"
	"grpc-ex/maximum/maximumpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func stream_maximum(c maximumpb.MaximumServiceClient) {
	fmt.Println("Streaming Maximum numbers")
	stream, _ := c.MaximumStream(context.Background())
	var arr = [6]int32{1, 2, 100, 23, 200, 312}
	ch := make(chan struct{})

	// Send numbers to server
	go func() {
		for _, val := range arr {
			fmt.Printf("Sent: %v\n", val)
			stream.Send(&maximumpb.MaximumRequest{
				Maximum: &maximumpb.Maximum{
					Mx: val,
				}})
		}
	}()

	// Receive if a maximum number is hit
	go func() {
		for {
			num, err := stream.Recv()
			if err == io.EOF {
				close(ch)
			}
			if err != nil {
				log.Fatalf("Failed to listen %v", err)
				close(ch)
			}
			fmt.Printf("Received: %v\n", num)
		}
	}()

	<-ch
}

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	defer cc.Close()

	c := maximumpb.NewMaximumServiceClient(cc)
	stream_maximum(c)
}
