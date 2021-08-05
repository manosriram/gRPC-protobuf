package main

import (
	"context"
	"grpc-ex/calculator/sumpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Sum Client Started")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer cc.Close()

	c := sumpb.NewSumServiceClient(cc)

	req := &sumpb.SumRequest{
		Sum: &sumpb.Sum{
			X: 3,
			Y: 10,
		},
	}
	log.Println(req)
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Response from sumservice: %v", res)
}
