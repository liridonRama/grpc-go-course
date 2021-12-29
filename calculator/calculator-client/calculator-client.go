package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
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

	// doUnary(c)

	// doServerStream(c)

	// doClientStream(c)

	doBiDiStream(c)
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

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	wSream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalln("Error while trying to start stream")
	}

	nums := []float64{100.1, 20.5}

	for _, num := range nums {
		wSream.Send(&calculatorpb.ComputeAverageRequest{
			Number: num,
		})
	}

	res, err := wSream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error while trying to reveice response:", err)
	}

	fmt.Println("Res received from server:", res.GetResult())

}

func doBiDiStream(c calculatorpb.CalculatorServiceClient) {
	wG := sync.WaitGroup{}
	wSream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalln("Error while trying to start stream")
	}

	nums := []int32{100, 20, 400, -20, 1440, 12, 1660, 12345}

	wG.Add(1)
	go func() {

		for _, num := range nums {
			wSream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})
		}

		wSream.CloseSend()
		wG.Done()
	}()

	wG.Add(1)
	go func() {
		for {
			res, err := wSream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("Error while trying to start stream")
			}

			fmt.Println(res.GetResult())
		}

		wG.Done()
	}()

	wG.Wait()
}
