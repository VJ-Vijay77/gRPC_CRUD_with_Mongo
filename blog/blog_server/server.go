package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/VJ-Vijay77/gRPC/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
}

//mongo struct model
type blogitem struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string `bson:"author_id"`
	Content  string `bson:"content"`
	Title    string `bson:"title"`
}

//implementing the interface in proto
func(*server) CreateBlog(ctx context.Context,req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := req.GetBlog()

	data := blogitem {
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}

}



func main() {

	//this shows the details like time and lines of codes when error or fatal occur
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	fmt.Println("Connecting to MongoDb")


	//? connecting mongo db with client
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
	}

	//creating a collection in mongodb
	collection = client.Database("mydb").Collection("blog")

	//connecting to port
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting the server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Print("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Stopping MongoDb")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")

}
