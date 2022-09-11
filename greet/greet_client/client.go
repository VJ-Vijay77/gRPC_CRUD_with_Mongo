package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/VJ-Vijay77/gRPC/greet/greetpb"
	"google.golang.org/grpc"
)


func main() {
	fmt.Println("Hello im a client!")

	cc,err := grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err != nil {
		log.Fatalln("could not connect ",err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created Client:%f",c)
	//doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
	
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming rpc...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vijay",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ajay",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
	}

	stream,err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error calling clinet streaming : %v",err)
	}
	for _,req := range requests {
	fmt.Printf("sending req :%v\n",req)
	stream.Send(req)
	time.Sleep(1500 *time.Millisecond)
	}
	res,err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error recieving response from LongGreet %v",err)
	}
	fmt.Printf("LongGreet Response: %v\n",res)
}



func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming rpc...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vijay",
			LastName: "Dinesh",
		},
	}
	res,err := c.GreetManyTimes(context.Background(),req)
	if err != nil {
		log.Fatalln(err)
	}
	for{
	msg,err := res.Recv()
	if err == io.EOF{
		break
	}
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(msg.GetResult())
	}
}



func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vijay",
			LastName: "Dinesh",
		},
	}

	res,err := c.Greet(context.Background(),req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC : %v", err)
	}
	log.Printf("Response from Greet : %v",res.Result)
}