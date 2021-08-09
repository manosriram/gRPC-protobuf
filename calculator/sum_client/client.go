package main

import (
	"context"
	"fmt"
	"grpc-ex/calculator/sumpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	log.Println("Sum Client Started")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer cc.Close()

	c := sumpb.NewSumServiceClient(cc)

	// _ := &sumpb.SumRequest{
	// Sum: &sumpb.Sum{
	// X: 3,
	// Y: 10,
	// },
	// }
	// log.Println(req)
	// res, err := c.Sum(context.Background(), req)
	res, err := c.SquareRoot(context.Background(), &sumpb.SquareRootRequest{
		Number: -123,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We sent a negative number")
			}
		} else {
			log.Fatalf("Error calling SquareRoot %v", err)
		}
	}

	log.Printf("Response from sumservice: %v", res)
}
