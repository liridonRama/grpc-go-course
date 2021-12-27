package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/liridonrama/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)

	fN, lN := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()

	result := fmt.Sprintf("Hello, %v %v", fN, lN)

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimesFunction was invoked:", req)

	fN := req.Greeting.GetFirstName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %v number %v", fN, strconv.Itoa(i+1)),
		}

		stream.Send(res)
		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	mux, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	err = s.Serve(mux)
	if err != nil {
		log.Fatal("failed to listen", err)
	}
}
