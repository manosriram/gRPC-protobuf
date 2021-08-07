package main

import (
	"context"
	"fmt"
	"grpc-ex/average/averagepb"
	"log"

	"google.golang.org/grpc"
)

func average_driver(c averagepb.AverageServiceClient) error {
	stream, _ := c.Average(context.Background())

	requests := []*averagepb.AverageRequest{
		&averagepb.AverageRequest{
			Avg: &averagepb.Average{
				X: 12,
			},
		},
		&averagepb.AverageRequest{
			Avg: &averagepb.Average{
				X: 21,
			},
		},
		&averagepb.AverageRequest{
			Avg: &averagepb.Average{
				X: 301,
			},
		},
	}

	for _, val := range requests {
		stream.Send(val)
	}

	resp, _ := stream.CloseAndRecv()
	fmt.Println(resp.GetResult())
	return nil
}

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer cc.Close()

	c := averagepb.NewAverageServiceClient(cc)
	average_driver(c)
}
