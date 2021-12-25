package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/liridonrama/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.GreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)

	fN := req.GetGreeting().GetFirstName()
	lN := req.GetGreeting().GetLastName()

	result := fmt.Sprintf("Hello, %v %v", fN, lN)

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
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
