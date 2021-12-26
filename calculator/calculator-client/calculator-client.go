package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/liridonrama/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ITERATIONS = 1

func main() {
	wg := sync.WaitGroup{}

	cc, err := grpc.Dial("localhost:6543", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewSumServiceClient(cc)

	start := time.Now()
	for i := 0; i < ITERATIONS; i++ {
		wg.Add(1)
		go getResult([]float64{rand.Float64()}, c, &wg)
	}

	wg.Wait()

	fmt.Println(time.Since(start))
}

func getResult(values []float64, c calculatorpb.SumServiceClient, wg *sync.WaitGroup) {
	res, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			Values: []float64{1.1, 1.1},
		},
	})
	if err != nil {
		log.Fatalf("request failed %v", err)
	}

	fmt.Println(res.GetResult())

	wg.Done()
}
