package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/liridonrama/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:6543", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	start := time.Now()

	fmt.Println("starting unary request")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Liridon",
			LastName:  "Rama",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("could resolve request %v", err)
	}

	fmt.Println(res.GetResult())

	fmt.Println(time.Since(start))
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting client streaming rpc")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("error while calling long greet")
	}

	requests := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Liri"}},
		{Greeting: &greetpb.Greeting{FirstName: "Hans"}},
		{Greeting: &greetpb.Greeting{FirstName: "Ueli"}},
		{Greeting: &greetpb.Greeting{FirstName: "Peter"}},
		{Greeting: &greetpb.Greeting{FirstName: "John"}},
		{Greeting: &greetpb.Greeting{FirstName: "Steve"}},
		{Greeting: &greetpb.Greeting{FirstName: "Sage"}},
	}

	for _, req := range requests {
		fmt.Println("sending request", req)
		stream.Send(req)

		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error while receiving long greet")
	}

	fmt.Println("Long Greet Response:", res.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting server streaming rpc")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Liridon",
			LastName:  "Rama",
		},
	}

	rStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("error while retrieving stream:", err)
	}

	for {
		msg, err := rStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("error while retrieving messages from stream:", err)
		}

		log.Printf("response from GreetManyTimes: %v\n", msg.GetResult())
	}

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalln("error while retrieving stream:", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Liri"}},
		{Greeting: &greetpb.Greeting{FirstName: "Hans"}},
		{Greeting: &greetpb.Greeting{FirstName: "Ueli"}},
		{Greeting: &greetpb.Greeting{FirstName: "Peter"}},
		{Greeting: &greetpb.Greeting{FirstName: "John"}},
		{Greeting: &greetpb.Greeting{FirstName: "Steve"}},
		{Greeting: &greetpb.Greeting{FirstName: "Sage"}},
	}

	//
	waitC := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Println("Sending Message:", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}

		stream.CloseSend()
	}()

	go func() {
		defer close(waitC)

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("error while retrieving stream:", err)
			}

			fmt.Println("Received:", res.GetResult())
		}
	}()

	<-waitC
}
