package main

import (
	"context"
	"io"
	"log"

	"github.com/1412270/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("error while dial %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)
	// log.Printf("service client %f", client)
	// callSum(client)
	// callPND(client)
	// callAverage(client)
	callMax(client)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})

	if err != nil {
		log.Fatalf("call sum api error %v", err)
	}
	log.Printf("sum api Response %v", resp.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 50,
	})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("prime number is %v", resp.GetResult())
	}
}

func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("sending nums to cal average...")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("err while cal average %v", err)
	}

	listReq := []calculatorpb.AverageRequest{
		{
			Num: 6,
		},
		{
			Num: 8,
		},
		{
			Num: 9.7,
		},
		{
			Num: 5,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("err while send num request %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err while receive average %v", err)
	}

	log.Printf("average response %v", resp.Result)
}

func callMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("call finding max...")
	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalf("err while cal max %v", err)
	}

	waitChannel := make(chan struct{})

	go func() {
		listReq := []calculatorpb.MaxRequest{
			{
				Num: 6,
			},
			{
				Num: 8,
			},
			{
				Num: 9,
			},
			{
				Num: 7,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("err while send num request %v", err)
				break
			}
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("find max complete")
				break
			}
			if err != nil {
				log.Fatalf("err while recv %v", err)
				break
			}
			log.Printf("max is: %v\n", resp.GetNum())
		}
		close(waitChannel)
	}()

	<-waitChannel
}
