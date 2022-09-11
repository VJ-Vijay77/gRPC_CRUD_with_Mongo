package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/VJ-Vijay77/gRPC/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{

}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v \n",req)
	first := req.GetFirstNumber()
	second := req.GetSecondNumber()
	sum := first + second

	res := calculatorpb.SumResponse{
		Sum: sum,
	}

	return &res,nil

}


func(*server)PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest,stream calculatorpb.SumService_PrimeNumberDecompositionServer) error{
fmt.Println("Recieved Prime Numbers Decompostition ")
number := req.GetNumber()
divisor := 2

	for number > 1 {
		if number%int64(divisor) == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: int64(divisor),
			})
			number =  number / int64(divisor)
		}else{
			divisor ++
			fmt.Printf("divisor has increased to %v",divisor)
		}
	}
return nil
}

func main() {
	fmt.Println("Starting")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}

}
