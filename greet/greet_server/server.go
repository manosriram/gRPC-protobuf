package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	greetpb "grpc-ex/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreeting() was invoked with stream")
	result := "Hello, "

	for {
		req, err := stream.Recv()
		fmt.Printf("Receiving : %v", req)
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			return err
		}

		first_name := req.GetGreet().FirstName
		result += first_name + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("Invoked GreetEveryone()")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Printf("Error receiving data: %v", err)
			return err
		}

		first_name := req.GetGreet().GetFirstName()
		result := "Hello, " + first_name + "\n"

		stream.Send(&greetpb.GreetManyTimesResponse{
			Result: result,
		})
	}

}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline() was invoked with %v", req)
	for t := 0; t < 3; t++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Client canceled the call")
			return nil, status.Error(codes.DeadlineExceeded, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	first_name := req.GetGreet().GetFirstName()
	result := "Hello, " + first_name
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Starting Greet Server")
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
