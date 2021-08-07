package main

import (
	"fmt"
	"grpc-ex/average/averagepb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	averagepb.UnimplementedAverageServiceServer
}

func (*server) Average(stream averagepb.AverageService_AverageServer) error {
	fmt.Println("Average() invoked with stream")
	var arr []int32
	for {
		req, err := stream.Recv()

		x := req.GetAvg().GetX()

		if err == io.EOF {
			break
		}
		arr = append(arr, x)
	}

	sum := int32(0)
	for _, val := range arr {
		sum += val
	}
	avg := float64(sum) / float64(len(arr))

	return stream.SendAndClose(&averagepb.AverageResponse{
		Result: avg,
	})
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	averagepb.RegisterAverageServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
