package main

import (
	"context"
	"fmt"
	"grpc-ex/prime-number-decomposition/prime_pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start client, %v", err)
	}

	defer cc.Close()

	c := prime_pb.NewPrimeServiceClient(cc)

	prime_request := &prime_pb.PrimeRequest{
		Prime: &prime_pb.Prime{
			X: 120,
		},
	}

	prime_response, _ := c.PrimeDecomposition(context.Background(), prime_request)

	for {
		msg, err := prime_response.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(msg)

		// Not needed
		time.Sleep(1000 * time.Millisecond)
	}

}
