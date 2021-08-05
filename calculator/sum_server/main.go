package main

import (
	"context"
	"fmt"
	"grpc-ex/calculator/sumpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	sumpb.UnimplementedSumServiceServer
}

func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	log.Printf("Sum() was invoked: %v", req)
	x := req.GetSum().GetX()
	y := req.GetSum().GetY()

	res := &sumpb.SumResponse{
		Z: x + y,
	}
	return res, nil
}

func main() {
	fmt.Println("Sum Server Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
