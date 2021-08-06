package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	greetpb "grpc-ex/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet() was invoked with %v", req)
	first_name := req.GetGreet().GetFirstName()
	result := "Hello, " + first_name
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes() was invoked with %v", req)
	first_name := req.GetGreet().GetFirstName()
	for t := 0; t < 10; t++ {
		result := "Hello, " + first_name + " with iteration " + strconv.Itoa(t)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println("Server started")
}
