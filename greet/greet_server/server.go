package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/VJ-Vijay77/gRPC/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{

}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n",req)
	firstName := req.GetGreeting().GetFirstName()
	result := "hello " + firstName

	res := greetpb.GreetResponse{
		Result: result,
	}

	return &res,nil

}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	
	for i := 0 ; i<10; i++ {
		result := "Hello "+firstName +" number"+strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 *time.Millisecond)
		
	}
	return nil
}


func(*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
fmt.Println("Long Greet Func was invoked with a streaming request")
	result := "Hello "
	for {
		req,err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v",err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName +"! "
	}

}


func main() {
	fmt.Println("Starting")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}

}
