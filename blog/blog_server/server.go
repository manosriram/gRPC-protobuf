package main

import (
	"context"
	"fmt"
	"grpc-ex/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog Service Started")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error connecting mongodb")
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to listen %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Closing MongoDB Connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB : %v", err)
	}

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("End of Program")
}
