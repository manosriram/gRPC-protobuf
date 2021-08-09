package main

import (
	"context"
	"fmt"
	"grpc-ex/calculator/sumpb"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) SquareRoot(ctx context.Context, req *sumpb.SquareRootRequest) (*sumpb.SquareRootResponse, error) {
	num := req.GetNumber()

	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Number less than 0")
	}
	return &sumpb.SquareRootResponse{
		Result: math.Sqrt(float64(num)),
	}, nil
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
