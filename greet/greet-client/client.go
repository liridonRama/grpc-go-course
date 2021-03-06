package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/liridonrama/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	var withCreds grpc.DialOption

	tls, err := strconv.ParseBool(os.Getenv("TLS_ENABLED"))
	if err != nil {
		tls = false
	}

	fmt.Println(tls)

	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("could not connect %v", err)
		}
		withCreds = grpc.WithTransportCredentials(creds)
	} else {
		withCreds = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	cc, err := grpc.Dial("localhost:6543", withCreds)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	// doUnaryWithDeadline(c, 5*time.Second)
	// doUnaryWithDeadline(c, 1*time.Second)
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

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	start := time.Now()

	fmt.Println("starting unary request")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Liridon",
			LastName:  "Rama",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statuserr, ok := status.FromError(err)
		if ok {
			if statuserr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Println("unexpected error:", statuserr.Details())
			}
		} else {
			log.Fatalln("error while calling GreetWithDeadline RPC:", err)
		}

		return
	}

	fmt.Println(res.GetResult())

	fmt.Println(time.Since(start))
}
