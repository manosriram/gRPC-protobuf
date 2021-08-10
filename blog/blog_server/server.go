package main

import (
	"context"
	"fmt"
	"grpc-ex/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	blogpb.BlogServiceServer
}

var collection *mongo.Collection

type blog_item struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func dataToBlog(data *blog_item) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.Id.Hex(),
		Title:    data.Title,
		Content:  data.Content,
		AuthorId: data.AuthorId,
	}
}

func (*server) Create(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	author_id := blog.GetAuthorId()
	content := blog.GetContent()
	title := blog.GetContent()

	data := blog_item{
		Title:    title,
		Content:  content,
		AuthorId: author_id,
	}
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Internal error %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "Cannot convert oid")
	}
	data.Id = oid
	return &blogpb.CreateBlogResponse{
		Blog: dataToBlog(&data),
	}, nil
}

func (*server) Update(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())

	updated_blog_item := &blog_item{}
	updated_blog_item.AuthorId = blog.GetAuthorId()
	updated_blog_item.Content = blog.GetContent()
	updated_blog_item.Title = blog.GetTitle()
	updated_blog_item.Id = oid

	filter := bson.M{"_id": oid}
	res, err := collection.ReplaceOne(ctx, filter, updated_blog_item)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Internal error %v", err))
	}

	fmt.Println(res)
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlog(updated_blog_item),
	}, nil
}

func (*server) Delete(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	doc_id := req.GetDocId()
	oid, err := primitive.ObjectIDFromHex(doc_id)
	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Internal error %v", err))
	}

	return &blogpb.DeleteBlogResponse{
		DeletedCount: int32(res.DeletedCount),
	}, nil
}

func (*server) Read(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	doc_id := req.GetDocId()
	oid, err := primitive.ObjectIDFromHex(doc_id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Internal error %v", err))
	}
	filter := bson.M{"_id": oid}

	result_data := blog_item{}
	collection.FindOne(ctx, filter).Decode(&result_data)

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlog(&result_data),
	}, nil
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
	reflection.Register(s)

	collection = client.Database("mydb").Collection("blog")

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
