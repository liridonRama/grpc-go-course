package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/liridonrama/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:6543", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

	doServerStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	start := time.Now()

	res, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			Values: []float64{1.1, 1.1},
		},
	})
	if err != nil {
		log.Fatalf("request failed %v", err)
	}

	fmt.Println(res.GetResult())

	fmt.Println(time.Since(start))
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 1902,
	})
	if err != nil {
		log.Fatalln("error while trying to retrieve stream:", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// reached end
			break
		}
		if err != nil {
			log.Fatalln("error while trying to retrieve stream:", err)
		}

		fmt.Println("Received prime number:", res.GetPrimeNumber())
	}
}