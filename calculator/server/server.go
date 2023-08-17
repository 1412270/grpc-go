package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/1412270/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum called")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
		} else {
			k++
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("average calculating...")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("err while cal avarage %v", err)
		}

		log.Printf("receive num %v", req.GetNum())
		total += req.GetNum()
		count++
	}
}

func (*server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	log.Println("find max calling...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while recv %v", err)
			return err
		}

		num := req.GetNum()
		if num > max {
			max = num
		}
		err = stream.Send(&calculatorpb.MaxResponse{
			Num: max,
		})
		if err != nil {
			log.Fatalf("err while find max %v", err)
			return err
		}
	}

}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Server connecting...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}
