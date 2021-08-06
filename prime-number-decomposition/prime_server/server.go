package main

import (
	"fmt"
	"grpc-ex/prime-number-decomposition/prime_pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	prime_pb.UnimplementedPrimeServiceServer
}

func (*server) PrimeDecomposition(req *prime_pb.PrimeRequest, stream prime_pb.PrimeService_PrimeDecompositionServer) error {
	x := req.GetPrime().GetX()
	k := int32(2)
	N := x
	fmt.Println("PrimeDecomposition() invoked")

	for N > 1 {
		if N%k == 0 {
			resp := &prime_pb.PrimeResponse{
				Number: k,
			}
			stream.Send(resp)
			N = N / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	prime_pb.RegisterPrimeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
}
