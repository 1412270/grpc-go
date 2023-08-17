package main

import (
	"context"
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
	callSum(client)
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
