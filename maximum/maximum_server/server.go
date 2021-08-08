package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"

	maximumpb "grpc-ex/maximum/maximumpb"

	"google.golang.org/grpc"
)

type server struct {
	maximumpb.UnimplementedMaximumServiceServer
}

func (*server) MaximumStream(stream maximumpb.MaximumService_MaximumStreamServer) error {
	fmt.Println("Invokeded max_stream()")
	var mn int32 = math.MinInt32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error receiving data: %v", err)
			return err
		}

		current_number := req.GetMaximum().GetMx()
		if current_number > mn {
			mn = current_number
			stream.Send(&maximumpb.MaxiumumResponse{
				Result: mn,
			})
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	maximumpb.RegisterMaximumServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println("Server started")

}
