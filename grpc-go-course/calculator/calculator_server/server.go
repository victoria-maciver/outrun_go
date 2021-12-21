package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/angelus_reprobi/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function invoked with %v", req)
	sum := req.Calculation.A + req.Calculation.B
	res := &calculatorpb.SumResponse{
		Result: sum,
	}

	return res, nil
}

func (*server) Prime(req *calculatorpb.PrimeRequest, stream calculatorpb.CalculatorService_PrimeServer) error {
	var k int32
	k = 2
	number := req.GetPrime().GetNumber()

	for number > 1 {
		if number%k == 0 {
			res := &calculatorpb.PrimeResponse{
				Result: k,
			}
			stream.Send(res)
			number = number / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	var result int32
	var count int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float32(result) / float32(count)
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		result += req.GetNumber()
		count++
	}

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		i := req.GetNumber()
		if i > max {
			max = i
			sendErr := stream.Send(&calculatorpb.MaximumResponse{
				Maximum: max,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to the client: %v", sendErr)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received negative number: %v", number))
	}
	return &calculatorpb.SquareRootResponse{
		Root: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Running calculator server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	-calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
