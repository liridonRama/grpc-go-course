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

	doUnary(c)
	doServerStreaming(c)
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
