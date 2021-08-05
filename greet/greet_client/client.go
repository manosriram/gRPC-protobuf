package main

import (
	"context"
	"fmt"
	"grpc-ex/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("This is client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	req := &greetpb.GreetRequest{
		Greet: &greetpb.Greeting{
			FirstName: "Mano",
			LastName:  "Sriram",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Err while calling Greet() rpc %v", err)
	}
	log.Printf("Response from Greet() %+v", res)

	// fmt.Printf("Client created %f", c)

}
